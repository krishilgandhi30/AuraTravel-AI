package services

import (
	"context"
	"fmt"
	"log"
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
