package services

import (
	"context"
	"log"
	"math"
	"sort"
	"time"
)

// DataValidator handles validation and reranking of retrieved data
type DataValidator struct {
	ragRetriever  *RAGRetriever
	dataConnector *DataSourceConnector
}

// NewDataValidator creates a new data validator
func NewDataValidator(ragRetriever *RAGRetriever, dataConnector *DataSourceConnector) *DataValidator {
	return &DataValidator{
		ragRetriever:  ragRetriever,
		dataConnector: dataConnector,
	}
}

// ValidationCriteria defines criteria for data validation
type ValidationCriteria struct {
	Budget            float64                `json:"budget"`
	MaxTravelTime     time.Duration          `json:"max_travel_time"`
	RequiredRating    float64                `json:"required_rating"`
	PreferredTypes    []string               `json:"preferred_types"`
	Accessibility     bool                   `json:"accessibility"`
	Preferences       map[string]interface{} `json:"preferences"`
	AvailabilityCheck bool                   `json:"availability_check"`
}

// RankingWeights defines weights for different ranking factors
type RankingWeights struct {
	Rating       float64 `json:"rating"`       // 0.3
	Price        float64 `json:"price"`        // 0.25
	Distance     float64 `json:"distance"`     // 0.2
	Availability float64 `json:"availability"` // 0.15
	UserMatch    float64 `json:"user_match"`   // 0.1
}

// DefaultRankingWeights returns default ranking weights
func DefaultRankingWeights() RankingWeights {
	return RankingWeights{
		Rating:       0.3,
		Price:        0.25,
		Distance:     0.2,
		Availability: 0.15,
		UserMatch:    0.1,
	}
}

// ValidateAndRankAttractions validates and ranks attractions
func (dv *DataValidator) ValidateAndRankAttractions(ctx context.Context, attractions []Attraction, criteria ValidationCriteria, weights RankingWeights) ([]Attraction, error) {
	log.Printf("Validating and ranking %d attractions", len(attractions))

	var validAttractions []Attraction

	// Step 1: Validate attractions
	for _, attraction := range attractions {
		if dv.isAttractionValid(attraction, criteria) {
			validAttractions = append(validAttractions, attraction)
		}
	}

	log.Printf("After validation: %d attractions remain", len(validAttractions))

	// Step 2: Calculate scores and rank
	scoredAttractions := dv.scoreAttractions(validAttractions, criteria, weights)

	// Step 3: Sort by score (descending)
	sort.Slice(scoredAttractions, func(i, j int) bool {
		return scoredAttractions[i].score > scoredAttractions[j].score
	})

	// Extract attractions without scores
	rankedAttractions := make([]Attraction, len(scoredAttractions))
	for i, scored := range scoredAttractions {
		rankedAttractions[i] = scored.attraction
	}

	return rankedAttractions, nil
}

// ValidateAndRankHotels validates and ranks hotels
func (dv *DataValidator) ValidateAndRankHotels(ctx context.Context, hotels []Hotel, criteria ValidationCriteria, weights RankingWeights) ([]Hotel, error) {
	log.Printf("Validating and ranking %d hotels", len(hotels))

	var validHotels []Hotel

	// Step 1: Validate hotels
	for _, hotel := range hotels {
		if dv.isHotelValid(hotel, criteria) {
			validHotels = append(validHotels, hotel)
		}
	}

	log.Printf("After validation: %d hotels remain", len(validHotels))

	// Step 2: Calculate scores and rank
	scoredHotels := dv.scoreHotels(validHotels, criteria, weights)

	// Step 3: Sort by score (descending)
	sort.Slice(scoredHotels, func(i, j int) bool {
		return scoredHotels[i].score > scoredHotels[j].score
	})

	// Extract hotels without scores
	rankedHotels := make([]Hotel, len(scoredHotels))
	for i, scored := range scoredHotels {
		rankedHotels[i] = scored.hotel
	}

	return rankedHotels, nil
}

// ValidateAvailability checks real-time availability
func (dv *DataValidator) ValidateAvailability(ctx context.Context, items interface{}, checkDate time.Time) error {
	// This would integrate with real booking systems
	// For now, mark random items as unavailable for simulation

	switch v := items.(type) {
	case []Attraction:
		for i := range v {
			// Simulate some attractions being unavailable
			if i%5 == 0 {
				v[i].Available = false
			}
		}
	case []Hotel:
		for i := range v {
			// Simulate some hotels being fully booked
			if i%7 == 0 {
				v[i].Available = false
			}
		}
	case []TransportOption:
		for i := range v {
			// Simulate some transport options being unavailable
			if i%4 == 0 {
				v[i].Available = false
			}
		}
	}

	return nil
}

