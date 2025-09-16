package services

import (
	"context"
	"log"
	"math"
	"strings"
	"time"
)

// TravelCostPredictor handles lightweight travel cost prediction
type TravelCostPredictor struct {
	// Historical data for cost prediction
	costFactors map[string]CostFactors
	seasonality map[string]SeasonalFactors
}

// CostFactors represents cost factors for different destinations
type CostFactors struct {
	BaseAccommodationCost float64 // Base cost per night
	BaseFoodCost          float64 // Base cost per day
	BaseTransportCost     float64 // Base cost per day
	BaseActivityCost      float64 // Base cost per day
	CostOfLivingIndex     float64 // Relative to global average
	PopularityMultiplier  float64 // Tourism demand factor
}

// SeasonalFactors represents seasonal cost variations
type SeasonalFactors struct {
	PeakSeasonMultiplier float64 // June-August
	MidSeasonMultiplier  float64 // April-May, Sept-Oct
	OffSeasonMultiplier  float64 // Nov-March
	HolidayMultiplier    float64 // Major holidays
	WeekendMultiplier    float64 // Weekend pricing
}

// TravelCostPrediction represents predicted costs
type TravelCostPrediction struct {
	TotalEstimatedCost float64            `json:"total_estimated_cost"`
	CostBreakdown      map[string]float64 `json:"cost_breakdown"`
	ConfidenceLevel    float64            `json:"confidence_level"`
	CostRange          CostRange          `json:"cost_range"`
	Recommendations    []string           `json:"recommendations"`
	SeasonalAdvice     string             `json:"seasonal_advice"`
}

// CostRange represents the range of possible costs
type CostRange struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

// NewTravelCostPredictor creates a new travel cost predictor
func NewTravelCostPredictor() *TravelCostPredictor {
	return &TravelCostPredictor{
		costFactors: initializeCostFactors(),
		seasonality: initializeSeasonalFactors(),
	}
}

// PredictTravelCost predicts travel costs using lightweight ML model
func (tcp *TravelCostPredictor) PredictTravelCost(ctx context.Context, req CostPredictionRequest) (*TravelCostPrediction, error) {
	// Get cost factors for destination
	factors, exists := tcp.costFactors[strings.ToLower(req.Destination)]
	if !exists {
		// Use global average for unknown destinations
		factors = tcp.getGlobalAverageCostFactors()
		log.Printf("Using global average cost factors for %s", req.Destination)
	}

	// Calculate base costs
	accommodationCost := factors.BaseAccommodationCost * float64(req.Duration) * factors.CostOfLivingIndex
	foodCost := factors.BaseFoodCost * float64(req.Duration) * factors.CostOfLivingIndex
	transportCost := factors.BaseTransportCost * float64(req.Duration)
	activityCost := factors.BaseActivityCost * float64(req.Duration)

	// Apply traveler count multiplier
	travelerMultiplier := tcp.calculateTravelerMultiplier(req.Travelers)
	accommodationCost *= travelerMultiplier
	foodCost *= float64(req.Travelers)
	transportCost *= travelerMultiplier
	activityCost *= float64(req.Travelers)

	// Apply seasonal adjustments
	seasonalMultiplier := tcp.calculateSeasonalMultiplier(req.TravelDate, req.Destination)
	accommodationCost *= seasonalMultiplier
	foodCost *= seasonalMultiplier

	// Apply budget preference adjustments
	budgetMultiplier := tcp.calculateBudgetMultiplier(req.BudgetPreference)
	accommodationCost *= budgetMultiplier
	foodCost *= budgetMultiplier
	activityCost *= budgetMultiplier

	// Calculate total cost
	totalCost := accommodationCost + foodCost + transportCost + activityCost

	// Calculate confidence level based on data availability
	confidence := tcp.calculateConfidenceLevel(req.Destination, factors)

	// Generate cost range (±20% for uncertainty)
	uncertainty := 0.2
	minCost := totalCost * (1 - uncertainty)
	maxCost := totalCost * (1 + uncertainty)

	// Generate recommendations
	recommendations := tcp.generateCostRecommendations(req, factors, totalCost)

	// Generate seasonal advice
	seasonalAdvice := tcp.generateSeasonalAdvice(req.TravelDate, req.Destination)

	return &TravelCostPrediction{
		TotalEstimatedCost: math.Round(totalCost*100) / 100,
		CostBreakdown: map[string]float64{
			"accommodation": math.Round(accommodationCost*100) / 100,
			"food":          math.Round(foodCost*100) / 100,
			"transport":     math.Round(transportCost*100) / 100,
			"activities":    math.Round(activityCost*100) / 100,
		},
		ConfidenceLevel: confidence,
		CostRange: CostRange{
			Minimum: math.Round(minCost*100) / 100,
			Maximum: math.Round(maxCost*100) / 100,
		},
		Recommendations: recommendations,
		SeasonalAdvice:  seasonalAdvice,
	}, nil
}

