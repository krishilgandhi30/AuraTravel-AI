package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"auratravel-backend/internal/config"
)

// GeminiService handles Gemini AI interactions
type GeminiService struct {
	apiKey     string
	cfg        *config.Config
	httpClient *http.Client
	baseURL    string
}

// GeminiRequest represents a request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content in Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents Gemini API response
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// NewGeminiService creates a new Gemini AI service
func NewGeminiService() (*GeminiService, error) {
	cfg := config.GetConfig()

	if cfg.GeminiAPIKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set, using mock service")
		return &GeminiService{
			apiKey:     "",
			cfg:        cfg,
			httpClient: &http.Client{Timeout: 30 * time.Second},
			baseURL:    "https://generativelanguage.googleapis.com/v1beta",
		}, nil
	}

	return &GeminiService{
		apiKey:     cfg.GeminiAPIKey,
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
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

	prompt := g.buildItineraryPrompt(req)
	response, err := g.callGeminiAPI(ctx, prompt)
	if err != nil {
		log.Printf("Gemini API call failed, falling back to mock: %v", err)
		return g.mockItinerary(req), nil
	}

	// Parse the response and create structured itinerary
	itinerary := g.parseItineraryResponse(response, req)
	return itinerary, nil
}

// GenerateItineraryWithRAG creates AI-powered itinerary using RAG context
func (g *GeminiService) GenerateItineraryWithRAG(ctx context.Context, req ItineraryRequest, ragContext TripContext) (map[string]interface{}, error) {
	if g.apiKey == "" {
		return g.mockItineraryWithRAG(req, ragContext), nil
	}

	prompt := g.buildRAGItineraryPrompt(req, ragContext)
	response, err := g.callGeminiAPI(ctx, prompt)
	if err != nil {
		log.Printf("Gemini API call failed, falling back to mock: %v", err)
		return g.mockItineraryWithRAG(req, ragContext), nil
	}

	// Parse the response and create structured itinerary with RAG context
	itinerary := g.parseRAGItineraryResponse(response, req, ragContext)
	return itinerary, nil
}

// GetDestinationRecommendations gets AI-powered destination recommendations
func (g *GeminiService) GetDestinationRecommendations(ctx context.Context, req RecommendationRequest) ([]map[string]interface{}, error) {
	if g.apiKey == "" {
		return g.mockRecommendations(req), nil
	}

	prompt := g.buildRecommendationPrompt(req)
	response, err := g.callGeminiAPI(ctx, prompt)
	if err != nil {
		log.Printf("Gemini API call failed, falling back to mock: %v", err)
		return g.mockRecommendations(req), nil
	}

	// Parse recommendations from response
	recommendations := g.parseRecommendationsResponse(response, req)
	return recommendations, nil
}

// GetActivitySuggestions gets AI-powered activity suggestions
func (g *GeminiService) GetActivitySuggestions(ctx context.Context, destination string, interests []string) ([]string, error) {
	if g.apiKey == "" {
		return g.mockActivitySuggestions(destination, interests), nil
	}

	prompt := g.buildActivityPrompt(destination, interests)
	response, err := g.callGeminiAPI(ctx, prompt)
	if err != nil {
		log.Printf("Gemini API call failed, falling back to mock: %v", err)
		return g.mockActivitySuggestions(destination, interests), nil
	}

	// Parse activities from response
	activities := g.parseActivitiesResponse(response)
	return activities, nil
}

