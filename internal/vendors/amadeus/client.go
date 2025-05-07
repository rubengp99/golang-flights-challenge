package amadeus

import (
	"encoding/json"
	"fmt"
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

type AuthResponse struct {
	Type            string `json:"type"`
	Username        string `json:"username"`
	ApplicationName string `json:"application_name"`
	ClientID        string `json:"client_id"`
	TokenType       string `json:"token_type"`
	AccessToken     string `json:"access_token"`
	ExpiresIn       int    `json:"expires_in"`
	State           string `json:"state"`
	Scope           string `json:"scope"`
}

type Meta struct {
	Count int       `json:"count"`
	Links MetaLinks `json:"links"`
}

type MetaLinks struct {
	Self string `json:"self"`
}

type APIResponse struct {
	Meta         Meta            `json:"meta"`
	Data         json.RawMessage `json:"data"`
	Dictionaries json.RawMessage `json:"dictionaries"` // just for documentation, we don't really need this
}

// Service is a representation of a Amadeus http client
type Service struct {
	config     vendors.Config
	httpclient *http.Client
}

// ConfigProviderFunc dinari config provider
type ConfigProviderFunc func(client infisical.InfisicalClientInterface, projectID string) vendors.Config

// DefaultConfigFromSecretsManager retrieves config from secrets manager
func DefaultConfigFromSecretsManager() ConfigProviderFunc {
	return func(client infisical.InfisicalClientInterface, projectID string) vendors.Config {
		var (
			c      = vendors.Config{}
			errors = make(chan error)
			wgdone = make(chan bool)
			wg     sync.WaitGroup
		)

		// retrieve all secrets from infisical
		secrets := []func(channel chan error){
			func(channel chan error) {
				clientID, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
					SecretKey:   "AMADEUS_CLIENT_ID",
					Environment: os.Getenv("STAGE"),
					ProjectID:   projectID,
					SecretPath:  "/",
				})
				wg.Done()
				c.ClientID = clientID.SecretValue
				if err != nil {
					channel <- err
				}
			},
			func(channel chan error) {
				clientID, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
					SecretKey:   "AMADEUS_SECRET_ID",
					Environment: os.Getenv("STAGE"),
					ProjectID:   projectID,
					SecretPath:  "/",
				})
				wg.Done()
				c.ClientSecret = clientID.SecretValue
				if err != nil {
					channel <- err
				}
			},
			func(channel chan error) {
				clientID, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
					SecretKey:   "AMADEUS_BASE_URL",
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

// NewService returns a new amadeus service
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
	var (
		response AuthResponse
		request  = vendors.Request{
			SkipAuth: true,
			BaseURL:  s.config.BaseURL,
			Resource: "v1/security/oauth2/token",
			Method:   http.MethodPost,
			Payload: url.Values{
				"client_id":     []string{s.config.ClientID},
				"client_secret": []string{s.config.ClientSecret},
				"grant_type":    []string{s.config.ClientID},
			},
		}
	)

	if err := vendors.MakeHTTPRequest(s, request, &response); err != nil {
		log.Printf("unable to authenticate with amadeus, error: %s", err)
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", response.TokenType, response.AccessToken))
	return nil
}

// Client exports local underlying http client settings for the current integration
func (s Service) Client() *http.Client {
	return s.httpclient
}

// RetrieveFlightOffers retrives all available flight offers from amadeus
func (s *Service) RetrieveFlightOffers(params pkg.QueryParams) ([]FlightOffer, []Airline, error) {
	var (
		response APIResponse
		offers   = []FlightOffer{}
		request  = vendors.Request{

			BaseURL:  s.config.BaseURL,
			Resource: "v2/shopping/flight-offers",
			Method:   http.MethodGet,
			Params: url.Values{
				"originLocationCode":      []string{params.Origin},
				"destinationLocationCode": []string{params.Destination},
				"departureDate":           []string{params.Date.Format("2006-01-02")},
				"adults":                  []string{params.Adults},
				"nonStop":                 []string{"true"}, // to keep things simple, only direct flights
				"currencyCode":            []string{"USD"},  // amadeus uses EUR as default, so we need to specify this
			},
		}
	)

	if err := vendors.MakeHTTPRequest(s, request, &response); err != nil {
		log.Printf("unable to retrieve flights from amadeus, error: %s", err)
		return nil, nil, err
	}

	if err := json.Unmarshal(response.Data, &offers); err != nil {
		log.Printf("unable to decode flights from amadeus, error: %s", err)
		return nil, nil, err
	}

	// unique airline codes
	dedupeAirlineCodes := map[string]bool{}
	airlineCodes := []string{}
	for _, o := range offers {
		for _, code := range o.ValidatingAirlineCodes {
			if !dedupeAirlineCodes[code] {
				dedupeAirlineCodes[code] = true
				airlineCodes = append(airlineCodes, code)
			}
		}
	}

	airlines, err := s.retrieveAirlines(airlineCodes)
	if err != nil {
		log.Printf("unable to retrieve airlines from amadeus, error: %s", err)
		return nil, nil, err
	}

	return offers, airlines, nil
}

func (s *Service) retrieveAirlines(codes []string) ([]Airline, error) {
	var (
		response APIResponse
		airlines = []Airline{}
		request  = vendors.Request{
			BaseURL:  s.config.BaseURL,
			Resource: "v1/reference-data",
			Method:   http.MethodGet,
			Params: url.Values{
				"airlineCodes": codes,
			},
		}
	)

	if err := vendors.MakeHTTPRequest(s, request, &response); err != nil {
		log.Printf("unable to retrieve flights from amadeus, error: %s", err)
		return nil, err
	}

	if err := json.Unmarshal(response.Data, &airlines); err != nil {
		log.Printf("unable to decode flights from amadeus, error: %s", err)
		return nil, err
	}

	return airlines, nil
}
