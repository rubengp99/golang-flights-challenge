package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	infisical "github.com/infisical/go-sdk"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/redis"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

// App is the representation of all the functionality exposed in this application
type App struct {
	LogWriter                           io.Writer
	SecretKey                           string
	GetBestFlightsHandler               http.HandlerFunc
	LoginHandler                        http.HandlerFunc
	SubcribeToFlightOfferUpdatesHandler http.HandlerFunc
}

// Options is a type for application options to modify the app
type Options func(o *Option)

type ProvideInfisicalClientFunc func() infisical.InfisicalClientInterface

// Option is a representation of configurable options for the app
type Option struct {
	LogWriter                  io.Writer
	ProjectUD                  string
	DisableRedis               bool
	TimeProvider               func() time.Time
	ProvideInfisicalClient     ProvideInfisicalClientFunc
	ProvideAmadeusConfig       amadeus.ConfigProviderFunc
	ProvideFlightskyConfig     flightsky.ConfigProviderFunc
	ProvideGoogleflightsConfig googleflights.ConfigProviderFunc
}

// New returns an instance of the default app
func New(options ...Options) App {
	o := Option{
		LogWriter:                  os.Stdout,
		TimeProvider:               time.Now,
		ProjectUD:                  os.Getenv("PROJECT_ID"),
		ProvideAmadeusConfig:       amadeus.DefaultConfigFromSecretsManager(),
		ProvideFlightskyConfig:     flightsky.DefaultConfigFromSecretsManager(),
		ProvideGoogleflightsConfig: googleflights.DefaultConfigFromSecretsManager(),
		ProvideInfisicalClient: func() infisical.InfisicalClientInterface {
			client := infisical.NewInfisicalClient(context.Background(), infisical.Config{})
			client.Auth().SetAccessToken(os.Getenv("INFISICAL_TOKEN"))
			return client
		},
	}

	for _, option := range options {
		option(&o)
	}

	var (
		errors              = make(chan error)
		wgdone              = make(chan bool)
		wg                  sync.WaitGroup
		secretKey           = ""
		clientID            = ""
		clientSecret        = ""
		amadeusClient       amadeus.Service
		flightskyClient     flightsky.Service
		googleflightsClient googleflights.Service
	)

	infisicalClient := o.ProvideInfisicalClient()
	redisClient := redis.NewRedisService(o.DisableRedis)
	// retrieve all secrets from infisical
	secrets := []func(channel chan error){
		func(channel chan error) {
			secret, err := infisicalClient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
				SecretKey:   "JOBSITY_SECRET_KEY",
				Environment: os.Getenv("STAGE"),
				ProjectID:   o.ProjectUD,
				SecretPath:  "/",
			})
			secretKey = secret.SecretValue
			if err != nil {
				channel <- err
			}
			wg.Done()
		},
		func(channel chan error) {
			secret, err := infisicalClient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
				SecretKey:   "JOBSITY_APP_CLIENT_ID",
				Environment: os.Getenv("STAGE"),
				ProjectID:   o.ProjectUD,
				SecretPath:  "/",
			})
			clientID = secret.SecretValue
			if err != nil {
				channel <- err
			}
			wg.Done()
		},
		func(channel chan error) {
			secret, err := infisicalClient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
				SecretKey:   "JOBSITY_APP_CLIENT_SECRET",
				Environment: os.Getenv("STAGE"),
				ProjectID:   o.ProjectUD,
				SecretPath:  "/",
			})
			clientSecret = secret.SecretValue
			if err != nil {
				channel <- err
			}
			wg.Done()
		},
		func(channel chan error) {
			googleflightsClient = googleflights.NewService(o.ProvideGoogleflightsConfig, infisicalClient, o.ProjectUD)
			wg.Done()
		},
		func(channel chan error) {
			amadeusClient = amadeus.NewService(o.ProvideAmadeusConfig, infisicalClient, o.ProjectUD)
			wg.Done()
		},
		func(channel chan error) {
			flightskyClient = flightsky.NewService(o.ProvideFlightskyConfig, infisicalClient, o.ProjectUD)
			wg.Done()
		},
	}

	wg.Add(len(secrets))
	go func() {
		wg.Wait()
		close(wgdone)
	}()

	for _, f := range secrets {
		go f(errors)
	}

	select {
	case <-wgdone:
		break
	case err := <-errors:
		close(errors)
		// we cannot proceed after this point, so we panic
		panic(err)
	}

	// default credentials for application, simple for demo purposes
	appCredentials := pkg.CrendetialsRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	return App{
		SecretKey:                           secretKey,
		LogWriter:                           o.LogWriter,
		LoginHandler:                        LoginHandler(appCredentials, secretKey),
		GetBestFlightsHandler:               RetrieveBestFlightsHandler(redisClient, googleflightsClient, amadeusClient, flightskyClient),
		SubcribeToFlightOfferUpdatesHandler: SubcribeToFlightOfferUpdatesHandler(redisClient, googleflightsClient, amadeusClient, flightskyClient),
	}
}

// Handler returns the main http handler for the application
func (a *App) Handler() http.HandlerFunc {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		/*AllowOriginFunc: func(r *http.Request, origin string) bool {
			// Modify this if we want to block origins some day
			return true
		},*/
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // this APP only includes GET/POST APIs for demo purposes
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Authorization",
		},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
	}))

	// auth required routes
	router.Group(func(r chi.Router) {
		r.Use(newMiddleware(a.LogWriter, a.SecretKey, true).Wrap)
		r.Get("/flights/search", a.GetBestFlightsHandler)
		r.Get("/subscribe", a.SubcribeToFlightOfferUpdatesHandler)
	})

	// no auth required routes
	router.Group(func(r chi.Router) {
		r.Use(newMiddleware(a.LogWriter, a.SecretKey, false).Wrap)
		r.Post("/login", a.LoginHandler)
	})

	router.Options("/flights/search", defaultOptionsHandler)
	router.Options("/login", defaultOptionsHandler)
	router.Options("/subscribe", defaultOptionsHandler)

	return router.ServeHTTP
}