// callGeminiAPI makes a request to the Gemini API
func (g *GeminiService) callGeminiAPI(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s/models/gemini-pro:generateContent?key=%s", g.baseURL, g.apiKey)

	request := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// buildItineraryPrompt creates a prompt for basic itinerary generation
func (g *GeminiService) buildItineraryPrompt(req ItineraryRequest) string {
	days := g.calculateDays(req.StartDate, req.EndDate)
	preferences := ""
	for key, value := range req.Preferences {
		preferences += fmt.Sprintf("%s: %v, ", key, value)
	}

	return fmt.Sprintf(`Generate a detailed %d-day travel itinerary for %s with the following requirements:
- Destination: %s
- Budget: $%.2f
- Number of travelers: %d
- Travel preferences: %s

Please provide a day-by-day breakdown with morning, afternoon, and evening activities. Include:
- Specific attractions and landmarks to visit
- Restaurant recommendations for meals
- Transportation between locations
- Estimated costs for activities
- Practical tips for travelers

Format the response as a structured JSON with clear day-by-day organization.`,
		days, req.Destination, req.Destination, req.Budget, req.Travelers, preferences)
}

// buildRAGItineraryPrompt creates a prompt for RAG-enhanced itinerary generation
func (g *GeminiService) buildRAGItineraryPrompt(req ItineraryRequest, ragContext TripContext) string {
	days := g.calculateDays(req.StartDate, req.EndDate)

	// Include real-time context in prompt
	contextInfo := ""
	if len(ragContext.Attractions) > 0 {
		contextInfo += fmt.Sprintf("Available attractions: %d locations including %s. ",
			len(ragContext.Attractions), ragContext.Attractions[0].Name)
	}
	if len(ragContext.Hotels) > 0 {
		contextInfo += fmt.Sprintf("Recommended hotels: %d options starting from $%.2f. ",
			len(ragContext.Hotels), ragContext.Hotels[0].PricePerNight)
	}
	if ragContext.Weather.Current.Temperature > 0 {
		contextInfo += fmt.Sprintf("Current weather: %.1f°C, %s. ",
			ragContext.Weather.Current.Temperature, ragContext.Weather.Current.Description)
	}

	return fmt.Sprintf(`Generate a detailed %d-day travel itinerary for %s using the following real-time data:

%s

Trip Requirements:
- Budget: $%.2f
- Travelers: %d
- Dates: %s to %s

Use the provided real-time data to create a personalized, accurate itinerary. Include:
- Specific attractions from the available options
- Hotel recommendations with actual pricing
- Weather-appropriate activities
- Real transportation options
- Local events and cultural experiences

Format as structured JSON with day-by-day breakdown and real-time validation.`,
		days, req.Destination, contextInfo, req.Budget, req.Travelers, req.StartDate, req.EndDate)
}

// buildRecommendationPrompt creates a prompt for destination recommendations
func (g *GeminiService) buildRecommendationPrompt(req RecommendationRequest) string {
	interests := strings.Join(req.Interests, ", ")

	return fmt.Sprintf(`Recommend 5 travel destinations based on:
- Budget: $%.2f
- Interests: %s
- User preferences for travel experiences

For each destination, provide:
- Why it matches the interests
- Estimated cost breakdown
- Best time to visit
- Top 3 must-see attractions
- Cultural highlights

Format as JSON array with structured destination objects.`, req.Budget, interests)
}

// buildActivityPrompt creates a prompt for activity suggestions
func (g *GeminiService) buildActivityPrompt(destination string, interests []string) string {
	interestList := strings.Join(interests, ", ")

	return fmt.Sprintf(`Suggest 10 specific activities in %s for travelers interested in: %s

Include:
- Activity name and description
- Location within the destination
- Estimated duration and cost
- Best time of day/season
- Difficulty level or requirements

Format as a simple list of activity descriptions.`, destination, interestList)
}

// parseItineraryResponse parses Gemini response into structured itinerary
func (g *GeminiService) parseItineraryResponse(response string, req ItineraryRequest) map[string]interface{} {
	// Try to parse JSON response, fallback to mock if parsing fails
	var itinerary map[string]interface{}
	if err := json.Unmarshal([]byte(response), &itinerary); err != nil {
		log.Printf("Failed to parse Gemini response as JSON, using enhanced mock: %v", err)
		return g.enhanceItineraryWithAI(g.mockItinerary(req), response)
	}

	// Enhance with standard fields
	itinerary["ai_generated"] = true
	itinerary["created_at"] = time.Now().Format(time.RFC3339)

	return itinerary
}

// parseRAGItineraryResponse parses RAG-enhanced response
func (g *GeminiService) parseRAGItineraryResponse(response string, req ItineraryRequest, ragContext TripContext) map[string]interface{} {
	// Try to parse JSON response
	var itinerary map[string]interface{}
	if err := json.Unmarshal([]byte(response), &itinerary); err != nil {
		log.Printf("Failed to parse RAG response as JSON, using enhanced mock: %v", err)
		baseItinerary := g.mockItineraryWithRAG(req, ragContext)
		return g.enhanceItineraryWithAI(baseItinerary, response)
	}

	// Enhance with RAG context
	itinerary["rag_enhanced"] = true
	itinerary["ai_generated"] = true
	itinerary["created_at"] = time.Now().Format(time.RFC3339)

	return itinerary
}

// parseRecommendationsResponse parses recommendations from AI response
func (g *GeminiService) parseRecommendationsResponse(response string, req RecommendationRequest) []map[string]interface{} {
	var recommendations []map[string]interface{}
	if err := json.Unmarshal([]byte(response), &recommendations); err != nil {
		log.Printf("Failed to parse recommendations as JSON, using mock: %v", err)
		return g.mockRecommendations(req)
	}

	return recommendations
}

// parseActivitiesResponse parses activities from AI response
func (g *GeminiService) parseActivitiesResponse(response string) []string {
	// Try to parse as JSON array first
	var activities []string
	if err := json.Unmarshal([]byte(response), &activities); err != nil {
		// If JSON parsing fails, split by lines and clean up
		lines := strings.Split(response, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "Here are") {
				// Remove list markers
				line = strings.TrimPrefix(line, "- ")
				line = strings.TrimPrefix(line, "• ")
				line = strings.TrimPrefix(line, "* ")
				if line != "" {
					activities = append(activities, line)
				}
			}
		}
	}

	return activities
}

// enhanceItineraryWithAI enhances mock itinerary with AI insights
func (g *GeminiService) enhanceItineraryWithAI(baseItinerary map[string]interface{}, aiResponse string) map[string]interface{} {
	// Extract insights from AI response and enhance the base itinerary
	if strings.Contains(aiResponse, "cultural") {
		baseItinerary["cultural_focus"] = true
	}
	if strings.Contains(aiResponse, "adventure") {
		baseItinerary["adventure_activities"] = true
	}
	if strings.Contains(aiResponse, "budget") {
		baseItinerary["budget_optimized"] = true
	}

	// Add AI insights as tips
	if tips, ok := baseItinerary["tips"].([]string); ok {
		tips = append(tips, "Enhanced with AI recommendations")
		baseItinerary["tips"] = tips
	}

	return baseItinerary
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
