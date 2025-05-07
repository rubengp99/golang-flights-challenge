package pkg

import "time"

// QueryParams represents our API available query parameters for decoding
type QueryParams struct {
	Origin      string    `json:"origin"`
	Adults      string    `json:"adults"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
}

type Location struct {
	Timestamp time.Time `json:"timestamp"`
	IataCode  string    `json:"iataCode"`
}

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type FlightOffer struct {
	Airline           string   `json:"airline"`
	FlightNumber      string   `json:"flightNumber"`
	Arrival           Location `json:"arrival"`
	Departure         Location `json:"departure"`
	DurationInMinutes float64  `json:"durationInMinutes"`
	Layovers          int      `json:"layovers"`
	Price             Amount   `json:"price"`
}

type GetBestFlightOffersResponse struct {
	Cheapest []FlightOffer `json:"cheapest"`
	Fastest  []FlightOffer `json:"fastest"`
}
