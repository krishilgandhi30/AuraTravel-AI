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

	trip := &models.Trip{
		ID:          time.Now().Format("20060102150405"),
		UserID:      "mock-user-id", // TODO: Replace with actual user ID from context
		Destination: req.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalBudget: req.TotalBudget,
		Travelers:   req.Travelers,
		Status:      "planning",
		CreatedAt:   time.Now(),
	}

	tripData := services.TripData{
		ID:          trip.ID,
		UserID:      trip.UserID,
		Destination: trip.Destination,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Status:      trip.Status,
		Budget:      trip.TotalBudget,
		Travelers:   trip.Travelers,
		CreatedAt:   trip.CreatedAt,
	}
	fb := h.services.Firebase
	ctx := c.Request.Context()
	if err := fb.SaveTrip(ctx, tripData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip in Firestore"})
		return
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

	fb := h.services.Firebase
	ctx := c.Request.Context()
	userID := "mock-user-id" // TODO: Replace with actual user ID from context
	tripDatas, err := fb.GetUserTrips(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trips from Firestore"})
		return
	}
	trips := make([]models.Trip, 0, len(tripDatas))
	for _, td := range tripDatas {
		trips = append(trips, models.Trip{
			ID:          td.ID,
			UserID:      td.UserID,
			Destination: td.Destination,
			StartDate:   toTime(td.StartDate),
			EndDate:     toTime(td.EndDate),
			Status:      td.Status,
			TotalBudget: td.Budget,
			Travelers:   td.Travelers,
			CreatedAt:   toTime(td.CreatedAt),
			UpdatedAt:   toTime(td.UpdatedAt),
		})
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

	fb := h.services.Firebase
	ctx := c.Request.Context()
	td, err := fb.GetTrip(ctx, tripID)
	if err != nil || td == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	trip := models.Trip{
		ID:          td.ID,
		UserID:      td.UserID,
		Destination: td.Destination,
		StartDate:   toTime(td.StartDate),
		EndDate:     toTime(td.EndDate),
		Status:      td.Status,
		TotalBudget: td.Budget,
		Travelers:   td.Travelers,
		CreatedAt:   toTime(td.CreatedAt),
		UpdatedAt:   toTime(td.UpdatedAt),
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

	fb := h.services.Firebase
	ctx := c.Request.Context()
	updates := map[string]interface{}{
		"destination": req.Destination,
		"start_date":  req.StartDate,
		"end_date":    req.EndDate,
		"budget":      req.TotalBudget,
		"travelers":   req.Travelers,
		"status":      "updated",
		"updated_at":  time.Now(),
	}
	if err := fb.UpdateTrip(ctx, tripID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip in Firestore"})
		return
	}
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
		"trip":          updates,
		"new_itinerary": newItinerary,
	})
}

// DeleteTrip deletes a trip
func (h *TripHandler) DeleteTrip(c *gin.Context) {
	// (userID, exists, id) removed: mock data does not use them
	tripID := c.Param("id")
	fb := h.services.Firebase
	ctx := c.Request.Context()
	if err := fb.DeleteTrip(ctx, tripID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip from Firestore"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Trip deleted successfully",
		"trip_id": tripID,
	})
}

// toTime safely converts Firestore timestamp/interface{} to time.Time
func toTime(val interface{}) time.Time {
	switch t := val.(type) {
	case time.Time:
		return t
	case *time.Time:
		if t != nil {
			return *t
		}
	}
	return time.Time{}
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
