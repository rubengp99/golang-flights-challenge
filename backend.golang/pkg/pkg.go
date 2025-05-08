package pkg

import "time"

// QueryParams represents our API available query parameters for decoding
type QueryParams struct {
	Origin      string    `json:"origin"`
	Adults      string    `json:"adults"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
}

// Location represents flight location and time
type Location struct {
	Timestamp time.Time `json:"timestamp"`
	IataCode  string    `json:"iataCode"`
}

// Amount represents flight pricing
type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// FlightOffer represents flight offer breakdown
type FlightOffer struct {
	Airline           string   `json:"airline"`
	FlightNumber      string   `json:"flightNumber"`
	Arrival           Location `json:"arrival"`
	Departure         Location `json:"departure"`
	DurationInMinutes float64  `json:"durationInMinutes"`
	Layovers          int      `json:"layovers"`
	Price             Amount   `json:"price"`
}

// GetBestFlightOffersResponse is the response for best flights API
type GetBestFlightOffersResponse struct {
	Cheapest []FlightOffer `json:"cheapest"`
	Fastest  []FlightOffer `json:"fastest"`
}

// CrendetialsRequest represents app credentials
type CrendetialsRequest struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"ClientSecret"`
}

// CredentialsResponse represents login credentials response
type CredentialsResponse struct {
	AccessToken string `json:"accessToken"`
	ExpIn       int64  `json:"expIn"`
}
