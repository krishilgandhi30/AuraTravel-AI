package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"auratravel-backend/internal/config"
)

// GeminiService handles Gemini AI interactions
type GeminiService struct {
	apiKey string
	cfg    *config.Config
}

// NewGeminiService creates a new Gemini AI service
func NewGeminiService() (*GeminiService, error) {
	cfg := config.GetConfig()

	if cfg.GeminiAPIKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set, using mock service")
		return &GeminiService{
			apiKey: "",
			cfg:    cfg,
		}, nil
	}

	return &GeminiService{
		apiKey: cfg.GeminiAPIKey,
		cfg:    cfg,
	}, nil
}

// ItineraryRequest represents itinerary generation request
type ItineraryRequest struct {
	Destination string                 `json:"destination"`
	StartDate   string                 `json:"start_date"`
	EndDate     string                 `json:"end_date"`
	Budget      float64                `json:"budget"`
	Travelers   int                    `json:"travelers"`
	Preferences map[string]interface{} `json:"preferences"`
}

// RecommendationRequest represents recommendation request
type RecommendationRequest struct {
	UserID    string   `json:"user_id"`
	Budget    float64  `json:"budget"`
	Interests []string `json:"interests"`
}

// GenerateItinerary creates AI-powered itinerary
func (g *GeminiService) GenerateItinerary(ctx context.Context, req ItineraryRequest) (map[string]interface{}, error) {
	if g.apiKey == "" {
		return g.mockItinerary(req), nil
	}

	// TODO: Implement actual Gemini API call via HTTP
	// For now, return mock data
	return g.mockItinerary(req), nil
}

// GenerateItineraryWithRAG creates AI-powered itinerary using RAG context
func (g *GeminiService) GenerateItineraryWithRAG(ctx context.Context, req ItineraryRequest, ragContext TripContext) (map[string]interface{}, error) {
	if g.apiKey == "" {
		return g.mockItineraryWithRAG(req, ragContext), nil
	}

	// TODO: Implement actual Gemini API call with RAG context
	// For now, return enhanced mock data
	return g.mockItineraryWithRAG(req, ragContext), nil
}

// GetDestinationRecommendations gets AI-powered destination recommendations
func (g *GeminiService) GetDestinationRecommendations(ctx context.Context, req RecommendationRequest) ([]map[string]interface{}, error) {
	if g.apiKey == "" {
		return g.mockRecommendations(req), nil
	}

	// TODO: Implement actual Gemini API call via HTTP
	// For now, return mock data
	return g.mockRecommendations(req), nil
}

// GetActivitySuggestions gets AI-powered activity suggestions
func (g *GeminiService) GetActivitySuggestions(ctx context.Context, destination string, interests []string) ([]string, error) {
	if g.apiKey == "" {
		return g.mockActivitySuggestions(destination, interests), nil
	}

	// TODO: Implement actual Gemini API call via HTTP
	// For now, return mock data
	return g.mockActivitySuggestions(destination, interests), nil
}

// Mock implementations

func (g *GeminiService) mockItinerary(req ItineraryRequest) map[string]interface{} {
	return map[string]interface{}{
		"destination": req.Destination,
		"duration":    g.calculateDays(req.StartDate, req.EndDate),
		"budget":      req.Budget,
		"travelers":   req.Travelers,
		"day_1": map[string]interface{}{
			"morning":   fmt.Sprintf("Arrive in %s, check into hotel", req.Destination),
			"afternoon": "City orientation tour and local lunch",
			"evening":   "Welcome dinner at local restaurant",
		},
		"day_2": map[string]interface{}{
			"morning":   "Visit main attractions and landmarks",
			"afternoon": "Cultural sites and museums",
			"evening":   "Local entertainment and dining",
		},
		"day_3": map[string]interface{}{
			"morning":   "Adventure activities or excursions",
			"afternoon": "Shopping and leisure time",
			"evening":   "Sunset viewing and farewell dinner",
		},
		"tips": []string{
			"Book accommodations in advance",
			"Try local cuisine and street food",
			"Respect local customs and traditions",
			"Keep important documents safe",
		},
		"ai_generated": true,
		"created_at":   time.Now().Format(time.RFC3339),
	}
}

