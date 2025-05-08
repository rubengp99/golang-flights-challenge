package flightsky

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
)

func TestRetrieveFlightOffers(t *testing.T) {
	run := testhelpers.Run(t)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/flights/search-one-way?adults=1&currencyCode=USD&departDate=2025-05-09&fromEntityId=SYD&stops=direct&toEntityId=BKK":
			var response APIResponse
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
	defer testServer.Close()

	mockConfigProvide := func(client infisical.InfisicalClientInterface, projectID string) vendors.Config {
		return vendors.Config{
			BaseURL: testServer.URL,
			APIKey:  "TestAPIKEY",
		}
	}

	service := NewService(mockConfigProvide, infisical.NewInfisicalClient(context.Background(), infisical.Config{}), "")

	date, _ := time.Parse("2006-01-02", "2025-05-09")

	flights, err := service.RetrieveFlightOffers(pkg.QueryParams{
		Origin:      "SYD",
		Destination: "BKK",
		Date:        date,
		Adults:      "1",
	})

	run("No errors", func(t *testing.T) {
		assert.NoError(t, err)
	})

	run("Flights as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "offers.json"), flights)
	})
}
