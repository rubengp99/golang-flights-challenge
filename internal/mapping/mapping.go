package mapping

import (
	"sort"
	"strconv"
	"time"

	"github.com/rubengp99/golang-flights-challenge/internal/vendors/amadeus"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/flightsky"
	"github.com/rubengp99/golang-flights-challenge/internal/vendors/googleflights"
	"github.com/rubengp99/golang-flights-challenge/pkg"
)

const (
	ISO8601TimeFormat = "2006-01-02T15:04:05"
)

// GoogleflightsToPkgFlights maps google flights response format to a generic pkg flight offer one
func GoogleflightsToPkgFlights(c chan error, gflights googleflights.FlightOffer) []pkg.FlightOffer {
	results := []pkg.FlightOffer{}

	itineraries := []googleflights.Itinerary{}
	itineraries = append(itineraries, gflights.Itineraries.TopFlights...)
	itineraries = append(itineraries, gflights.Itineraries.OtherFlights...)

	for _, itinerary := range itineraries {
		for _, flight := range itinerary.Flights {
			arrivalTime, err := time.Parse(ISO8601TimeFormat, flight.ArrivalAirport.Time)
			if err != nil {
				c <- err
				return []pkg.FlightOffer{}
			}

			departureTime, err := time.Parse(ISO8601TimeFormat, flight.DepartureAirport.Time)
			if err != nil {
				c <- err
				return []pkg.FlightOffer{}
			}

			mapped := pkg.FlightOffer{
				Airline:      flight.Airline,
				FlightNumber: flight.FlightNumber,
				Arrival: pkg.Location{
					Timestamp: arrivalTime,
					IataCode:  flight.DepartureAirport.AirportCode,
				},
				Departure: pkg.Location{
					Timestamp: departureTime,
					IataCode:  flight.DepartureAirport.AirportCode,
				},
				DurationInMinutes: float64(flight.Duration.Raw),
				Price: pkg.Amount{
					Value:    itinerary.Price,
					Currency: "USD",
				},
				Layovers: len(itinerary.Layovers),
			}

			results = append(results, mapped)
		}
	}

	return results
}

// AmadeusToPkgFlights maps Amadeus flight offers to a generic pkg one
func AmadeusToPkgFlights(c chan error, aflights []amadeus.FlightOffer, airlines []amadeus.Airline) []pkg.FlightOffer {
	results := []pkg.FlightOffer{}
	mapAirlines := map[string]amadeus.Airline{}
	for _, a := range airlines {
		mapAirlines[a.IataCode] = a
	}

	for _, offer := range aflights {
		for _, flight := range offer.Itineraries {
			if len(flight.Segments) == 0 {
				// invalid data, ignore
				continue
			}

			// the first segment represents the departure time
			departureTime, err := time.Parse(ISO8601TimeFormat, flight.Segments[0].Departure.At)
			if err != nil {
				c <- err
				return []pkg.FlightOffer{}
			}

			length := len(flight.Segments) - 1

			// the last segment, even if it contains a single segment, represents the arrival time
			arrivalTime, err := time.Parse(ISO8601TimeFormat, flight.Segments[length].Arrival.At)
			if err != nil {
				c <- err
				return []pkg.FlightOffer{}
			}

			price, err := strconv.ParseFloat(offer.Price.Total, 64)
			if err != nil {
				c <- err
				return []pkg.FlightOffer{}
			}

			airlineName := ""
			if len(offer.ValidatingAirlineCodes) > 0 {
				code := offer.ValidatingAirlineCodes[0]
				airlineName = mapAirlines[code].BusinessName
			}

			mapped := pkg.FlightOffer{
				Airline:      airlineName,
				FlightNumber: flight.Segments[0].Number,
				Arrival: pkg.Location{
					Timestamp: arrivalTime,
					IataCode:  flight.Segments[length].Arrival.IataCode,
				},
				Departure: pkg.Location{
					Timestamp: departureTime,
					IataCode:  flight.Segments[0].Departure.IataCode,
				},
				DurationInMinutes: arrivalTime.Sub(departureTime).Minutes(),
				Price: pkg.Amount{
					Value:    price,
					Currency: "USD",
				},
				Layovers: len(flight.Segments),
			}

			results = append(results, mapped)
		}
	}

	return results
}

// FlightskyToPkgFlights maps flightsky flights format to a generic pkg one
func FlightskyToPkgFlights(c chan error, fsflights flightsky.FlightOffer) []pkg.FlightOffer {
	results := []pkg.FlightOffer{}

	for _, flight := range fsflights.Itineraries {
		if len(flight.Legs) == 0 {
			// invalid data, ignore
			continue
		}

		// the first segment represents the departure time
		departureTime, err := time.Parse(ISO8601TimeFormat, flight.Legs[0].Departure)
		if err != nil {
			c <- err
			return []pkg.FlightOffer{}
		}

		length := len(flight.Legs) - 1

		// the last segment, even if it contains a single segment, represents the arrival time
		arrivalTime, err := time.Parse(ISO8601TimeFormat, flight.Legs[length].Arrival)
		if err != nil {
			c <- err
			return []pkg.FlightOffer{}
		}

		airlineName := ""
		if len(flight.Legs[0].Segments) > 0 {
			airlineName = flight.Legs[0].Segments[0].MarketingCarrier.Name
		} else {
			// invalid data, ignore
			continue
		}

		mapped := pkg.FlightOffer{
			Airline:      airlineName,
			FlightNumber: flight.Legs[0].Segments[0].FlightNumber,
			Arrival: pkg.Location{
				Timestamp: arrivalTime,
				IataCode:  flight.Legs[0].Destination.ID,
			},
			Departure: pkg.Location{
				Timestamp: departureTime,
				IataCode:  flight.Legs[0].Origin.ID,
			},
			DurationInMinutes: arrivalTime.Sub(departureTime).Minutes(),
			Price: pkg.Amount{
				Value:    flight.Price.Raw,
				Currency: "USD",
			},
			Layovers: len(flight.Legs),
		}

		results = append(results, mapped)
	}

	return results
}

func NewBestFlightsOffersResponse(flights ...pkg.FlightOffer) pkg.GetBestFlightOffersResponse {
	response := pkg.GetBestFlightOffersResponse{
		Cheapest: flights,
		Fastest:  flights,
	}

	// sort fastests
	sort.SliceStable(response.Fastest, func(i, j int) bool {
		return response.Fastest[i].DurationInMinutes < response.Fastest[j].DurationInMinutes
	})

	// sort cheapest
	sort.SliceStable(response.Cheapest, func(i, j int) bool {
		return response.Cheapest[i].Price.Value < response.Cheapest[j].Price.Value
	})

	return response
}
