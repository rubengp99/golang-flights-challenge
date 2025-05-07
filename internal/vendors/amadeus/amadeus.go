package amadeus

// FlightOffer represents a flight offer in Amadeus API format
type FlightOffer struct {
	Airline                  Airline            `json:"airline,omitempty"`
	Type                     string             `json:"type"`
	ID                       string             `json:"id"`
	Source                   string             `json:"source"`
	InstantTicketingRequired bool               `json:"instantTicketingRequired"`
	NonHomogeneous           bool               `json:"nonHomogeneous"`
	OneWay                   bool               `json:"oneWay"`
	IsUpsellOffer            bool               `json:"isUpsellOffer"`
	LastTicketingDate        string             `json:"lastTicketingDate"`
	LastTicketingDateTime    string             `json:"lastTicketingDateTime"`
	NumberOfBookableSeats    int                `json:"numberOfBookableSeats"`
	Itineraries              []Itinerary        `json:"itineraries"`
	Price                    Price              `json:"price"`
	PricingOptions           PriceOptions       `json:"pricingOptions"`
	ValidatingAirlineCodes   []string           `json:"validatingAirlineCodes"`
	TravelerPricings         []TravelerPricings `json:"travelerPricings"`
}

// Itinerary represents itinerary information for a given flight offer
type Itinerary struct {
	Duration string    `json:"duration"`
	Segments []Segment `json:"segments"`
}

// Segment represents a segment within a given itinerary for a flight offer
type Segment struct {
	Departure       Location         `json:"departure"`
	Arrival         Location         `json:"arrival"`
	CarrierCode     string           `json:"carrierCode"`
	Number          string           `json:"number"`
	Aircraft        SegmentAircraft  `json:"aircraft"`
	Operating       SegmentOperating `json:"operating"`
	Duration        string           `json:"duration"`
	ID              string           `json:"id"`
	NumberOfStops   int              `json:"numberOfStops"`
	BlacklistedInEU bool             `json:"blacklistedInEU"`
}

// SegmentOperating represents the info for a carrier operating a segment in a given itinerary
type SegmentOperating struct {
	CarrierCode string `json:"carrierCode"`
}

// SegmentAircraft represents the aircraft info for a segment in a given itinerary
type SegmentAircraft struct {
	Code string `json:"code"`
}

// SegmentDeparture represents the departure/arrival info for a segment in a given itinerary
type Location struct {
	IataCode string `json:"iataCode"`
	Terminal string `json:"terminal"`
	At       string `json:"at"`
}

// Amount represents a basic pricing info
type Amount struct {
	Value string `json:"amount"`
	Type  string `json:"type"`
}

// Price presents pricing information about a flight offer
type Price struct {
	Currency           string   `json:"currency"`
	Total              string   `json:"total"`
	Base               string   `json:"base"`
	Fees               []Amount `json:"fees"`
	GrandTotal         string   `json:"grandTotal"`
	AdditionalServices []Amount `json:"additionalServices"`
}

// PriceOptions represents additional and options regarding a flight offer
type PriceOptions struct {
	FareType                []string `json:"fareType"`
	IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
}

// TravelerPricingsPrice represents pricings over traveler fee
type TravelerPricingsPrice struct {
	Currency string `json:"currency"`
	Total    string `json:"total"`
	Base     string `json:"base"`
}

// FareDetailsBySegment represents fee breakdown per segment in a given itinerary
type FareDetailsBySegment struct {
	SegmentID           string              `json:"segmentId"`
	Cabin               string              `json:"cabin"`
	FareBasis           string              `json:"fareBasis"`
	BrandedFare         string              `json:"brandedFare"`
	BrandedFareLabel    string              `json:"brandedFareLabel"`
	Class               string              `json:"class"`
	IncludedCheckedBags IncludedCheckedBags `json:"includedCheckedBags"`
	IncludedCabinBags   IncludedCabinBags   `json:"includedCabinBags"`
	Amenities           []Amenity           `json:"amenities"`
}

// IncludedCheckedBags represents the quantity of bags a traveler can check in their flight
type IncludedCheckedBags struct {
	Quantity int `json:"quantity"`
}

// IncludedCabinBags represents the quantity of hand bags a traveler can check in their flight
type IncludedCabinBags struct {
	Quantity int `json:"quantity"`
}

// Amenity represents information about the Amenities available for a given flight offer
type Amenity struct {
	Description     string          `json:"description"`
	IsChargeable    bool            `json:"isChargeable"`
	AmenityType     string          `json:"amenityType"`
	AmenityProvider AmenityProvider `json:"amenityProvider"`
}

// AmenityProvider represents information about who is providing amenities for a given flight offerr
type AmenityProvider struct {
	Name string `json:"name"`
}

// TravelerPricings represents a breakdown of fees, pricings and details about the traveler spending
type TravelerPricings struct {
	TravelerID           string                 `json:"travelerId"`
	FareOption           string                 `json:"fareOption"`
	TravelerType         string                 `json:"travelerType"`
	Price                TravelerPricingsPrice  `json:"price"`
	FareDetailsBySegment []FareDetailsBySegment `json:"fareDetailsBySegment"`
}

// Airline represents airline information for a given flight offer
type Airline struct {
	Type         string `json:"type"`
	IataCode     string `json:"iataCode"`
	IcaoCode     string `json:"icaoCode"`
	BusinessName string `json:"businessName"`
}
