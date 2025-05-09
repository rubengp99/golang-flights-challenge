package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	infisical "github.com/infisical/go-sdk"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/testhelpers"
	"github.com/rubengp99/golang-flights-challenge/pkg"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("STAGE", "dev")
	os.Setenv("PROJECT_ID", "testInfisical")
	os.Setenv("INFISICAL_TOKEN", "testInfisicalToken")
}

func mockInfisicalClient(baseURL string) Options {
	return func(o *Option) {
		o.ProvideInfisicalClient = func() infisical.InfisicalClientInterface {
			return infisical.NewInfisicalClient(context.Background(), infisical.Config{
				SiteUrl: baseURL,
			})
		}
	}
}

func mockGoogleflightsConfig() Options {
	return func(o *Option) {
		o.ProvideGoogleflightsConfig = func(infclient infisical.InfisicalClientInterface, projectID string) vendors.Config {
			return vendors.Config{
				Disabled: true,
			}
		}
	}
}

func mockAmadeusServer(t *testing.T) *httptest.Server {
	run := testhelpers.Run(t)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/v1/reference-data/airlines?airlineCodes=TG%2CQF":
			var response amadeus.APIResponse
			reader := testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-airlines.json"), &response)
			data, err := io.ReadAll(reader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			run("Method as expected", func(t *testing.T) {
				assert.Equal(t, http.MethodGet, r.Method)
			})

			run("Token as expected", func(t *testing.T) {
				assert.Equal(t, "Bearer TestAccessToken", r.Header.Get("Authorization"))
			})

			w.WriteHeader(http.StatusOK)
			w.Write(data)
			break
		case "/v2/shopping/flight-offers?adults=1&currencyCode=USD&departureDate=2025-05-09&destinationLocationCode=BKK&nonStop=true&originLocationCode=SYD":
			var response amadeus.APIResponse
			reader := testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-offers.json"), &response)
			data, err := io.ReadAll(reader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			run("Method as expected", func(t *testing.T) {
				assert.Equal(t, http.MethodGet, r.Method)
			})

			run("Token as expected", func(t *testing.T) {
				assert.Equal(t, "Bearer TestAccessToken", r.Header.Get("Authorization"))
			})

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		case "/v1/security/oauth2/token?":
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			run("Payload is as expected", func(t *testing.T) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "testClientId", r.Form.Get("client_id"))
				assert.Equal(t, "testClientSecret", r.Form.Get("client_secret"))
				assert.Equal(t, "client_credentials", r.Form.Get("grant_type"))
			})

			response := amadeus.AuthResponse{
				TokenType:   "Bearer",
				AccessToken: "TestAccessToken",
			}
			data, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		default:
			t.Errorf("unexpected path %s", r.URL.String())
			t.FailNow()
		}

	}))
}

func mockFlightskyServer(t *testing.T) *httptest.Server {
	run := testhelpers.Run(t)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/flights/search-one-way?adults=1&currencyCode=USD&departDate=2025-05-09&fromEntityId=SYD&stops=direct&toEntityId=BKK":
			var response flightsky.APIResponse
			reader := testhelpers.FileToStruct(t, filepath.Join("testdata", "flightsky-offers.json"), &response)
			data, err := io.ReadAll(reader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			run("Method as expected", func(t *testing.T) {
				assert.Equal(t, http.MethodGet, r.Method)
			})

			run("Token as expected", func(t *testing.T) {
				assert.Equal(t, "TestAPIKEY", r.Header.Get("x-rapidapi-key"))
			})

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		default:
			t.Errorf("unexpected path %s", r.URL.String())
			t.FailNow()
		}

	}))
}

type secretResponse struct {
	Secret secret `json:"secret"`
}

type secret struct {
	SecretValue string `json:"secretValue"`
}

func mockInfisicalServer(t *testing.T, amadeusURL, flightskyURL string) *httptest.Server {
	run := testhelpers.Run(t)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlSecrets := map[string]string{
			"/api/v3/secrets/raw/AMADEUS_BASE_URL?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":          amadeusURL,
			"/api/v3/secrets/raw/FLIGHTS_SKY_BASE_URL?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":      flightskyURL,
			"/api/v3/secrets/raw/JOBSITY_APP_CLIENT_SECRET?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical": "TEST_PWD",
			"/api/v3/secrets/raw/AMADEUS_CLIENT_ID?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":         "testClientId",
			"/api/v3/secrets/raw/AMADEUS_SECRET_ID?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":         "testClientSecret",
			"/api/v3/secrets/raw/JOBSITY_APP_CLIENT_ID?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":     "TEST_USER",
			"/api/v3/secrets/raw/JOBSITY_SECRET_KEY?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":        "TEST_SECRET",
			"/api/v3/secrets/raw/RAPID_API_KEY?environment=dev&include_imports=false&secretPath=%2F&type=shared&workspaceId=testInfisical":             "TestAPIKEY",
		}

		s, ok := urlSecrets[r.URL.String()]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			t.Errorf("unexpected path %s", r.URL.String())
			t.FailNow()
		}

		response := secretResponse{
			Secret: secret{
				SecretValue: s,
			},
		}

		run("Method as expected", func(t *testing.T) {
			assert.Equal(t, http.MethodGet, r.Method)
		})

		data, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)

	}))
}

