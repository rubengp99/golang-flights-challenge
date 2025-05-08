package workflow

import (
	"log"
	"sync"

	"github.com/rubengp99/golang-flights-challenge/internal/mapping"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

type RetrieveBestFlightsFunc func(params pkg.QueryParams) (pkg.GetBestFlightOffersResponse, error)

func RetrieveBestFlights(googleflightService googleflights.Service,
	amadeusService amadeus.Service,
	flightskyService flightsky.Service) RetrieveBestFlightsFunc {
	return func(params pkg.QueryParams) (pkg.GetBestFlightOffersResponse, error) {
		var (
			errors       = make(chan error)
			wgdone       = make(chan bool)
			wg           sync.WaitGroup
			flightOffers = []pkg.FlightOffer{}
		)

		// retrieve all secrets from infisical
		retrieveFlightRequests := []func(channel chan error){
			func(channel chan error) {
				flights, err := googleflightService.RetrieveFlightOffers(params)
				flightOffers = append(flightOffers, mapping.GoogleflightsToPkgFlights(errors, flights)...)
				if err != nil {
					channel <- err
				}
				log.Printf("found %v flights with google flights", len(flightOffers))
				wg.Done()
			},
			func(channel chan error) {
				flights, airlines, err := amadeusService.RetrieveFlightOffers(params)
				flightOffers = append(flightOffers, mapping.AmadeusToPkgFlights(errors, flights, airlines)...)
				if err != nil {
					channel <- err
				}
				log.Printf("found %v flights with amadeus", len(flightOffers))
				wg.Done()
			},
			func(channel chan error) {
				flights, err := flightskyService.RetrieveFlightOffers(params)
				flightOffers = append(flightOffers, mapping.FlightskyToPkgFlights(errors, flights)...)
				if err != nil {
					channel <- err
				}
				log.Printf("found %v flights with flightsky", len(flightOffers))
				wg.Done()
			},
		}

		wg.Add(len(retrieveFlightRequests))
		go func() {
			wg.Wait()
			close(wgdone)
		}()

		for _, f := range retrieveFlightRequests {
			go f(errors)
		}

		select {
		case <-wgdone:
			break
		case err := <-errors:
			close(errors)
			return pkg.GetBestFlightOffersResponse{}, err
		}

		return mapping.NewBestFlightsOffersResponse(flightOffers...), nil
	}
}
