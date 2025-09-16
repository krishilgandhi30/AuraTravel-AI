package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"auratravel-backend/internal/config"
)

// EmbeddingService handles text embedding generation
type EmbeddingService struct {
	projectID   string
	location    string
	httpClient  *http.Client
	accessToken string
	cfg         *config.Config
}

// EmbeddingRequest represents a request to the Vertex AI Embeddings API
type EmbeddingRequest struct {
	Instances []EmbeddingInstance `json:"instances"`
}

// EmbeddingInstance represents an instance for embedding generation
type EmbeddingInstance struct {
	Content string `json:"content"`
}

// EmbeddingResponse represents the response from Vertex AI Embeddings API
type EmbeddingResponse struct {
	Predictions []EmbeddingPrediction `json:"predictions"`
}

// EmbeddingPrediction represents a prediction containing embeddings
type EmbeddingPrediction struct {
	Embeddings EmbeddingValues `json:"embeddings"`
}

// EmbeddingValues contains the actual embedding values
type EmbeddingValues struct {
	Values []float64 `json:"values"`
}

// NewEmbeddingService creates a new embedding service
func NewEmbeddingService() (*EmbeddingService, error) {
	cfg := config.GetConfig()

	return &EmbeddingService{
		projectID:  cfg.GoogleCloudProjectID,
		location:   "us-central1", // Default location for Vertex AI
		httpClient: &http.Client{Timeout: 30 * time.Second},
		cfg:        cfg,
	}, nil
}

// GenerateEmbedding generates embeddings for the given text
func (e *EmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if e.projectID == "" {
		log.Println("Project ID not set, using mock embeddings")
		return e.generateMockEmbedding(text), nil
	}

	// Use the textembedding-gecko model for generating embeddings
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/textembedding-gecko:predict",
		e.location, e.projectID, e.location)

	request := EmbeddingRequest{
		Instances: []EmbeddingInstance{
			{Content: text},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if e.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.accessToken))
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to call embedding API, using mock: %v", err)
		return e.generateMockEmbedding(text), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Embedding API failed with status %d: %s, using mock", resp.StatusCode, string(body))
		return e.generateMockEmbedding(text), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var response EmbeddingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Predictions) == 0 {
		return nil, fmt.Errorf("no predictions in response")
	}

	return response.Predictions[0].Embeddings.Values, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (e *EmbeddingService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	var embeddings [][]float64

	// For simplicity, process one by one. In production, you'd want to batch these.
	for _, text := range texts {
		embedding, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			log.Printf("Failed to generate embedding for text, using mock: %v", err)
			embedding = e.generateMockEmbedding(text)
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

// SetAccessToken sets the access token for API authentication
func (e *EmbeddingService) SetAccessToken(token string) {
	e.accessToken = token
}

// generateMockEmbedding creates a deterministic mock embedding based on text
func (e *EmbeddingService) generateMockEmbedding(text string) []float64 {
	// Create a 768-dimensional embedding (same as textembedding-gecko)
	embedding := make([]float64, 768)

	// Simple deterministic hash-based embedding
	hash := 0
	for _, char := range text {
		hash = hash*31 + int(char)
	}

	// Fill embedding with deterministic values based on hash
	for i := 0; i < 768; i++ {
		// Create pseudo-random but deterministic values
		seed := hash + i
		embedding[i] = float64((seed%2000)-1000) / 1000.0 // Range [-1, 1]
	}

	return e.normalizeEmbedding(embedding)
}

// normalizeEmbedding normalizes an embedding vector to unit length
func (e *EmbeddingService) normalizeEmbedding(embedding []float64) []float64 {
	var norm float64
	for _, val := range embedding {
		norm += val * val
	}

	if norm == 0.0 {
		return embedding
	}

	norm = math.Sqrt(norm)
	normalized := make([]float64, len(embedding))
	for i, val := range embedding {
		normalized[i] = val / norm
	}

	return normalized
}
