package pkg

import "time"

// QueryParams represents our API available query parameters for decoding
type QueryParams struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Date        time.Time `json:"date"`
}

type FlightOffer struct{}
