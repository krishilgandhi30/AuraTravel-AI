package services

import (
	"context"
	"log"
)

// Services contains all service dependencies
type Services struct {
	Gemini   *GeminiService
	Vertex   *VertexService
	Vision   *VisionService
	BigQuery *BigQueryService
	Firebase *FirebaseService
}

// NewServices initializes and returns all services
func NewServices() (*Services, error) {
	// Initialize Google AI services
	geminiService, err := NewGeminiService()
	if err != nil {
		log.Printf("Warning: Failed to initialize Gemini service: %v", err)
	}

	vertexService, err := NewVertexService()
	if err != nil {
		log.Printf("Warning: Failed to initialize Vertex AI service: %v", err)
	}

	visionService, err := NewVisionService()
	if err != nil {
		log.Printf("Warning: Failed to initialize Vision service: %v", err)
	}

	bigQueryService, err := NewBigQueryService()
	if err != nil {
		log.Printf("Warning: Failed to initialize BigQuery service: %v", err)
	}

	firebaseService, err := NewFirebaseService()
	if err != nil {
		log.Printf("Warning: Failed to initialize Firebase service: %v", err)
	}

	log.Println("Services initialized successfully")

	return &Services{
		Gemini:   geminiService,
		Vertex:   vertexService,
		Vision:   visionService,
		BigQuery: bigQueryService,
		Firebase: firebaseService,
	}, nil
}

// Shutdown gracefully shuts down all services
func (s *Services) Shutdown(ctx context.Context) error {
	var lastError error

	// Shutdown AI services
	if s.Gemini != nil {
		if err := s.Gemini.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down Gemini service: %v", err)
			lastError = err
		}
	}

	if s.Vertex != nil {
		if err := s.Vertex.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down Vertex service: %v", err)
			lastError = err
		}
	}

	if s.Vision != nil {
		if err := s.Vision.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down Vision service: %v", err)
			lastError = err
		}
	}

	if s.BigQuery != nil {
		if err := s.BigQuery.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down BigQuery service: %v", err)
			lastError = err
		}
	}

	if s.Firebase != nil {
		if err := s.Firebase.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down Firebase service: %v", err)
			lastError = err
		}
	}

	log.Println("All services shut down")
	return lastError
}