func (g *GeminiService) mockRecommendations(req RecommendationRequest) []map[string]interface{} {
	recommendations := []map[string]interface{}{
		{
			"destination":  "Paris, France",
			"type":         "cultural",
			"budget_match": req.Budget >= 1000,
			"reasons": []string{
				"Rich cultural heritage",
				"World-class museums",
				"Excellent cuisine",
			},
			"estimated_cost": 1200,
			"best_time":      "April-June, September-October",
			"highlights": []string{
				"Eiffel Tower",
				"Louvre Museum",
				"Notre-Dame Cathedral",
			},
		},
		{
			"destination":  "Tokyo, Japan",
			"type":         "adventure",
			"budget_match": req.Budget >= 1500,
			"reasons": []string{
				"Unique blend of tradition and modernity",
				"Amazing food scene",
				"Efficient transportation",
			},
			"estimated_cost": 1800,
			"best_time":      "March-May, September-November",
			"highlights": []string{
				"Shibuya Crossing",
				"Senso-ji Temple",
				"Tsukiji Fish Market",
			},
		},
		{
			"destination":  "Bali, Indonesia",
			"type":         "relaxation",
			"budget_match": req.Budget >= 600,
			"reasons": []string{
				"Beautiful beaches and landscapes",
				"Rich spiritual culture",
				"Affordable luxury",
			},
			"estimated_cost": 800,
			"best_time":      "April-October",
			"highlights": []string{
				"Ubud Rice Terraces",
				"Tanah Lot Temple",
				"Seminyak Beach",
			},
		},
	}

	return recommendations
}

func (g *GeminiService) mockItineraryWithRAG(req ItineraryRequest, ragContext TripContext) map[string]interface{} {
	itinerary := map[string]interface{}{
		"destination":  req.Destination,
		"duration":     g.calculateDays(req.StartDate, req.EndDate),
		"budget":       req.Budget,
		"travelers":    req.Travelers,
		"ai_generated": true,
		"rag_enhanced": true,
		"created_at":   time.Now().Format(time.RFC3339),
	}

	// Add real-time weather context
	if len(ragContext.Weather.Forecast) > 0 {
		itinerary["weather_info"] = map[string]interface{}{
			"current":  ragContext.Weather.Current,
			"forecast": ragContext.Weather.Forecast,
			"tips":     g.generateWeatherTips(ragContext.Weather.Forecast),
		}
	}

	// Build day-by-day itinerary using real attractions
	days := g.calculateDays(req.StartDate, req.EndDate)
	for day := 1; day <= days; day++ {
		dayKey := fmt.Sprintf("day_%d", day)
		dayPlan := g.buildDayPlan(day, ragContext, req.Preferences)
		itinerary[dayKey] = dayPlan
	}

	// Add hotel recommendations
	if len(ragContext.Hotels) > 0 {
		var hotelRecs []map[string]interface{}
		for _, hotel := range ragContext.Hotels[:min(3, len(ragContext.Hotels))] {
			hotelRecs = append(hotelRecs, map[string]interface{}{
				"name":            hotel.Name,
				"rating":          hotel.Rating,
				"price_per_night": hotel.PricePerNight,
				"availability":    hotel.Available,
				"booking_url":     hotel.BookingURL,
			})
		}
		itinerary["recommended_hotels"] = hotelRecs
	}

	// Add transportation options
	if len(ragContext.Transportation) > 0 {
		var transportRecs []map[string]interface{}
		for _, transport := range ragContext.Transportation {
			transportRecs = append(transportRecs, map[string]interface{}{
				"type":         transport.Type,
				"provider":     transport.Provider,
				"price":        transport.Price,
				"availability": transport.Available,
				"booking_url":  transport.BookingURL,
			})
		}
		itinerary["transportation"] = transportRecs
	}

	// Add local events
	if len(ragContext.LocalEvents) > 0 {
		var eventRecs []map[string]interface{}
		for _, event := range ragContext.LocalEvents {
			eventRecs = append(eventRecs, map[string]interface{}{
				"name":        event.Name,
				"date":        event.Date.Format("2006-01-02"),
				"description": event.Description,
				"location":    event.Location,
				"price":       event.Price,
			})
		}
		itinerary["local_events"] = eventRecs
	}

	// Add EMT inventory items
	if len(ragContext.EMTInventory) > 0 {
		var emtRecs []map[string]interface{}
		for _, item := range ragContext.EMTInventory[:min(5, len(ragContext.EMTInventory))] {
			emtRecs = append(emtRecs, map[string]interface{}{
				"name":        item.Name,
				"type":        item.Type,
				"description": item.Description,
				"available":   item.Available,
				"contact":     item.Contact,
			})
		}
		itinerary["emt_services"] = emtRecs
	}

	// Generate contextual tips based on real data
	itinerary["tips"] = g.generateContextualTips(ragContext, req)

	return itinerary
}

