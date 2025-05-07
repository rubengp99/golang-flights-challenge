package googleflights

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	infisical "github.com/infisical/go-sdk"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

type APIResponse struct {
	Status    bool            `json:"status"`
	Message   string          `json:"message"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
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
					SecretKey:   "RAPID_API_KEY",
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
	req.Header.Add("x-rapidapi-key", s.config.APIKey)
	req.Header.Add("x-rapidapi-host", s.config.BaseURL)
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
		offers   = FlightOffer{}
		request  = vendors.Request{

			BaseURL:  s.config.BaseURL,
			Resource: "api/v1/searchFlights",
			Method:   http.MethodGet,
			Params: url.Values{
				"departure_id":  []string{params.Origin},
				"arrival_id":    []string{params.Destination},
				"outbound_date": []string{params.Date.Format("2006-01-02")},
				"adults":        []string{params.Adults},
				"stops":         []string{"direct"}, // to keep things simple, only direct flights
				"currencyCode":  []string{"USD"},
			},
		}
	)

	if err := vendors.MakeHTTPRequest(s, request, &response); err != nil {
		log.Printf("unable to retrieve flights from google flights, error: %s", err)
		return FlightOffer{}, err
	}

	if err := json.Unmarshal(response.Data, &offers); err != nil {
		log.Printf("unable to decode flights from google flights, error: %s", err)
		return FlightOffer{}, err
	}

	return offers, nil
}
