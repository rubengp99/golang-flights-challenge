package googleflights

// Duration represents flight duration in google flights
type Duration struct {
	Raw  int    `json:"raw"` // is represent in minutes, e.g 1H23M is 83 here
	Text string `json:"text"`
}

// AirportInfo represents arrival and departure airport info in google flights
type AirportInfo struct {
	AirportName string `json:"airport_name"`
	AirportCode string `json:"airport_code"`
	Time        string `json:"time"`
}

// Delay represents flight delay info in google flights
type Delay struct {
	Values bool `json:"values"`
	Text   int  `json:"text"`
}

// BagsInfo represents allowed bags info for flights in google flight
type BagsInfo struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

// CarbonEmissionInfo represents information about the carbon emissions for a given airplane
type CarbonEmissionInfo struct {
	DifferencePercent   int `json:"difference_percent"`
	CO2E                int `json:"CO2e"`
	TypicalForThisRoute int `json:"typical_for_this_route"`
	Higher              int `json:"higher"`
}

// Flight represents a breakdown of a flight information
type Flight struct {
	DepartureAirport AirportInfo `json:"departure_airport"`
	ArrivalAirport   AirportInfo `json:"arrival_airport"`
	Duration         Duration    `json:"duration"`
	Airline          string      `json:"airline"`
	AirlineLogo      string      `json:"airline_logo"`
	FlightNumber     string      `json:"flight_number"`
	Aircraft         string      `json:"aircraft"`
	Seat             string      `json:"seat"`
	Legroom          string      `json:"legroom"`
	Extensions       []string    `json:"extensions"`
}

// Itinerary represents a breakdown of a flight itinerary
type Itinerary struct {
	DepartureTime   string             `json:"departure_time"`
	ArrivalTime     string             `json:"arrival_time"`
	Duration        Duration           `json:"duration"`
	Flights         []Flight           `json:"flights"`
	Delay           Delay              `json:"delay"`
	SelfTransfer    bool               `json:"self_transfer"`
	Layovers        []any              `json:"layovers"`
	Bags            BagsInfo           `json:"bags"`
	CarbonEmissions CarbonEmissionInfo `json:"carbon_emissions"`
	Price           float64            `json:"price"`
	Stops           int                `json:"stops"`
	AirlineLogo     string             `json:"airline_logo"`
	BookingToken    string             `json:"booking_token"`
}

// Itineraries represents a list of itineraries in google flights apis
type Itineraries struct {
	TopFlights   []Itinerary `json:"topFlights"`
	OtherFlights []Itinerary `json:"otherFlights"`
}

// Price represents a price for historical record
type Price struct {
	Operation string `json:"operation"`
	Value     int    `json:"value"`
}

// PriceSummary represents
type PriceSummary struct {
	Current float64 `json:"current"`
	Low     []Price `json:"low"`
	Typical []Price `json:"typical"`
	High    []Price `json:"high"`
}

// History represents a timespan in historical price summarry
type History struct {
	Time  int64 `json:"time"`
	Value int   `json:"value"`
}

// PriceHistory represents historical pricing for the given flight search criteria
type PriceHistory struct {
	Summary PriceSummary `json:"summary"`
	History []History    `json:"history"`
}

// FlightOffer represents a list of flights in google flights api
type FlightOffer struct {
	Itineraries  Itineraries  `json:"itineraries"`
	PriceHistory PriceHistory `json:"priceHistory"`
}