// CostPredictionRequest represents a request for cost prediction
type CostPredictionRequest struct {
	Destination      string    `json:"destination"`
	TravelDate       time.Time `json:"travel_date"`
	Duration         int       `json:"duration"` // days
	Travelers        int       `json:"travelers"`
	BudgetPreference string    `json:"budget_preference"` // budget, mid-range, luxury
}

// calculateTravelerMultiplier calculates cost multiplier based on number of travelers
func (tcp *TravelCostPredictor) calculateTravelerMultiplier(travelers int) float64 {
	// Accommodations often have group discounts
	if travelers == 1 {
		return 1.0
	} else if travelers == 2 {
		return 1.6 // Shared room
	} else if travelers <= 4 {
		return 1.8 // Family rooms
	} else {
		return 2.0 // Group accommodations
	}
}

// calculateSeasonalMultiplier calculates seasonal pricing multiplier
func (tcp *TravelCostPredictor) calculateSeasonalMultiplier(travelDate time.Time, destination string) float64 {
	seasonal, exists := tcp.seasonality[strings.ToLower(destination)]
	if !exists {
		seasonal = tcp.getGlobalSeasonalFactors()
	}

	month := travelDate.Month()

	// Peak season (summer in northern hemisphere)
	if month >= 6 && month <= 8 {
		return seasonal.PeakSeasonMultiplier
	}

	// Mid season (spring/fall)
	if month == 4 || month == 5 || month == 9 || month == 10 {
		return seasonal.MidSeasonMultiplier
	}

	// Off season (winter)
	return seasonal.OffSeasonMultiplier
}

// calculateBudgetMultiplier calculates multiplier based on budget preference
func (tcp *TravelCostPredictor) calculateBudgetMultiplier(budgetPreference string) float64 {
	switch strings.ToLower(budgetPreference) {
	case "budget":
		return 0.7
	case "mid-range":
		return 1.0
	case "luxury":
		return 1.8
	case "ultra-luxury":
		return 3.0
	default:
		return 1.0
	}
}

// calculateConfidenceLevel calculates prediction confidence
func (tcp *TravelCostPredictor) calculateConfidenceLevel(destination string, factors CostFactors) float64 {
	// Higher confidence for destinations with more data
	_, exists := tcp.costFactors[strings.ToLower(destination)]
	if !exists {
		return 0.6 // Lower confidence for unknown destinations
	}

	// Base confidence varies by destination popularity
	if factors.PopularityMultiplier > 1.5 {
		return 0.9 // High confidence for popular destinations
	} else if factors.PopularityMultiplier > 1.0 {
		return 0.8 // Medium confidence
	}

	return 0.7 // Standard confidence
}

