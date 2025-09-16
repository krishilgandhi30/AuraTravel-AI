package main

import (
	"log"
	"os"

	"auratravel-backend/internal/config"
	"auratravel-backend/internal/database"
	"auratravel-backend/internal/routes"
	"auratravel-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AuraTravel AI Backend API
// @version 1.0
// @description AI-powered travel planning platform backend
// @termsOfService http://swagger.io/terms/

// @contact.name AuraTravel Team
// @contact.url http://www.auratravel.ai/support
// @contact.email support@auratravel.ai

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg := config.GetConfig()

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Server will continue without database connection...")
	}

	// Initialize services (for future use)
	services, err := services.NewServices()
	if err != nil {
		log.Printf("Warning: Failed to initialize services: %v", err)
	}

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"https://auratravel.ai",
		"https://*.auratravel.ai",
	}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
	}
	corsConfig.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}
	corsConfig.AllowCredentials = true

	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auratravel-backend",
			"version": "1.0.0",
		})
	})

	// Setup routes
	routes.SetupRoutes(router, services)

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting AuraTravel AI Backend on port %s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/docs/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
