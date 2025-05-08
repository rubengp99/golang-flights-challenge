package mapping_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/rubengp99/golang-flights-challenge/internal/mapping"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/testhelpers"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

func TestAmadeusToPkgFlights(t *testing.T) {
	var amadeusFlights amadeus.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-offers.json"), &amadeusFlights)

	var amadeusOffers []amadeus.FlightOffer
	if err := json.Unmarshal(amadeusFlights.Data, &amadeusOffers); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var amadeusAirlinesResp amadeus.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-airlines.json"), &amadeusAirlinesResp)

	var amadeusAirlines []amadeus.Airline
	if err := json.Unmarshal(amadeusAirlinesResp.Data, &amadeusAirlines); err != nil {
		t.Error(err)
		t.FailNow()
	}

	actual := mapping.AmadeusToPkgFlights(make(chan error), amadeusOffers, amadeusAirlines)

	run := testhelpers.Run(t)

	run("Parsed pkg response is as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "amadeus-offers-pkg-expected.json"), actual)
	})
}

func TestGoogleflightsToPkgFlights(t *testing.T) {
	var googleflightsFlights googleflights.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "googleflights-offers.json"), &googleflightsFlights)

	actual := mapping.GoogleflightsToPkgFlights(make(chan error), googleflightsFlights.FlightOffer)

	run := testhelpers.Run(t)

	run("Parsed pkg response is as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "googleflights-offers-pkg-expected.json"), actual)
	})
}

func TestFlightskyToPkgFlights(t *testing.T) {
	var flightskyFlights flightsky.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "flightsky-offers.json"), &flightskyFlights)

	var flightskyOffers flightsky.FlightOffer
	if err := json.Unmarshal(flightskyFlights.Data, &flightskyOffers); err != nil {
		t.Error(err)
		t.FailNow()
	}

	actual := mapping.FlightskyToPkgFlights(make(chan error), flightskyOffers)

	run := testhelpers.Run(t)

	run("Parsed pkg response is as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "flightsky-offers-pkg-expected.json"), actual)
	})
}

func TestNewBestFlightsOffersResponse(t *testing.T) {
	var amadeusFlights amadeus.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-offers.json"), &amadeusFlights)

	var amadeusOffers []amadeus.FlightOffer
	if err := json.Unmarshal(amadeusFlights.Data, &amadeusOffers); err != nil {
		t.Error(err)
		t.FailNow()
	}

	var amadeusAirlinesResp amadeus.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "amadeus-airlines.json"), &amadeusAirlinesResp)

	var amadeusAirlines []amadeus.Airline
	if err := json.Unmarshal(amadeusAirlinesResp.Data, &amadeusAirlines); err != nil {
		t.Error(err)
		t.FailNow()
	}

	amadeusList := mapping.AmadeusToPkgFlights(make(chan error), amadeusOffers, amadeusAirlines)

	var flightskyFlights flightsky.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "flightsky-offers.json"), &flightskyFlights)

	var flightskyOffers flightsky.FlightOffer
	if err := json.Unmarshal(flightskyFlights.Data, &flightskyOffers); err != nil {
		t.Error(err)
		t.FailNow()
	}

	flightskyList := mapping.FlightskyToPkgFlights(make(chan error), flightskyOffers)

	var googleflightsFlights googleflights.APIResponse
	testhelpers.FileToStruct(t, filepath.Join("testdata", "googleflights-offers.json"), &googleflightsFlights)

	googleflightsList := mapping.GoogleflightsToPkgFlights(make(chan error), googleflightsFlights.FlightOffer)

	wholelist := []pkg.FlightOffer{}
	wholelist = append(wholelist, amadeusList...)
	wholelist = append(wholelist, googleflightsList...)
	wholelist = append(wholelist, flightskyList...)

	actual := mapping.NewBestFlightsOffersResponse(wholelist...)

	run := testhelpers.Run(t)

	run("Parsed pkg response is as expected", func(t *testing.T) {
		testhelpers.AssertJSONEquals(t, filepath.Join("testdata", "best-flight-offers-pkg-expected.json"), actual)
	})
}