// Internal types for scoring
type scoredAttraction struct {
	attraction Attraction
	score      float64
}

type scoredHotel struct {
	hotel Hotel
	score float64
}

// Validation methods

func (dv *DataValidator) isAttractionValid(attraction Attraction, criteria ValidationCriteria) bool {
	// Check availability
	if criteria.AvailabilityCheck && !attraction.Available {
		return false
	}

	// Check rating requirement
	if attraction.Rating < criteria.RequiredRating {
		return false
	}

	// Check if type matches preferences
	if len(criteria.PreferredTypes) > 0 {
		found := false
		for _, preferredType := range criteria.PreferredTypes {
			if attraction.Type == preferredType {
				found = true
				break
			}
			// Also check tags
			for _, tag := range attraction.Tags {
				if tag == preferredType {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (dv *DataValidator) isHotelValid(hotel Hotel, criteria ValidationCriteria) bool {
	// Check availability
	if criteria.AvailabilityCheck && !hotel.Available {
		return false
	}

	// Check rating requirement
	if hotel.Rating < criteria.RequiredRating {
		return false
	}

	// Check budget (assuming per night, multiply by estimated nights)
	if criteria.Budget > 0 {
		estimatedNights := 7.0 // Default assumption
		totalCost := hotel.PricePerNight * estimatedNights
		if totalCost > criteria.Budget*0.5 { // Hotels shouldn't exceed 50% of total budget
			return false
		}
	}

	return true
}

// Scoring methods

func (dv *DataValidator) scoreAttractions(attractions []Attraction, criteria ValidationCriteria, weights RankingWeights) []scoredAttraction {
	var scored []scoredAttraction

	maxRating := 5.0
	maxPrice := dv.findMaxAttractionPrice(attractions)

	for _, attraction := range attractions {
		score := 0.0

		// Rating score (normalized)
		if weights.Rating > 0 {
			ratingScore := attraction.Rating / maxRating
			score += weights.Rating * ratingScore
		}

		// Price score (lower price = higher score, normalized)
		if weights.Price > 0 && maxPrice > 0 {
			priceScore := 1.0 - (float64(attraction.PriceLevel) / 4.0) // Price level 0-4
			score += weights.Price * priceScore
		}

		// Availability score
		if weights.Availability > 0 {
			availabilityScore := 0.0
			if attraction.Available {
				availabilityScore = 1.0
			}
			score += weights.Availability * availabilityScore
		}

		// User preference match score
		if weights.UserMatch > 0 {
			userMatchScore := dv.calculateAttractionUserMatch(attraction, criteria)
			score += weights.UserMatch * userMatchScore
		}

		scored = append(scored, scoredAttraction{
			attraction: attraction,
			score:      score,
		})
	}

	return scored
}

func (dv *DataValidator) scoreHotels(hotels []Hotel, criteria ValidationCriteria, weights RankingWeights) []scoredHotel {
	var scored []scoredHotel

	maxRating := 5.0
	maxPrice := dv.findMaxHotelPrice(hotels)

	for _, hotel := range hotels {
		score := 0.0

		// Rating score (normalized)
		if weights.Rating > 0 {
			ratingScore := hotel.Rating / maxRating
			score += weights.Rating * ratingScore
		}

		// Price score (lower price = higher score, normalized)
		if weights.Price > 0 && maxPrice > 0 {
			priceScore := 1.0 - (hotel.PricePerNight / maxPrice)
			score += weights.Price * priceScore
		}

		// Availability score
		if weights.Availability > 0 {
			availabilityScore := 0.0
			if hotel.Available {
				availabilityScore = 1.0
			}
			score += weights.Availability * availabilityScore
		}

		// User preference match score
		if weights.UserMatch > 0 {
			userMatchScore := dv.calculateHotelUserMatch(hotel, criteria)
			score += weights.UserMatch * userMatchScore
		}

		scored = append(scored, scoredHotel{
			hotel: hotel,
			score: score,
		})
	}

	return scored
}

// Helper methods

func (dv *DataValidator) findMaxAttractionPrice(attractions []Attraction) float64 {
	max := 0.0
	for _, attraction := range attractions {
		if float64(attraction.PriceLevel) > max {
			max = float64(attraction.PriceLevel)
		}
	}
	return math.Max(max, 1.0) // Avoid division by zero
}

func (dv *DataValidator) findMaxHotelPrice(hotels []Hotel) float64 {
	max := 0.0
	for _, hotel := range hotels {
		if hotel.PricePerNight > max {
			max = hotel.PricePerNight
		}
	}
	return math.Max(max, 1.0) // Avoid division by zero
}

func (dv *DataValidator) calculateAttractionUserMatch(attraction Attraction, criteria ValidationCriteria) float64 {
	score := 0.0
	totalFactors := 0.0

	// Check if attraction type matches preferred types
	if len(criteria.PreferredTypes) > 0 {
		for _, preferredType := range criteria.PreferredTypes {
			if attraction.Type == preferredType {
				score += 1.0
			}
			// Check tags too
			for _, tag := range attraction.Tags {
				if tag == preferredType {
					score += 0.5
					break
				}
			}
		}
		totalFactors += 1.0
	}

	// Check specific preferences
	if criteria.Preferences != nil {
		// Example: cultural preference
		if cultural, ok := criteria.Preferences["cultural"]; ok && cultural == true {
			if attraction.Type == "culture" || attraction.Type == "museum" {
				score += 1.0
			}
		}

		// Example: outdoor preference
		if outdoor, ok := criteria.Preferences["outdoor"]; ok && outdoor == true {
			if attraction.Type == "nature" || attraction.Type == "park" {
				score += 1.0
			}
		}

		totalFactors += 1.0
	}

	if totalFactors == 0 {
		return 0.5 // Neutral score if no preferences
	}

	return math.Min(score/totalFactors, 1.0)
}

func (dv *DataValidator) calculateHotelUserMatch(hotel Hotel, criteria ValidationCriteria) float64 {
	score := 0.5 // Base score

	// Budget alignment
	if criteria.Budget > 0 {
		estimatedNights := 7.0
		hotelBudgetRatio := (hotel.PricePerNight * estimatedNights) / criteria.Budget

		if hotelBudgetRatio <= 0.3 { // Very affordable
			score += 0.3
		} else if hotelBudgetRatio <= 0.5 { // Moderately affordable
			score += 0.1
		} else if hotelBudgetRatio > 0.7 { // Expensive
			score -= 0.2
		}
	}

	// Amenity preferences
	if criteria.Preferences != nil {
		if luxury, ok := criteria.Preferences["luxury"]; ok && luxury == true {
			if hotel.Rating >= 4.0 {
				score += 0.2
			}
		}

		if budget, ok := criteria.Preferences["budget"]; ok && budget == true {
			if hotel.PricePerNight <= 100 {
				score += 0.2
			}
		}
	}

	return math.Min(math.Max(score, 0.0), 1.0)
}

// ApplyBudgetConstraints applies budget constraints to the entire trip context
func (dv *DataValidator) ApplyBudgetConstraints(ctx context.Context, tripContext *TripContext, totalBudget float64) error {
	if totalBudget <= 0 {
		return nil // No budget constraints
	}

	log.Printf("Applying budget constraints: $%.2f", totalBudget)

	// Allocate budget percentages
	hotelBudget := totalBudget * 0.4     // 40% for accommodation
	activityBudget := totalBudget * 0.3  // 30% for activities
	transportBudget := totalBudget * 0.2 // 20% for transportation
	// Remaining 10% for miscellaneous

	// Filter hotels by budget
	var affordableHotels []Hotel
	for _, hotel := range tripContext.Hotels {
		estimatedNights := 7.0 // Default
		totalHotelCost := hotel.PricePerNight * estimatedNights
		if totalHotelCost <= hotelBudget {
			affordableHotels = append(affordableHotels, hotel)
		}
	}
	tripContext.Hotels = affordableHotels

	// Filter attractions by price level and budget
	var affordableAttractions []Attraction
	for _, attraction := range tripContext.Attractions {
		estimatedCost := float64(attraction.PriceLevel) * 25 // Rough estimation
		if estimatedCost <= activityBudget/5 {               // Assuming 5 activities
			affordableAttractions = append(affordableAttractions, attraction)
		}
	}
	tripContext.Attractions = affordableAttractions

	// Filter transportation by budget
	var affordableTransport []TransportOption
	for _, transport := range tripContext.Transportation {
		if transport.Price <= transportBudget {
			affordableTransport = append(affordableTransport, transport)
		}
	}
	tripContext.Transportation = affordableTransport

	log.Printf("After budget filtering: %d hotels, %d attractions, %d transport options",
		len(tripContext.Hotels), len(tripContext.Attractions), len(tripContext.Transportation))

	return nil
}
