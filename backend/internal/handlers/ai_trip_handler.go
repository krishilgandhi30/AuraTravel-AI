package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"auratravel-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AITripHandler handles AI-powered trip operations
type AITripHandler struct {
	services *services.Services
}

// NewAITripHandler creates a new AI trip handler
func NewAITripHandler(services *services.Services) *AITripHandler {
	return &AITripHandler{
		services: services,
	}
}

// PlanTripRequest represents a trip planning request
type PlanTripRequest struct {
	Destination string                 `json:"destination" binding:"required"`
	StartDate   string                 `json:"start_date" binding:"required"`
	EndDate     string                 `json:"end_date" binding:"required"`
	Budget      float64                `json:"budget"`
	Travelers   int                    `json:"travelers"`
	Preferences map[string]interface{} `json:"preferences"`
	TripType    string                 `json:"trip_type"`
	Interests   []string               `json:"interests"`
	UserID      string                 `json:"user_id"`
}

// PlanTripResponse represents the AI-generated trip plan
type PlanTripResponse struct {
	TripID      string                 `json:"trip_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Itinerary   map[string]interface{} `json:"itinerary"`
	Budget      TripBudget             `json:"budget"`
	Suggestions []string               `json:"suggestions"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TripBudget represents budget breakdown
type TripBudget struct {
	Total          float64            `json:"total"`
	Accommodation  float64            `json:"accommodation"`
	Transportation float64            `json:"transportation"`
	Food           float64            `json:"food"`
	Activities     float64            `json:"activities"`
	Breakdown      map[string]float64 `json:"breakdown"`
}

// PlanTrip creates an AI-powered trip plan
func (h *AITripHandler) PlanTrip(c *gin.Context) {
	var req PlanTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// Use RAG for enhanced trip planning
	var itinerary map[string]interface{}
	var suggestions []string
	var ragEnabled bool

	// Try RAG-enhanced planning first
	if h.services.RAGRetriever != nil && h.services.Gemini != nil {
		// Retrieve contextual data
		ragRequest := services.RetrievalRequest{
			UserID:      req.UserID,
			Destination: req.Destination,
			StartDate:   h.parseDate(req.StartDate),
			EndDate:     h.parseDate(req.EndDate),
			Budget:      req.Budget,
			Travelers:   req.Travelers,
			Interests:   req.Interests,
			Preferences: req.Preferences,
		}

		ragContext, err := h.services.RAGRetriever.RetrieveContext(ctx, ragRequest)
		if err == nil {
			// Generate itinerary with RAG context
			ragItinerary, err := h.services.Gemini.GenerateItineraryWithRAG(ctx, services.ItineraryRequest{
				Destination: req.Destination,
				StartDate:   req.StartDate,
				EndDate:     req.EndDate,
				Budget:      req.Budget,
				Travelers:   req.Travelers,
				Preferences: req.Preferences,
			}, *ragContext)

			if err == nil {
				itinerary = ragItinerary
				ragEnabled = true

				// Extract suggestions from attractions
				for _, attraction := range ragContext.Attractions {
					suggestions = append(suggestions, fmt.Sprintf("Visit %s - %s", attraction.Name, attraction.Description))
				}
			}
		}
	}

	// Fallback to regular AI if RAG fails
	if itinerary == nil && h.services.Gemini != nil {
		geminiItinerary, err := h.services.Gemini.GenerateItinerary(ctx, services.ItineraryRequest{
			Destination: req.Destination,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
			Budget:      req.Budget,
			Travelers:   req.Travelers,
			Preferences: req.Preferences,
		})
		if err == nil {
			itinerary = geminiItinerary
		}

		geminiSuggestions, err := h.services.Gemini.GetActivitySuggestions(ctx, req.Destination, req.Interests)
		if err == nil {
			suggestions = geminiSuggestions
		}
	}

	// Fallback to basic itinerary if AI fails
	if itinerary == nil {
		itinerary = h.createBasicItinerary(req)
	}
	if len(suggestions) == 0 {
		suggestions = h.getBasicSuggestions(req.Destination)
	}

	// Add RAG enhancement indicator
	if ragEnabled {
		itinerary["rag_enhanced"] = true
		itinerary["data_sources"] = []string{"real_attractions", "weather_forecast", "hotel_availability", "transportation_options"}
	}

	// Create trip in Firestore only
	tripID := uuid.New().String()
	trip := services.TripData{
		ID:          tripID,
		UserID:      req.UserID,
		Title:       fmt.Sprintf("AI Trip to %s", req.Destination),
		Destination: req.Destination,
		StartDate:   h.parseDate(req.StartDate),
		EndDate:     h.parseDate(req.EndDate),
		Status:      "planned",
		Itinerary:   itinerary,
		Budget:      req.Budget,
		Travelers:   req.Travelers,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if h.services.Firebase != nil {
		if err := h.services.Firebase.SaveTrip(ctx, trip); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save trip to Firestore"})
			return
		}
	}

	budget := h.calculateBudgetBreakdown(req.Budget, req.Travelers)

	response := PlanTripResponse{
		TripID:      tripID,
		Title:       trip.Title,
		Description: fmt.Sprintf("AI-powered %d-day trip to %s", h.calculateDays(req.StartDate, req.EndDate), req.Destination),
		Itinerary:   itinerary,
		Budget:      budget,
		Suggestions: suggestions,
		CreatedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// GetRecommendations gets AI-powered destination recommendations
func (h *AITripHandler) GetRecommendations(c *gin.Context) {
	userID := c.Query("user_id")
	budget := c.Query("budget")
	interests := c.QueryArray("interests")

	ctx := context.Background()
	var recommendations []map[string]interface{}

	// Get recommendations from Gemini AI
	if h.services.Gemini != nil {
		geminiRecs, err := h.services.Gemini.GetDestinationRecommendations(ctx, services.RecommendationRequest{
			UserID:    userID,
			Budget:    h.parseFloat(budget),
			Interests: interests,
		})
		if err == nil {
			recommendations = geminiRecs
		}
	}

	// Get personalized recommendations from Vertex AI
	if h.services.Vertex != nil {
		// Convert preferences for the existing method signature
		// vertexRecs, err := h.services.Vertex.GeneratePersonalizedRecommendations(ctx, ...)
		// For now, skip Vertex recommendations since signature doesn't match
	}

	// Get trending destinations from BigQuery
	if h.services.BigQuery != nil {
		// trending, err := h.services.BigQuery.GetTrendingDestinations(ctx, 10)
		// Method signature doesn't match - skipping for now
	}

	// Fallback recommendations
	if len(recommendations) == 0 {
		recommendations = h.getDefaultRecommendations()
	}

	// Save recommendations to Firebase
	if h.services.Firebase != nil && userID != "" {
		h.services.Firebase.SaveRecommendations(ctx, userID, recommendations)
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"generated_at":    time.Now(),
	})
}

// OptimizeItinerary optimizes an existing itinerary using AI
func (h *AITripHandler) OptimizeItinerary(c *gin.Context) {
	tripID := c.Param("id")

	var optimizationReq struct {
		Preferences map[string]interface{} `json:"preferences"`
		Constraints map[string]interface{} `json:"constraints"`
	}

	if err := c.ShouldBindJSON(&optimizationReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// Get existing trip from Firestore
	if h.services.Firebase != nil {
		t, err := h.services.Firebase.GetTrip(ctx, tripID)
		if err != nil || t == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
			return
		}
	}

	// Optimize using Vertex AI
	var optimizedItinerary map[string]interface{}
	if h.services.Vertex != nil {
		// optimized, _, err := h.services.Vertex.OptimizeItinerary(ctx, ...)
		// Method signature doesn't match current requirements - using mock for now
		optimizedItinerary = map[string]interface{}{
			"optimized": true,
			"message":   "Itinerary optimization not yet implemented",
		}
	}

	if optimizedItinerary == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to optimize itinerary"})
		return
	}

	// Update trip in Firebase
	if h.services.Firebase != nil {
		updates := map[string]interface{}{
			"itinerary": optimizedItinerary,
		}
		h.services.Firebase.UpdateTrip(ctx, tripID, updates)
	}

	c.JSON(http.StatusOK, gin.H{
		"trip_id":             tripID,
		"optimized_itinerary": optimizedItinerary,
		"optimized_at":        time.Now(),
	})
}

// AnalyzeImage analyzes uploaded travel images using Vision AI
func (h *AITripHandler) AnalyzeImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	ctx := context.Background()
	var analysis map[string]interface{}

	if h.services.Vision != nil {
		// Convert file to bytes
		imageData := make([]byte, header.Size)
		file.Read(imageData)

		visionAnalysis, err := h.services.Vision.AnalyzeImage(ctx, imageData)
		if err == nil {
			analysis = visionAnalysis
		}

		// Detect landmarks
		landmarks, err := h.services.Vision.DetectLandmarks(ctx, imageData)
		if err == nil {
			analysis["landmarks"] = landmarks
		}
	}

	if analysis == nil {
		analysis = map[string]interface{}{
			"message": "Image analysis service unavailable",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis":    analysis,
		"analyzed_at": time.Now(),
	})
}