func (g *GeminiService) buildDayPlan(day int, ragContext TripContext, preferences map[string]interface{}) map[string]interface{} {
	dayPlan := map[string]interface{}{}

	// Distribute attractions across time slots
	attractionsPerDay := 2
	startIdx := (day - 1) * attractionsPerDay
	endIdx := min(startIdx+attractionsPerDay, len(ragContext.Attractions))

	if startIdx < len(ragContext.Attractions) {
		dayAttractions := ragContext.Attractions[startIdx:endIdx]

		if len(dayAttractions) > 0 {
			dayPlan["morning"] = fmt.Sprintf("Visit %s - %s",
				dayAttractions[0].Name,
				dayAttractions[0].Description)
		}

		if len(dayAttractions) > 1 {
			dayPlan["afternoon"] = fmt.Sprintf("Explore %s - %s",
				dayAttractions[1].Name,
				dayAttractions[1].Description)
		}
	}

	// Add evening recommendations based on day
	switch day % 3 {
	case 1:
		dayPlan["evening"] = "Welcome dinner at local restaurant, orientation walk"
	case 2:
		dayPlan["evening"] = "Cultural performance or local entertainment"
	case 0:
		dayPlan["evening"] = "Sunset viewing and leisure dining"
	}

	return dayPlan
}

func (g *GeminiService) generateWeatherTips(forecasts []WeatherCondition) []string {
	tips := []string{}

	for _, forecast := range forecasts {
		if forecast.Temperature > 30 {
			tips = append(tips, "Pack light, breathable clothing and sunscreen")
		} else if forecast.Temperature < 10 {
			tips = append(tips, "Bring warm clothes and layers")
		}

		if strings.Contains(strings.ToLower(forecast.Description), "rain") {
			tips = append(tips, "Pack rain gear and plan indoor activities")
		}

		if forecast.Humidity > 80 {
			tips = append(tips, "Stay hydrated and take breaks in air-conditioned spaces")
		}
	}

	if len(tips) == 0 {
		tips = append(tips, "Weather looks favorable for outdoor activities")
	}

	return removeDuplicateStrings(tips)
}

func (g *GeminiService) generateContextualTips(ragContext TripContext, req ItineraryRequest) []string {
	tips := []string{
		"Book accommodations and attractions in advance",
		"Keep important documents and emergency contacts handy",
		"Respect local customs and traditions",
	}

	// Budget-based tips
	if req.Budget < 500 {
		tips = append(tips, "Look for free walking tours and public transportation")
		tips = append(tips, "Try street food and local markets for affordable meals")
	} else if req.Budget > 2000 {
		tips = append(tips, "Consider premium experiences and fine dining")
		tips = append(tips, "Private tours and luxury accommodations recommended")
	}

	// Traveler count tips
	if req.Travelers > 4 {
		tips = append(tips, "Book group accommodations and activities")
		tips = append(tips, "Consider splitting into smaller groups for some activities")
	}

	// EMT availability tips
	if len(ragContext.EMTInventory) > 0 {
		tips = append(tips, "Emergency medical services and facilities are available")
	}

	// Transportation tips
	hasPublicTransport := false
	for _, transport := range ragContext.Transportation {
		if transport.Type == "public_transport" {
			hasPublicTransport = true
			break
		}
	}

	if hasPublicTransport {
		tips = append(tips, "Public transportation is available and cost-effective")
	} else {
		tips = append(tips, "Consider renting a car or using ride-sharing services")
	}

	return tips
}

func removeDuplicateStrings(strings []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, str := range strings {
		if !keys[str] {
			keys[str] = true
			result = append(result, str)
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *GeminiService) mockActivitySuggestions(destination string, interests []string) []string {
	baseActivities := []string{
		fmt.Sprintf("Guided walking tour of %s", destination),
		"Local cooking class experience",
		"Visit to historical landmarks",
		"Traditional market exploration",
		"Sunset photography session",
	}

	// Add interest-specific activities
	for _, interest := range interests {
		switch interest {
		case "adventure":
			baseActivities = append(baseActivities, "Hiking and outdoor activities", "Water sports and adventures")
		case "culture":
			baseActivities = append(baseActivities, "Museum visits", "Cultural performances")
		case "food":
			baseActivities = append(baseActivities, "Food tours", "Wine tasting experiences")
		case "nature":
			baseActivities = append(baseActivities, "Nature reserves visit", "Botanical gardens tour")
		}
	}

	return baseActivities
}

func (g *GeminiService) calculateDays(startDate, endDate string) int {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return int(end.Sub(start).Hours()/24) + 1
}

// Shutdown closes the Gemini service
func (g *GeminiService) Shutdown(ctx context.Context) error {
	log.Println("Gemini service shut down successfully")
	return nil
}
