package routes

import (
	"auratravel-backend/internal/handlers"
	"auratravel-backend/internal/middleware"
	"auratravel-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, services *services.Services) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(services)
	userHandler := handlers.NewUserHandler(services)
	tripHandler := handlers.NewTripHandler(services)
	aiTripHandler := handlers.NewAITripHandler(services)
	vectorHandler := handlers.NewVectorHandler(services)

	// Public routes
	public := router.Group("/api/v1")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "AuraTravel API is running",
				"services": gin.H{
					"gemini":   services.Gemini != nil,
					"vertex":   services.Vertex != nil,
					"vision":   services.Vision != nil,
					"bigquery": services.BigQuery != nil,
					"firebase": services.Firebase != nil,
				},
			})
		})

		// Authentication routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/firebase-auth", authHandler.FirebaseAuth)
		}

		// Public AI endpoints (limited functionality)
		public.GET("/recommendations", aiTripHandler.GetRecommendations)
		public.GET("/insights", aiTripHandler.GetTravelInsights)
		public.POST("/analyze-image", aiTripHandler.AnalyzeImage)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile routes
		users := protected.Group("/users")
		{
			users.POST("/register", userHandler.RegisterUser)
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.DELETE("/profile", userHandler.DeleteAccount)
		}

		// Trip management routes
		trips := protected.Group("/trips")
		{
			trips.POST("/", tripHandler.CreateTrip)
			trips.GET("/", tripHandler.GetTrips)
			trips.GET("/:id", tripHandler.GetTrip)
			trips.PUT("/:id", tripHandler.UpdateTrip)
			trips.DELETE("/:id", tripHandler.DeleteTrip)
			trips.GET("/recommendations", tripHandler.GenerateRecommendations)
		}

		// AI-powered trip routes
		aiTrips := protected.Group("/ai")
		{
			aiTrips.POST("/plan-trip", aiTripHandler.PlanTrip)
			aiTrips.GET("/recommendations", aiTripHandler.GetRecommendations)
			aiTrips.POST("/optimize/:id", aiTripHandler.OptimizeItinerary)
			aiTrips.POST("/analyze-image", aiTripHandler.AnalyzeImage)
			aiTrips.GET("/insights", aiTripHandler.GetTravelInsights)
			aiTrips.GET("/rag-context", vectorHandler.GetRAGContext)
			aiTrips.POST("/validate-availability", vectorHandler.ValidateAvailability)
		}

		// Vector database routes
		vector := protected.Group("/vector")
		{
			vector.POST("/search-attractions", vectorHandler.SearchSimilarAttractions)
			vector.POST("/search-trips", vectorHandler.SearchSimilarTrips)
			vector.POST("/store-preferences", vectorHandler.StoreUserPreferences)
			vector.GET("/predict-cost", vectorHandler.PredictTravelCost)
		}

		// Booking routes (future implementation)
		bookings := protected.Group("/bookings")
		{
			bookings.POST("/", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Booking endpoint - coming soon"})
			})
			bookings.GET("/", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Get bookings endpoint - coming soon"})
			})
		}

		// Review routes (future implementation)
		reviews := protected.Group("/reviews")
		{
			reviews.POST("/", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Create review endpoint - coming soon"})
			})
			reviews.GET("/destination/:destination", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Get destination reviews endpoint - coming soon"})
			})
		}
	}
}