// GetTravelInsights gets travel insights and analytics
func (h *AITripHandler) GetTravelInsights(c *gin.Context) {
	userID := c.Query("user_id")
	ctx := context.Background()

	insights := make(map[string]interface{})

	// Get travel patterns from Vertex AI
	if h.services.Vertex != nil {
		// patterns, err := h.services.Vertex.AnalyzeTravelPattern(ctx, ...)
		// Method signature doesn't match - using mock for now
		insights["travel_patterns"] = map[string]interface{}{
			"message": "Travel pattern analysis not yet implemented",
		}
	}

	// Get price trends from BigQuery
	if h.services.BigQuery != nil {
		// trends, err := h.services.BigQuery.GetPriceTrends(ctx, "global", 30)
		// Method not available - using mock for now
		insights["price_trends"] = map[string]interface{}{
			"message": "Price trend analysis not yet implemented",
		}

		// popular, err := h.services.BigQuery.GetTrendingDestinations(ctx, 5)
		// Method not available - using mock for now
		insights["trending_destinations"] = []string{
			"Paris", "Tokyo", "Bali", "London", "New York",
		}
	}

	// Get user insights from Firebase
	if h.services.Firebase != nil && userID != "" {
		userTrips, err := h.services.Firebase.GetUserTrips(ctx, userID)
		if err == nil {
			insights["user_stats"] = map[string]interface{}{
				"total_trips":     len(userTrips),
				"total_countries": h.countCountries(userTrips),
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"insights":     insights,
		"generated_at": time.Now(),
	})
}

// Helper functions

func (h *AITripHandler) createBasicItinerary(req PlanTripRequest) map[string]interface{} {
	days := h.calculateDays(req.StartDate, req.EndDate)
	itinerary := make(map[string]interface{})

	for i := 1; i <= days; i++ {
		dayKey := fmt.Sprintf("day_%d", i)
		itinerary[dayKey] = map[string]interface{}{
			"morning":   fmt.Sprintf("Explore %s attractions", req.Destination),
			"afternoon": "Local dining and shopping",
			"evening":   "Relaxation and local entertainment",
		}
	}

	return itinerary
}

func (h *AITripHandler) getBasicSuggestions(destination string) []string {
	return []string{
		fmt.Sprintf("Visit popular attractions in %s", destination),
		"Try local cuisine and street food",
		"Take a guided city tour",
		"Visit local markets and shops",
		"Experience nightlife and entertainment",
	}
}

func (h *AITripHandler) calculateBudgetBreakdown(total float64, travelers int) TripBudget {
	perPerson := total / float64(travelers)

	return TripBudget{
		Total:          total,
		Accommodation:  total * 0.4,
		Transportation: total * 0.25,
		Food:           total * 0.2,
		Activities:     total * 0.15,
		Breakdown: map[string]float64{
			"per_person":     perPerson,
			"accommodation":  total * 0.4,
			"transportation": total * 0.25,
			"food":           total * 0.2,
			"activities":     total * 0.15,
		},
	}
}

func (h *AITripHandler) calculateDays(startDate, endDate string) int {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	return int(end.Sub(start).Hours()/24) + 1
}

func (h *AITripHandler) parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return date
}

func (h *AITripHandler) parseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return val
}

func (h *AITripHandler) getDefaultRecommendations() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"destination": "Paris, France",
			"type":        "cultural",
			"budget":      "medium",
			"rating":      4.5,
		},
		{
			"destination": "Tokyo, Japan",
			"type":        "adventure",
			"budget":      "high",
			"rating":      4.7,
		},
		{
			"destination": "Bali, Indonesia",
			"type":        "relaxation",
			"budget":      "low",
			"rating":      4.4,
		},
	}
}

func (h *AITripHandler) countCountries(trips []services.TripData) int {
	countries := make(map[string]bool)
	for _, trip := range trips {
		countries[trip.Destination] = true
	}
	return len(countries)
}
