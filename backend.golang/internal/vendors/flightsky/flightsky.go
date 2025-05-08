package flightsky

// Context represents metadata about the API response
type Context struct {
	Status       string `json:"status"`
	SessionID    string `json:"sessionId"`
	TotalResults int    `json:"totalResults"`
}

// Price price represents itinerary pricing
type Price struct {
	Raw             float64 `json:"raw"`
	Formatted       string  `json:"formatted"`
	PricingOptionID string  `json:"pricingOptionId"`
}

// Location represents the info about arrival, destiny
type Location struct {
	ID            string `json:"id"`
	EntityID      string `json:"entityId"`
	Name          string `json:"name"`
	DisplayCode   string `json:"displayCode"`
	City          string `json:"city"`
	Country       string `json:"country"`
	IsHighlighted bool   `json:"isHighlighted"`
}

// MarketingInfo represents metadata about the service providers
type MarketingInfo struct {
	ID          int    `json:"id"`
	AlternateID string `json:"alternateId"`
	Name        string `json:"name"`
}

// Carriers represents carrier or service provider information
type Carriers struct {
	Marketing     []MarketingInfo `json:"marketing"`
	OperationType string          `json:"operationType"`
}

// SegmentLocation represents location information about a trip segment
type SegmentLocation struct {
	FlightPlaceID string `json:"flightPlaceId"`
	DisplayCode   string `json:"displayCode"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Country       string `json:"country"`
}

// MarketingCarrrier represents marketing info about our flight carrier
type MarketingCarrier struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	AlternateID string `json:"alternateId"`
	AllianceID  int    `json:"allianceId"`
	DisplayCode string `json:"displayCode"`
}

// OperatingCarrier reprersents information about the carrier operating our flight
type OperatingCarrier struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	AlternateID string `json:"alternateId"`
	AllianceID  int    `json:"allianceId"`
	DisplayCode string `json:"displayCode"`
}

// Segment represents information about the flight stops
type Segment struct {
	ID                string           `json:"id"`
	Origin            SegmentLocation  `json:"origin"`
	Destination       SegmentLocation  `json:"destination"`
	Departure         string           `json:"departure"`
	Arrival           string           `json:"arrival"`
	DurationInMinutes int              `json:"durationInMinutes"`
	FlightNumber      string           `json:"flightNumber"`
	MarketingCarrier  MarketingCarrier `json:"marketingCarrier"`
	OperatingCarrier  OperatingCarrier `json:"operatingCarrier"`
}

// Flight represents the main information about the flight offer
type Flight struct {
	ID                string    `json:"id"`
	Origin            Location  `json:"origin"`
	Destination       Location  `json:"destination"`
	DurationInMinutes int       `json:"durationInMinutes"`
	StopCount         int       `json:"stopCount"`
	IsSmallestStops   bool      `json:"isSmallestStops"`
	Departure         string    `json:"departure"`
	Arrival           string    `json:"arrival"`
	TimeDeltaInDays   int       `json:"timeDeltaInDays"`
	Carriers          Carriers  `json:"carriers"`
	Segments          []Segment `json:"segments"`
}

// FarePolicy represents information about flight changes made by customer allowed
type FarePolicy struct {
	IsChangeAllowed       bool `json:"isChangeAllowed"`
	IsPartiallyChangeable bool `json:"isPartiallyChangeable"`
	IsCancellationAllowed bool `json:"isCancellationAllowed"`
	IsPartiallyRefundable bool `json:"isPartiallyRefundable"`
}

// Itinerary represents information about the flight itinerary listing
type Itinerary struct {
	ID                      string     `json:"id"`
	Price                   Price      `json:"price"`
	Legs                    []Flight   `json:"legs"`
	IsSelfTransfer          bool       `json:"isSelfTransfer"`
	IsProtectedSelfTransfer bool       `json:"isProtectedSelfTransfer"`
	FarePolicy              FarePolicy `json:"farePolicy"`
	FareAttributes          any        `json:"fareAttributes"`
	Tags                    []string   `json:"tags,omitempty"`
	IsMashUp                bool       `json:"isMashUp"`
	HasFlexibleOptions      bool       `json:"hasFlexibleOptions"`
	Score                   float64    `json:"score"`
}

// Duration represents flight duration for arrival
type Duration struct {
	Min          int `json:"min"`
	Max          int `json:"max"`
	MultiCityMin int `json:"multiCityMin"`
	MultiCityMax int `json:"multiCityMax"`
}

// Airport represents information about provider airport
type Airport struct {
	ID       string `json:"id"`
	EntityID string `json:"entityId"`
	Name     string `json:"name"`
}

// AirportLocation represents information about
type AirportLocation struct {
	City     string    `json:"city"`
	Airports []Airport `json:"airports"`
}

// Carrier represents information about flight carrier
type Carrier struct {
	ID          int    `json:"id"`
	AlternateID string `json:"alternateId"`
	Name        string `json:"name"`
}

// StopPrice represents information about the pricing per stop
type StopPrice struct {
	IsPresent      bool   `json:"isPresent"`
	FormattedPrice string `json:"formattedPrice"`
}

// StopPriceSummary represents a breakdown of stop pricings
type StopPriceSummary struct {
	Direct    StopPrice `json:"direct"`
	One       StopPrice `json:"one"`
	TwoOrMore StopPrice `json:"twoOrMore"`
}

// FilterStats represents metadata about the current flight offer
type FilterStats struct {
	Duration   Duration          `json:"duration"`
	Airports   []AirportLocation `json:"airports"`
	Carriers   []Carrier         `json:"carriers"`
	StopPrices StopPriceSummary  `json:"stopPrices"`
}

// FlightOffer represents a breakdown of all available flight offers per search criteria
type FlightOffer struct {
	Context             Context     `json:"context"`
	Itineraries         []Itinerary `json:"itineraries"`
	Messages            []any       `json:"messages"`
	FilterStats         FilterStats `json:"filterStats"`
	FlightsSessionID    string      `json:"flightsSessionId"`
	DestinationImageURL string      `json:"destinationImageUrl"`
	Token               string      `json:"token"`
}
