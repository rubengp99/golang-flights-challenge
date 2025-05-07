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
)

// App is the representation of all the functionality exposed in this application
type App struct {
	LogWriter                   io.Writer
	SecretKey                   string
	GetBestFlightsHandler       http.HandlerFunc
	GetLocationIATACodesHandler http.HandlerFunc
	LoginHandler                http.HandlerFunc
}

// Options is a type for application options to modify the app
type Options func(o *Option)

type ProvideInfisicalClientFunc func() infisical.InfisicalClientInterface

// Option is a representation of configurable options for the app
type Option struct {
	LogWriter                  io.Writer
	ProjectUD                  string
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
			return infisical.NewInfisicalClient(context.Background(), infisical.Config{})
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
		amadeusClient       amadeus.Service
		flightskyClient     flightsky.Service
		googleflightsClient googleflights.Service
	)

	infisicalClient := o.ProvideInfisicalClient()
	// retrieve all secrets from infisical
	secrets := []func(channel chan error){
		func(channel chan error) {
			APIKey, err := infisicalClient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
				SecretKey:   "JOBSITY_SECRET_KEY",
				Environment: os.Getenv("STAGE"),
				ProjectID:   o.ProjectUD,
				SecretPath:  "/",
			})
			wg.Done()
			secretKey = APIKey.SecretValue
			channel <- err
		},
		func(channel chan error) {
			googleflightsClient = googleflights.NewService(o.ProvideGoogleflightsConfig, infisicalClient, o.ProjectUD)
			wg.Done()
			channel <- nil
		},
		func(channel chan error) {
			amadeusClient = amadeus.NewService(o.ProvideAmadeusConfig, infisicalClient, o.ProjectUD)
			wg.Done()
			channel <- nil
		},
		func(channel chan error) {
			flightskyClient = flightsky.NewService(o.ProvideFlightskyConfig, infisicalClient, o.ProjectUD)
			wg.Done()
			channel <- nil
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

	return App{
		SecretKey:             secretKey,
		LogWriter:             o.LogWriter,
		GetBestFlightsHandler: RetrieveBestFlightsHandler(googleflightsClient, amadeusClient, flightskyClient),
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
		AllowedMethods: []string{"GET", "OPTIONS"}, // this APP only includes GET APIs for demo purposes
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
	})

	// no auth required routes
	router.Group(func(r chi.Router) {
		r.Use(newMiddleware(a.LogWriter, a.SecretKey, true).Wrap)
		r.Get("/login", a.GetBestFlightsHandler)
	})

	return router.ServeHTTP
}
