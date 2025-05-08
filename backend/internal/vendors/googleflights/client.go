package googleflights

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	infisical "github.com/infisical/go-sdk"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors"
	"github.com/rubengp99/golang-flights-challenge/pkg"
	g "github.com/serpapi/google-search-results-golang"
)

type SearchMetadata struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	JSONEndpoint     string  `json:"json_endpoint"`
	CreatedAt        string  `json:"created_at"`
	ProcessedAt      string  `json:"processed_at"`
	GoogleFlightsURL string  `json:"google_flights_url"`
	RawHTMLFile      string  `json:"raw_html_file"`
	PrettifyHTMLFile string  `json:"prettify_html_file"`
	TotalTimeTaken   float64 `json:"total_time_taken"`
}

type SearchParameters struct {
	Engine       string `json:"engine"`
	Hl           string `json:"hl"`
	Gl           string `json:"gl"`
	DepartureID  string `json:"departure_id"`
	ArrivalID    string `json:"arrival_id"`
	OutboundDate string `json:"outbound_date"`
	ReturnDate   string `json:"return_date"`
	Currency     string `json:"currency"`
}

type APIResponse struct {
	FlightOffer
	SearchMetadata   SearchMetadata   `json:"search_metadata"`
	SearchParameters SearchParameters `json:"search_parameters"`
}

// Service is a representation of a google flights http client
type Service struct {
	config     vendors.Config
	httpclient *http.Client
}

// ConfigProviderFunc dinari config provider
type ConfigProviderFunc func(infclient infisical.InfisicalClientInterface, projectID string) vendors.Config

// DefaultConfigFromSecretsManager retrieves config from secrets manager
func DefaultConfigFromSecretsManager() ConfigProviderFunc {
	return func(infclient infisical.InfisicalClientInterface, projectID string) vendors.Config {
		var (
			c      = vendors.Config{}
			errors = make(chan error)
			wgdone = make(chan bool)
			wg     sync.WaitGroup
		)

		// retrieve all secrets from infisical
		secrets := []func(channel chan error){
			func(channel chan error) {
				APIKey, err := infclient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
					SecretKey:   "SERPAPI_API_KEY",
					Environment: os.Getenv("STAGE"),
					ProjectID:   projectID,
					SecretPath:  "/",
				})
				wg.Done()
				c.APIKey = APIKey.SecretValue
				if err != nil {
					channel <- err
				}
			},
			func(channel chan error) {
				clientID, err := infclient.Secrets().Retrieve(infisical.RetrieveSecretOptions{
					SecretKey:   "GOOGLE_FLIGHTS_BASE_URL",
					Environment: os.Getenv("STAGE"),
					ProjectID:   projectID,
					SecretPath:  "/",
				})
				wg.Done()
				c.BaseURL = clientID.SecretValue
				if err != nil {
					channel <- err
				}
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
		return c
	}
}

// NewService returns a new google flights service
func NewService(c ConfigProviderFunc, infclient infisical.InfisicalClientInterface, projectID string) Service {
	client := http.DefaultClient
	client.Timeout = 60 * time.Second
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	client.Transport = transport

	return Service{
		config:     c(infclient, projectID),
		httpclient: client,
	}
}

// Authenticate generates all the necessary http headers and settings required by our integration in order to authorize our requests
func (s *Service) Authenticate(req *http.Request) error {
	// auth on query params here
	return nil
}

// Client exports local underlying http client settings for the current integration
func (s Service) Client() *http.Client {
	return s.httpclient
}

// RetrieveFlightOffers retrives all available flight offers from google flights
func (s *Service) RetrieveFlightOffers(params pkg.QueryParams) (FlightOffer, error) {
	var (
		response APIResponse
		request  = map[string]string{
			"engine":        "google_flights",
			"departure_id":  params.Origin,
			"arrival_id":    params.Destination,
			"outbound_date": params.Date.Format("2006-01-02"),
			"adults":        params.Adults,
			"stops":         "direct", // to keep things simple, only direct flights
			"currencyCode":  "USD",
			"type":          "2", // one way
			"hl":            "en",
		}
	)

	search := g.NewGoogleSearch(request, s.config.APIKey)
	results, err := search.GetJSON()
	if err != nil {
		log.Printf("unable to retrieve flights from google flights, error: %s", err)
		return FlightOffer{}, err
	}

	bodyBytes, err := json.Marshal(results)
	if err != nil {
		log.Printf("unable to marshal flights from google flights, error: %s", err)
		return FlightOffer{}, err
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Printf("unable to marshal flights from google flights, error: %s", err)
		return FlightOffer{}, err
	}

	return response.FlightOffer, nil
}
