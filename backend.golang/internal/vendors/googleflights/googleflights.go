package googleflights

// Duration represents flight duration in google flights
type Duration struct {
	Raw  int    `json:"raw"` // is represent in minutes, e.g 1H23M is 83 here
	Text string `json:"text"`
}

// AirportInfo represents arrival and departure airport info in google flights
type AirportInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Time string `json:"time"`
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
	ThisFlight          int `json:"this_flight"`
	TypicalForThisRoute int `json:"typical_for_this_route"`
	DifferencePercent   int `json:"difference_percent"`
}

// Flight represents a breakdown of a flight information
type Flight struct {
	DepartureAirport AirportInfo `json:"departure_airport"`
	ArrivalAirport   AirportInfo `json:"arrival_airport"`
	Duration         int         `json:"duration"`
	Airplane         string      `json:"airplane"`
	Airline          string      `json:"airline"`
	AirlineLogo      string      `json:"airline_logo"`
	TravelClass      string      `json:"travel_class"`
	FlightNumber     string      `json:"flight_number"`
	Legroom          string      `json:"legroom"`
	Extensions       []string    `json:"extensions"`
	Overnight        bool        `json:"overnight,omitempty"`
}

// Itinerary represents a breakdown of a flight itinerary
type Itinerary struct {
	Flights         []Flight           `json:"flights"`
	Layovers        []any              `json:"layovers"`
	TotalDuration   int                `json:"total_duration"`
	CarbonEmissions CarbonEmissionInfo `json:"carbon_emissions"`
	Price           float64            `json:"price"`
	Type            string             `json:"type"`
	AirlineLogo     string             `json:"airline_logo"`
	BookingToken    string             `json:"booking_token"`
}

// Price represents a price for historical record
type Price struct {
	Operation string `json:"operation"`
	Value     int    `json:"value"`
}

// PriceSummary represents
type PriceSummary struct {
	LowestPrice       int     `json:"lowest_price"`
	PriceLevel        string  `json:"price_level"`
	TypicalPriceRange []int   `json:"typical_price_range"`
	PriceHistory      [][]int `json:"price_history"`
}

// History represents a timespan in historical price summarry
type History struct {
	Time  int64 `json:"time"`
	Value int   `json:"value"`
}

// FlightOffer represents a list of flights in google flights api
type FlightOffer struct {
	BestFlights   []Itinerary  `json:"best_flights"`
	OtherFlights  []Itinerary  `json:"other_flights"`
	PriceInsights PriceSummary `json:"price_insights"`
	Airports      []any        `json:"airports"`
}
