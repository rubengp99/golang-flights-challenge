package app

import (
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/internal/workflow"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

func getQueryParams(value interface{}, r *http.Request) error {
	decoder := schema.NewDecoder()
	// decoder lookup for values on the json tag, instead of the default schema tag
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(value, r.URL.Query()); err != nil {
		return errors.Wrap(err, "handler - failed to decode query params")
	}

	return nil
}

// RetrieveBestFlightsHandler handles best flights lookup
func RetrieveBestFlightsHandler(googleflightService googleflights.Service,
	amadeusService amadeus.Service,
	flightskyService flightsky.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params pkg.QueryParams
		if err := getQueryParams(&params, r); err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}

		wf := workflow.RetrieveBestFlights(googleflightService, amadeusService, flightskyService)
		res, err := wf(params)
		if err != nil {
			serveResponse(newError(err.Error()), http.StatusInternalServerError, w)
			return
		}
		serveResponse(res, http.StatusOK, w)
	})
}
