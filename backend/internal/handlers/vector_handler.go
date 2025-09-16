package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"auratravel-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// VectorHandler handles vector database operations
type VectorHandler struct {
	services *services.Services
}

// NewVectorHandler creates a new vector handler
func NewVectorHandler(services *services.Services) *VectorHandler {
	return &VectorHandler{
		services: services,
	}
}

// SearchSimilarAttractionsRequest represents the request for similarity search
type SearchSimilarAttractionsRequest struct {
	Interests []string `json:"interests" binding:"required"`
	Limit     int      `json:"limit"`
}

// SearchSimilarTripsRequest represents the request for similar trips search
type SearchSimilarTripsRequest struct {
	Destination string                 `json:"destination" binding:"required"`
	Preferences map[string]interface{} `json:"preferences" binding:"required"`
	Limit       int                    `json:"limit"`
}

// StoreUserPreferencesRequest represents the request for storing user preferences
type StoreUserPreferencesRequest struct {
	UserID      string                 `json:"user_id" binding:"required"`
	Preferences map[string]interface{} `json:"preferences" binding:"required"`
	TripHistory []string               `json:"trip_history"`
}

// SearchSimilarAttractions searches for similar attractions using vector similarity
func (h *VectorHandler) SearchSimilarAttractions(c *gin.Context) {
	var req SearchSimilarAttractionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	ctx := context.Background()

	if h.services.VectorDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Vector database not available"})
		return
	}

	attractions, err := h.services.VectorDB.FindSimilarAttractions(ctx, req.Interests, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search attractions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attractions": attractions,
		"count":       len(attractions),
		"query":       req.Interests,
	})
}

// SearchSimilarTrips searches for similar trips using vector similarity
func (h *VectorHandler) SearchSimilarTrips(c *gin.Context) {
	var req SearchSimilarTripsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 5
	}

	ctx := context.Background()

	if h.services.VectorDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Vector database not available"})
		return
	}

	trips, err := h.services.VectorDB.FindSimilarTrips(ctx, req.Destination, req.Preferences, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search trips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trips":       trips,
		"count":       len(trips),
		"destination": req.Destination,
	})
}

// StoreUserPreferences stores user preferences in vector database
func (h *VectorHandler) StoreUserPreferences(c *gin.Context) {
	var req StoreUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	if h.services.VectorDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Vector database not available"})
		return
	}

	// Create user profile for embedding
	userProfile := services.UserProfile{
		UID:               req.UserID,
		Email:             "", // Would be filled from auth context in production
		DisplayName:       req.UserID,
		TravelPreferences: req.Preferences,
		TripHistory:       req.TripHistory,
	}

	err := h.services.VectorDB.StoreUserPreferencesEmbedding(ctx, userProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User preferences stored successfully",
		"user_id": req.UserID,
	})
}

// GetRAGContext retrieves RAG context for a destination
func (h *VectorHandler) GetRAGContext(c *gin.Context) {
	destination := c.Query("destination")
	userID := c.Query("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	budgetStr := c.Query("budget")
	travelersStr := c.Query("travelers")

	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination is required"})
		return
	}

	// Parse optional parameters
	budget := 1000.0
	if budgetStr != "" {
		if b, err := strconv.ParseFloat(budgetStr, 64); err == nil {
			budget = b
		}
	}

	travelers := 1
	if travelersStr != "" {
		if t, err := strconv.Atoi(travelersStr); err == nil {
			travelers = t
		}
	}

	interests := c.QueryArray("interests")
	preferencesStr := c.Query("preferences")
	var preferences map[string]interface{}
	if preferencesStr != "" {
		// In a real implementation, you'd parse this properly
		preferences = map[string]interface{}{"raw": preferencesStr}
	}

	ctx := context.Background()

	if h.services.RAGRetriever == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "RAG retriever not available"})
		return
	}

	// Build retrieval request
		layout := "2006-01-02"

		start, _ := time.Parse(layout, h.parseDate(startDate))
		end, _ := time.Parse(layout, h.parseDate(endDate))
		ragRequest := services.RetrievalRequest{
			UserID:      userID,
			Destination: destination,
			StartDate:   start,
			EndDate:     end,
			Budget:      budget,
			Travelers:   travelers,
			Interests:   interests,
			Preferences: preferences,
		}

	ragContext, err := h.services.RAGRetriever.RetrieveContext(ctx, ragRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve RAG context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"context":     ragContext,
		"destination": destination,
		"user_id":     userID,
	})
}

// ValidateAvailability validates real-time availability of items
func (h *VectorHandler) ValidateAvailability(c *gin.Context) {
	var req struct {
		Type      string        `json:"type" binding:"required"` // attractions, hotels, transportation
		Data      []interface{} `json:"data" binding:"required"`
		CheckDate string        `json:"check_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ctx := context.Background()

	if h.services.RAGRetriever == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "RAG retriever not available"})
		return
	}

	// Validate availability using the data validator
	validatedData := req.Data         // In real implementation, would validate each item
	availabilityStatus := "available" // Mock status

	c.JSON(http.StatusOK, gin.H{
		"validated_data": validatedData,
		"status":         availabilityStatus,
		"check_date":     req.CheckDate,
		"type":           req.Type,
	})
}

// PredictTravelCost predicts travel costs using lightweight ML model
func (h *VectorHandler) PredictTravelCost(c *gin.Context) {
	destination := c.Query("destination")
	travelDateStr := c.Query("travel_date")
	durationStr := c.Query("duration")
	travelersStr := c.Query("travelers")
	budgetPreference := c.Query("budget_preference")

	if destination == "" || travelDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination and travel_date are required"})
		return
	}

	// Parse parameters
	travelDate, err := time.Parse("2006-01-02", travelDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid travel_date format. Use YYYY-MM-DD"})
		return
	}

	duration := 7 // Default 7 days
	if durationStr != "" {
		if d, err := strconv.Atoi(durationStr); err == nil {
			duration = d
		}
	}

	travelers := 1 // Default 1 traveler
	if travelersStr != "" {
		if t, err := strconv.Atoi(travelersStr); err == nil {
			travelers = t
		}
	}

	if budgetPreference == "" {
		budgetPreference = "mid-range"
	}

	ctx := context.Background()

	if h.services.CostPredictor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cost predictor not available"})
		return
	}

	// Create prediction request
	predictionReq := services.CostPredictionRequest{
		Destination:      destination,
		TravelDate:       travelDate,
		Duration:         duration,
		Travelers:        travelers,
		BudgetPreference: budgetPreference,
	}

	prediction, err := h.services.CostPredictor.PredictTravelCost(ctx, predictionReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to predict travel cost"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prediction": prediction,
		"request":    predictionReq,
	})
}

// Helper function to parse date strings
func (h *VectorHandler) parseDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	// In production, you'd want proper date parsing and validation
	return dateStr
}
