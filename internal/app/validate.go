package app

import (
	"fmt"

	"github.com/rubengp99/golang-flights-challenge/pkg"
)

func validateBestFlightsParamsRequest(req pkg.QueryParams) error {
	if req.Origin == "" {
		return fmt.Errorf("ORIGIN should not be empty")
	}

	if req.Destination == "" {
		return fmt.Errorf("DESTINATION should not be empty")
	}

	if req.Adults == "" {
		return fmt.Errorf("ADULTS should not be empty")
	}

	if req.Date.Unix() <= 0 {
		return fmt.Errorf("DATE should not be empty")
	}

	return nil
}

func validateCrendetialsRequest(req pkg.CrendetialsRequest) error {
	if req.ClientID == "" {
		return fmt.Errorf("USERNAME should not be empty")
	}

	if req.ClientSecret == "" {
		return fmt.Errorf("USERNAME should not be empty")
	}

	return nil
}
