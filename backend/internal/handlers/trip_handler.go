package handlers

import (
	"log"
	"net/http"
	"time"

	"auratravel-backend/internal/models"
	"auratravel-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// TripHandler handles trip operations
type TripHandler struct {
	services *services.Services
}

// NewTripHandler creates a new trip handler
func NewTripHandler(services *services.Services) *TripHandler {
	return &TripHandler{
		services: services,
	}
}

// CreateTripRequest represents the request payload for creating a trip
type CreateTripRequest struct {
	Destination string    `json:"destination" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	TotalBudget float64   `json:"total_budget" binding:"required"`
	Travelers   int       `json:"travelers" binding:"required,min=1"`
}

// CreateTrip creates a new trip with AI-powered itinerary generation
func (h *TripHandler) CreateTrip(c *gin.Context) {
	// (userID, exists) removed: mock data does not use them

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Create trip model (mock fields only)
	trip := &models.Trip{
		ID:          time.Now().Format("20060102150405"),
		UserID:      "mock-user-id",
		Destination: req.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalBudget: req.TotalBudget,
		Travelers:   req.Travelers,
		Status:      "planning",
	}

	// Generate AI-powered itinerary using Gemini (mock)
	itinerary, err := h.services.Gemini.GenerateItinerary(c, services.ItineraryRequest{
		Destination: req.Destination,
		StartDate:   req.StartDate.Format("2006-01-02"),
		EndDate:     req.EndDate.Format("2006-01-02"),
		Budget:      req.TotalBudget,
		Travelers:   req.Travelers,
	})
	if err != nil {
		log.Printf("Failed to generate itinerary: %v", err)
	}

	// VertexAI and BigQuery mock responses (since methods do not exist)
	insights := map[string]interface{}{"insight": "Mock destination insights"}
	budgetAnalysis := map[string]interface{}{"budget": "Mock budget analysis"}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Trip created successfully",
		"trip":            trip,
		"itinerary":       itinerary,
		"insights":        insights,
		"budget_analysis": budgetAnalysis,
	})
}

// GetTrips gets user trips
func (h *TripHandler) GetTrips(c *gin.Context) {
	// (userID, exists) removed: mock data does not use them

	// Mock trips data - replace with actual database query
	trips := []models.Trip{
		{
			ID:          "1",
			UserID:      "mock-user-id",
			Destination: "Paris, France",
			StartDate:   time.Now().AddDate(0, 1, 0),
			EndDate:     time.Now().AddDate(0, 1, 7),
			TotalBudget: 3000.0,
			Travelers:   2,
			Status:      "confirmed",
			CreatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			ID:          "2",
			UserID:      "mock-user-id",
			Destination: "Tokyo, Japan",
			StartDate:   time.Now().AddDate(0, 2, 0),
			EndDate:     time.Now().AddDate(0, 2, 10),
			TotalBudget: 4500.0,
			Travelers:   1,
			Status:      "planning",
			CreatedAt:   time.Now().AddDate(0, 0, -2),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"trips": trips,
		"total": len(trips),
	})
}

// GetTrip gets a specific trip with detailed itinerary
func (h *TripHandler) GetTrip(c *gin.Context) {
	// (userID, exists, id) removed: mock data does not use them
	tripID := c.Param("id")

	// Mock trip data - replace with actual database query
	trip := &models.Trip{
		ID:          tripID,
		UserID:      "mock-user-id",
		Destination: "Paris, France",
		StartDate:   time.Now().AddDate(0, 1, 0),
		EndDate:     time.Now().AddDate(0, 1, 7),
		TotalBudget: 3000.0,
		Travelers:   2,
		Status:      "confirmed",
		CreatedAt:   time.Now().AddDate(0, 0, -7),
	}

	// Mock recommendations and visual insights
	recommendations := []string{"Mock recommendation 1", "Mock recommendation 2"}
	visualInsights := map[string]interface{}{"insight": "Mock visual insight"}

	c.JSON(http.StatusOK, gin.H{
		"trip":            trip,
		"recommendations": recommendations,
		"visual_insights": visualInsights,
	})
}

// UpdateTrip updates a trip
func (h *TripHandler) UpdateTrip(c *gin.Context) {
	// (userID, exists, id) removed: mock data does not use them
	tripID := c.Param("id")

	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock update - replace with actual database update
	updatedTrip := &models.Trip{
		ID:          tripID,
		UserID:      "mock-user-id",
		Destination: req.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalBudget: req.TotalBudget,
		Travelers:   req.Travelers,
		Status:      "updated",
		UpdatedAt:   time.Now(),
	}

	// Regenerate itinerary with updated preferences (mock)
	newItinerary, err := h.services.Gemini.GenerateItinerary(c, services.ItineraryRequest{
		Destination: req.Destination,
		StartDate:   req.StartDate.Format("2006-01-02"),
		EndDate:     req.EndDate.Format("2006-01-02"),
		Budget:      req.TotalBudget,
		Travelers:   req.Travelers,
	})
	if err != nil {
		log.Printf("Failed to regenerate itinerary: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Trip updated successfully",
		"trip":          updatedTrip,
		"new_itinerary": newItinerary,
	})
}

// DeleteTrip deletes a trip
func (h *TripHandler) DeleteTrip(c *gin.Context) {
	// (userID, exists, id) removed: mock data does not use them
	tripID := c.Param("id")
	// Mock deletion - replace with actual database deletion
	log.Printf("Deleting trip %s", tripID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Trip deleted successfully",
		"trip_id": tripID,
	})
}

// GenerateRecommendations generates AI-powered recommendations for a destination
func (h *TripHandler) GenerateRecommendations(c *gin.Context) {
	destination := c.Query("destination")
	tripType := c.Query("trip_type")

	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination is required"})
		return
	}

	// Mock recommendations and visual insights
	geminiRecs := []string{"Mock Gemini recommendation 1", "Mock Gemini recommendation 2"}
	vertexRecs := []string{"Mock Vertex recommendation 1", "Mock Vertex recommendation 2"}
	visualInsights := map[string]interface{}{"insight": "Mock visual insight"}

	c.JSON(http.StatusOK, gin.H{
		"destination":            destination,
		"trip_type":              tripType,
		"gemini_recommendations": geminiRecs,
		"vertex_recommendations": vertexRecs,
		"visual_insights":        visualInsights,
	})
}