// generateCostRecommendations generates cost-saving recommendations
func (tcp *TravelCostPredictor) generateCostRecommendations(req CostPredictionRequest, factors CostFactors, totalCost float64) []string {
	var recommendations []string

	if factors.CostOfLivingIndex > 1.5 {
		recommendations = append(recommendations, "Consider staying in budget accommodations or hostels to reduce costs")
		recommendations = append(recommendations, "Try local street food and markets for affordable dining")
	}

	if req.Duration > 7 {
		recommendations = append(recommendations, "Look for weekly accommodation discounts")
		recommendations = append(recommendations, "Consider public transportation passes for extended stays")
	}

	if req.Travelers > 2 {
		recommendations = append(recommendations, "Book group accommodations and activities for discounts")
		recommendations = append(recommendations, "Split costs for private tours and transportation")
	}

	return recommendations
}

// generateSeasonalAdvice generates advice based on travel timing
func (tcp *TravelCostPredictor) generateSeasonalAdvice(travelDate time.Time, destination string) string {
	month := travelDate.Month()

	if month >= 6 && month <= 8 {
		return "Peak season - expect higher prices but better weather and more activities"
	} else if month == 4 || month == 5 || month == 9 || month == 10 {
		return "Shoulder season - good balance of weather and pricing"
	} else {
		return "Off season - lower prices but weather may be less ideal"
	}
}

// initializeCostFactors initializes cost factors for major destinations
func initializeCostFactors() map[string]CostFactors {
	return map[string]CostFactors{
		"paris": {
			BaseAccommodationCost: 120.0,
			BaseFoodCost:          45.0,
			BaseTransportCost:     25.0,
			BaseActivityCost:      35.0,
			CostOfLivingIndex:     1.3,
			PopularityMultiplier:  1.8,
		},
		"tokyo": {
			BaseAccommodationCost: 100.0,
			BaseFoodCost:          40.0,
			BaseTransportCost:     20.0,
			BaseActivityCost:      30.0,
			CostOfLivingIndex:     1.2,
			PopularityMultiplier:  1.7,
		},
		"london": {
			BaseAccommodationCost: 140.0,
			BaseFoodCost:          50.0,
			BaseTransportCost:     30.0,
			BaseActivityCost:      40.0,
			CostOfLivingIndex:     1.4,
			PopularityMultiplier:  1.9,
		},
		"new york": {
			BaseAccommodationCost: 180.0,
			BaseFoodCost:          55.0,
			BaseTransportCost:     35.0,
			BaseActivityCost:      45.0,
			CostOfLivingIndex:     1.5,
			PopularityMultiplier:  2.0,
		},
		"bangkok": {
			BaseAccommodationCost: 40.0,
			BaseFoodCost:          15.0,
			BaseTransportCost:     10.0,
			BaseActivityCost:      20.0,
			CostOfLivingIndex:     0.6,
			PopularityMultiplier:  1.4,
		},
		"rome": {
			BaseAccommodationCost: 90.0,
			BaseFoodCost:          35.0,
			BaseTransportCost:     20.0,
			BaseActivityCost:      25.0,
			CostOfLivingIndex:     1.1,
			PopularityMultiplier:  1.6,
		},
	}
}

// initializeSeasonalFactors initializes seasonal factors
func initializeSeasonalFactors() map[string]SeasonalFactors {
	return map[string]SeasonalFactors{
		"global": {
			PeakSeasonMultiplier: 1.4,
			MidSeasonMultiplier:  1.1,
			OffSeasonMultiplier:  0.8,
			HolidayMultiplier:    1.6,
			WeekendMultiplier:    1.2,
		},
	}
}

// getGlobalAverageCostFactors returns global average cost factors
func (tcp *TravelCostPredictor) getGlobalAverageCostFactors() CostFactors {
	return CostFactors{
		BaseAccommodationCost: 80.0,
		BaseFoodCost:          30.0,
		BaseTransportCost:     20.0,
		BaseActivityCost:      25.0,
		CostOfLivingIndex:     1.0,
		PopularityMultiplier:  1.0,
	}
}

// getGlobalSeasonalFactors returns global seasonal factors
func (tcp *TravelCostPredictor) getGlobalSeasonalFactors() SeasonalFactors {
	return tcp.seasonality["global"]
}
