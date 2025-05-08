package amadeus

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	infisical "github.com/infisical/go-sdk"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/testhelpers"
	"github.com/rubengp99/golang-flights-challenge/pkg"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2/json"
)

func TestRetrieveFlightOffers(t *testing.T) {
	run := testhelpers.Run(t)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/v1/reference-data/airlines?airlineCodes=TG%2CQF":
			var response APIResponse
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
			var response APIResponse
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

			response := AuthResponse{
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
	defer testServer.Close()

	mockConfigProvide := func(client infisical.InfisicalClientInterface, projectID string) vendors.Config {
		return vendors.Config{
			BaseURL:      testServer.URL,
			ClientID:     "testClientId",
			ClientSecret: "testClientSecret",
		}
	}

	service := NewService(mockConfigProvide, infisical.NewInfisicalClient(context.Background(), infisical.Config{}), "")

	date, _ := time.Parse("2006-01-02", "2025-05-09")

	flights, airlines, err := service.RetrieveFlightOffers(pkg.QueryParams{
		Origin:      "SYD",
		Destination: "BKK",
		Date:        date,
		Adults:      "1",
	})

	run("No errors", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("Airlines as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "airlines.json"), airlines)
	})

	run("Flights as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "offers.json"), flights)
	})
}
