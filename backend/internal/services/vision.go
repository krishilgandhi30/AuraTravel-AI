package services

import (
	"context"
	"log"

	"auratravel-backend/internal/config"
)

// VisionService handles Google Cloud Vision API interactions
type VisionService struct {
	projectID string
	cfg       *config.Config
}

// NewVisionService creates a new Vision service
func NewVisionService() (*VisionService, error) {
	cfg := config.GetConfig()

	if cfg.GoogleCloudProjectID == "" {
		log.Println("Warning: GOOGLE_PROJECT_ID not set, using mock service")
		return &VisionService{
			projectID: "",
			cfg:       cfg,
		}, nil
	}

	return &VisionService{
		projectID: cfg.GoogleCloudProjectID,
		cfg:       cfg,
	}, nil
}

// AnalyzeImage analyzes an image using Cloud Vision API
func (v *VisionService) AnalyzeImage(ctx context.Context, imageData []byte) (map[string]interface{}, error) {
	if v.projectID == "" {
		return v.mockImageAnalysis(imageData), nil
	}

	// TODO: Implement actual Vision API call
	// For now, return mock data
	return v.mockImageAnalysis(imageData), nil
}

// DetectLandmarks detects landmarks in an image
func (v *VisionService) DetectLandmarks(ctx context.Context, imageData []byte) ([]map[string]interface{}, error) {
	if v.projectID == "" {
		return v.mockLandmarkDetection(imageData), nil
	}

	// TODO: Implement actual Vision API call
	// For now, return mock data
	return v.mockLandmarkDetection(imageData), nil
}

// ExtractText extracts text from an image
func (v *VisionService) ExtractText(ctx context.Context, imageData []byte) (string, error) {
	if v.projectID == "" {
		return v.mockTextExtraction(imageData), nil
	}

	// TODO: Implement actual Vision API call
	// For now, return mock data
	return v.mockTextExtraction(imageData), nil
}

// Mock implementations

func (v *VisionService) mockImageAnalysis(imageData []byte) map[string]interface{} {
	return map[string]interface{}{
		"labels": []map[string]interface{}{
			{
				"description": "Travel",
				"score":       0.95,
				"confidence":  "high",
			},
			{
				"description": "Tourism",
				"score":       0.89,
				"confidence":  "high",
			},
			{
				"description": "Landmark",
				"score":       0.82,
				"confidence":  "medium",
			},
		},
		"dominant_colors": []string{"#3498db", "#e74c3c", "#f39c12"},
		"image_properties": map[string]interface{}{
			"format": "JPEG",
			"width":  1920,
			"height": 1080,
		},
		"safe_search": map[string]string{
			"adult":    "very_unlikely",
			"violence": "very_unlikely",
			"racy":     "unlikely",
		},
		"analyzed_at": "2024-01-01T12:00:00Z",
	}
}

func (v *VisionService) mockLandmarkDetection(imageData []byte) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "Eiffel Tower",
			"confidence":  0.94,
			"description": "Famous iron lattice tower in Paris, France",
			"location": map[string]interface{}{
				"latitude":  48.8584,
				"longitude": 2.2945,
			},
			"bounding_box": map[string]interface{}{
				"vertices": []map[string]int{
					{"x": 100, "y": 50},
					{"x": 800, "y": 50},
					{"x": 800, "y": 900},
					{"x": 100, "y": 900},
				},
			},
		},
		{
			"name":        "Arc de Triomphe",
			"confidence":  0.87,
			"description": "Triumphal arch in Paris, France",
			"location": map[string]interface{}{
				"latitude":  48.8738,
				"longitude": 2.2950,
			},
		},
	}
}

func (v *VisionService) mockTextExtraction(imageData []byte) string {
	return `Welcome to Paris!
Visit the Eiffel Tower
Open: 9:30 AM - 11:45 PM
Admission: €29.40
Metro: Bir-Hakeim (Line 6)
Website: www.toureiffel.paris`
}

// Shutdown closes the Vision service
func (v *VisionService) Shutdown(ctx context.Context) error {
	log.Println("Vision service shut down successfully")
	return nil
}