func TestLoginSuccess(t *testing.T) {
	run := testhelpers.Run(t)

	amadeusServer := mockAmadeusServer(t)
	defer amadeusServer.Close()

	flightskyServer := mockFlightskyServer(t)
	defer flightskyServer.Close()

	testInfisical := mockInfisicalServer(t, amadeusServer.URL, flightskyServer.URL)
	defer testInfisical.Close()

	a := New(
		mockInfisicalClient(testInfisical.URL),
		mockGoogleflightsConfig(),
		func(o *Option) {
			o.DisableRedis = true
		},
	)

	testServer := httptest.NewServer(a.Handler())
	defer testServer.Close()

	// generate a new fresh token, as these have 1 day expiration
	var reqDTO pkg.CrendetialsRequest
	payload := testhelpers.FileToStruct(t, filepath.Join("testdata", "login-request.json"), &reqDTO)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/login", testServer.URL), payload)
	res, err := http.DefaultClient.Do(req)

	run("No http error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("HTTP Status response is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	data, err := io.ReadAll(res.Body)
	run("No read error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	var resDTO pkg.CredentialsResponse
	run("No unmarshal error", func(t *testing.T) {
		assert.NoError(t, json.Unmarshal(data, &resDTO))
	})

	run("Token is present", func(t *testing.T) {
		assert.NotEmpty(t, resDTO.AccessToken)
	})
}

func TestLoginFailWrongCredentials(t *testing.T) {
	run := testhelpers.Run(t)

	amadeusServer := mockAmadeusServer(t)
	defer amadeusServer.Close()

	flightskyServer := mockFlightskyServer(t)
	defer flightskyServer.Close()

	testInfisical := mockInfisicalServer(t, amadeusServer.URL, flightskyServer.URL)
	defer testInfisical.Close()

	a := New(
		mockInfisicalClient(testInfisical.URL),
		mockGoogleflightsConfig(),
		func(o *Option) {
			o.DisableRedis = true
		},
	)

	testServer := httptest.NewServer(a.Handler())
	defer testServer.Close()

	// generate a new fresh token, as these have 1 day expiration
	var reqDTO pkg.CrendetialsRequest
	payload := testhelpers.FileToStruct(t, filepath.Join("testdata", "login-bad-credentials-request.json"), &reqDTO)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/login", testServer.URL), payload)
	res, err := http.DefaultClient.Do(req)

	run("No http error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("HTTP Status response is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
}

func TestGetBestFlightOffersResponse(t *testing.T) {
	run := testhelpers.Run(t)

	amadeusServer := mockAmadeusServer(t)
	defer amadeusServer.Close()

	flightskyServer := mockFlightskyServer(t)
	defer flightskyServer.Close()

	testInfisical := mockInfisicalServer(t, amadeusServer.URL, flightskyServer.URL)
	defer testInfisical.Close()

	a := New(
		mockInfisicalClient(testInfisical.URL),
		mockGoogleflightsConfig(),
		func(o *Option) {
			o.DisableRedis = true
		},
	)

	testServer := httptest.NewServer(a.Handler())
	defer testServer.Close()

	// generate a new fresh token, as these have 1 day expiration
	var reqDTO pkg.CrendetialsRequest
	payload := testhelpers.FileToStruct(t, filepath.Join("testdata", "login-request.json"), &reqDTO)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/login", testServer.URL), payload)
	res, err := http.DefaultClient.Do(req)

	run("No http error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("HTTP Status response is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	data, err := io.ReadAll(res.Body)
	run("No read error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	var resDTO pkg.CredentialsResponse
	run("No unmarshal error", func(t *testing.T) {
		assert.NoError(t, json.Unmarshal(data, &resDTO))
	})

	run("Token is present", func(t *testing.T) {
		assert.NotEmpty(t, resDTO.AccessToken)
	})

	// use the token and now access to flights API
	params := url.Values{}
	params.Add("date", "2025-05-09")
	params.Add("origin", "SYD")
	params.Add("adults", "1")
	params.Add("destination", "BKK")
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/flights/search?%s", testServer.URL, params.Encode()), nil)
	req.Header.Add("Authorization", "Bearer "+resDTO.AccessToken)

	res, err = http.DefaultClient.Do(req)
	run("No error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("HTTP Status response is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	run("Response body is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestGetBestFlightOffersResponseUnAuthorized(t *testing.T) {
	run := testhelpers.Run(t)

	amadeusServer := mockAmadeusServer(t)
	defer amadeusServer.Close()

	flightskyServer := mockFlightskyServer(t)
	defer flightskyServer.Close()

	testInfisical := mockInfisicalServer(t, amadeusServer.URL, flightskyServer.URL)
	defer testInfisical.Close()

	a := New(
		mockInfisicalClient(testInfisical.URL),
		mockGoogleflightsConfig(),
		func(o *Option) {
			o.DisableRedis = true
		},
	)

	testServer := httptest.NewServer(a.Handler())
	defer testServer.Close()

	// use the token and now access to flights API
	params := url.Values{}
	params.Add("date", "2025-05-09")
	params.Add("origin", "SYD")
	params.Add("adults", "1")
	params.Add("destination", "BKK")
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/flights/search?%s", testServer.URL, params.Encode()), nil)
	req.Header.Add("Authorization", "") // no token

	res, err := http.DefaultClient.Do(req)
	run("No error", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("HTTP Status response is as expected", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

}
