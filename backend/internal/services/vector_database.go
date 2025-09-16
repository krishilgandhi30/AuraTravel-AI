package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

// VectorDatabase handles embedding storage and similarity search
type VectorDatabase struct {
	firestore *firestore.Client
	gemini    *GeminiService
}

// NewVectorDatabase creates a new vector database instance
func NewVectorDatabase(firestoreClient *firestore.Client, gemini *GeminiService) *VectorDatabase {
	return &VectorDatabase{
		firestore: firestoreClient,
		gemini:    gemini,
	}
}

// EmbeddingDocument represents a document with embeddings
type EmbeddingDocument struct {
	ID        string                 `firestore:"id" json:"id"`
	Type      string                 `firestore:"type" json:"type"` // attraction, trip, preference, user_profile
	Content   string                 `firestore:"content" json:"content"`
	Metadata  map[string]interface{} `firestore:"metadata" json:"metadata"`
	Embedding []float64              `firestore:"embedding" json:"embedding"`
	CreatedAt time.Time              `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time              `firestore:"updated_at" json:"updated_at"`
}

// SimilarityResult represents a similarity search result
type SimilarityResult struct {
	Document   EmbeddingDocument `json:"document"`
	Similarity float64           `json:"similarity"`
}

// StoreEmbedding stores a document with its embedding
func (vdb *VectorDatabase) StoreEmbedding(ctx context.Context, doc EmbeddingDocument) error {
	if doc.ID == "" {
		return fmt.Errorf("document ID is required")
	}

	// Generate embedding if not provided
	if len(doc.Embedding) == 0 && doc.Content != "" {
		embedding, err := vdb.generateEmbedding(ctx, doc.Content)
		if err != nil {
			log.Printf("Failed to generate embedding: %v", err)
			// Continue without embedding - store document anyway
		} else {
			doc.Embedding = embedding
		}
	}

	doc.UpdatedAt = time.Now()
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}

	collection := vdb.getCollectionName(doc.Type)
	_, err := vdb.firestore.Collection(collection).Doc(doc.ID).Set(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to store embedding document: %v", err)
	}

	log.Printf("Stored embedding document: %s in collection: %s", doc.ID, collection)
	return nil
}

// SearchSimilar finds similar documents using cosine similarity
func (vdb *VectorDatabase) SearchSimilar(ctx context.Context, query string, docType string, limit int) ([]SimilarityResult, error) {
	// Generate embedding for query
	queryEmbedding, err := vdb.generateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %v", err)
	}

	// Retrieve all documents of the specified type
	collection := vdb.getCollectionName(docType)
	docs, err := vdb.firestore.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}

	var results []SimilarityResult
	for _, doc := range docs {
		var embeddingDoc EmbeddingDocument
		if err := doc.DataTo(&embeddingDoc); err != nil {
			log.Printf("Failed to unmarshal document %s: %v", doc.Ref.ID, err)
			continue
		}

		// Skip documents without embeddings
		if len(embeddingDoc.Embedding) == 0 {
			continue
		}

		// Calculate cosine similarity
		similarity := vdb.cosineSimilarity(queryEmbedding, embeddingDoc.Embedding)
		results = append(results, SimilarityResult{
			Document:   embeddingDoc,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit results
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// StoreAttractionEmbedding stores an attraction with embedding
func (vdb *VectorDatabase) StoreAttractionEmbedding(ctx context.Context, attraction Attraction) error {
	content := fmt.Sprintf("%s %s %s %s",
		attraction.Name,
		attraction.Type,
		attraction.Description,
		strings.Join(attraction.Tags, " "))

	metadata := map[string]interface{}{
		"name":        attraction.Name,
		"type":        attraction.Type,
		"rating":      attraction.Rating,
		"price_level": attraction.PriceLevel,
		"location":    attraction.Location,
		"tags":        attraction.Tags,
		"available":   attraction.Available,
	}

	doc := EmbeddingDocument{
		ID:       attraction.ID,
		Type:     "attraction",
		Content:  content,
		Metadata: metadata,
	}

	return vdb.StoreEmbedding(ctx, doc)
}

// StoreTripEmbedding stores a trip with embedding
func (vdb *VectorDatabase) StoreTripEmbedding(ctx context.Context, trip TripData) error {
	// Create content from trip details
	itineraryStr, _ := json.Marshal(trip.Itinerary)
	content := fmt.Sprintf("%s %s %s", trip.Title, trip.Destination, string(itineraryStr))

	metadata := map[string]interface{}{
		"user_id":     trip.UserID,
		"destination": trip.Destination,
		"budget":      trip.Budget,
		"travelers":   trip.Travelers,
		"status":      trip.Status,
		"start_date":  trip.StartDate,
		"end_date":    trip.EndDate,
	}

	doc := EmbeddingDocument{
		ID:       trip.ID,
		Type:     "trip",
		Content:  content,
		Metadata: metadata,
	}

	return vdb.StoreEmbedding(ctx, doc)
}

// StoreUserPreferencesEmbedding stores user preferences with embedding
func (vdb *VectorDatabase) StoreUserPreferencesEmbedding(ctx context.Context, userProfile UserProfile) error {
	// Create content from user preferences
	preferencesStr, _ := json.Marshal(userProfile.TravelPreferences)
	content := fmt.Sprintf("%s %s", userProfile.DisplayName, string(preferencesStr))

	metadata := map[string]interface{}{
		"user_id":            userProfile.UID,
		"email":              userProfile.Email,
		"travel_preferences": userProfile.TravelPreferences,
		"trip_history":       userProfile.TripHistory,
	}

	doc := EmbeddingDocument{
		ID:       userProfile.UID,
		Type:     "user_profile",
		Content:  content,
		Metadata: metadata,
	}

	return vdb.StoreEmbedding(ctx, doc)
}

// FindSimilarAttractions finds attractions similar to the given interests
func (vdb *VectorDatabase) FindSimilarAttractions(ctx context.Context, interests []string, limit int) ([]Attraction, error) {
	query := strings.Join(interests, " ")
	results, err := vdb.SearchSimilar(ctx, query, "attraction", limit)
	if err != nil {
		return nil, err
	}

	var attractions []Attraction
	for _, result := range results {
		// Convert metadata back to Attraction
		attraction := Attraction{
			ID:          result.Document.ID,
			Name:        getStringFromMetadata(result.Document.Metadata, "name"),
			Type:        getStringFromMetadata(result.Document.Metadata, "type"),
			Rating:      getFloatFromMetadata(result.Document.Metadata, "rating"),
			PriceLevel:  getIntFromMetadata(result.Document.Metadata, "price_level"),
			Description: result.Document.Content,
			Available:   getBoolFromMetadata(result.Document.Metadata, "available"),
		}

		// Extract tags
		if tagsInterface, ok := result.Document.Metadata["tags"]; ok {
			if tags, ok := tagsInterface.([]interface{}); ok {
				for _, tag := range tags {
					if tagStr, ok := tag.(string); ok {
						attraction.Tags = append(attraction.Tags, tagStr)
					}
				}
			}
		}

		// Extract location
		if locationInterface, ok := result.Document.Metadata["location"]; ok {
			if locationMap, ok := locationInterface.(map[string]interface{}); ok {
				attraction.Location = Location{
					Latitude:  getFloatFromMetadata(locationMap, "latitude"),
					Longitude: getFloatFromMetadata(locationMap, "longitude"),
					Address:   getStringFromMetadata(locationMap, "address"),
				}
			}
		}

		attractions = append(attractions, attraction)
	}

	return attractions, nil
}

// FindSimilarTrips finds trips similar to the given destination and preferences
func (vdb *VectorDatabase) FindSimilarTrips(ctx context.Context, destination string, preferences map[string]interface{}, limit int) ([]TripData, error) {
	preferencesStr, _ := json.Marshal(preferences)
	query := fmt.Sprintf("%s %s", destination, string(preferencesStr))

	results, err := vdb.SearchSimilar(ctx, query, "trip", limit)
	if err != nil {
		return nil, err
	}

	var trips []TripData
	for _, result := range results {
		trip := TripData{
			ID:          result.Document.ID,
			UserID:      getStringFromMetadata(result.Document.Metadata, "user_id"),
			Title:       getStringFromMetadata(result.Document.Metadata, "title"),
			Destination: getStringFromMetadata(result.Document.Metadata, "destination"),
			Budget:      getFloatFromMetadata(result.Document.Metadata, "budget"),
			Travelers:   getIntFromMetadata(result.Document.Metadata, "travelers"),
			Status:      getStringFromMetadata(result.Document.Metadata, "status"),
		}

		// Parse dates
		if startDateInterface, ok := result.Document.Metadata["start_date"]; ok {
			if startDate, ok := startDateInterface.(time.Time); ok {
				trip.StartDate = startDate
			}
		}

		if endDateInterface, ok := result.Document.Metadata["end_date"]; ok {
			if endDate, ok := endDateInterface.(time.Time); ok {
				trip.EndDate = endDate
			}
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

// generateEmbedding generates an embedding for the given text
func (vdb *VectorDatabase) generateEmbedding(ctx context.Context, text string) ([]float64, error) {
	if vdb.gemini == nil {
		// Return mock embedding if Gemini is not available
		return vdb.generateMockEmbedding(text), nil
	}

	// For now, use a simple mock embedding
	// In production, you would call Vertex AI Embedding API or similar
	return vdb.generateMockEmbedding(text), nil
}

// generateMockEmbedding creates a simple mock embedding based on text
func (vdb *VectorDatabase) generateMockEmbedding(text string) []float64 {
	// Simple hash-based embedding (not suitable for production)
	embedding := make([]float64, 128) // 128-dimensional embedding

	words := strings.Fields(strings.ToLower(text))
	for i, word := range words {
		if i >= len(embedding) {
			break
		}
		// Simple hash function to generate embedding values
		hash := 0
		for _, char := range word {
			hash = hash*31 + int(char)
		}
		embedding[i] = float64(hash%1000) / 1000.0 // Normalize to [0, 1]
	}

	// Normalize the embedding vector
	return vdb.normalizeVector(embedding)
}

// cosineSimilarity calculates cosine similarity between two vectors
func (vdb *VectorDatabase) cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0.0 || normB == 0.0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// normalizeVector normalizes a vector to unit length
func (vdb *VectorDatabase) normalizeVector(vector []float64) []float64 {
	var norm float64
	for _, val := range vector {
		norm += val * val
	}
	norm = math.Sqrt(norm)

	if norm == 0.0 {
		return vector
	}

	normalized := make([]float64, len(vector))
	for i, val := range vector {
		normalized[i] = val / norm
	}
	return normalized
}

// getCollectionName returns the Firestore collection name for the document type
func (vdb *VectorDatabase) getCollectionName(docType string) string {
	return fmt.Sprintf("embeddings_%s", docType)
}

// Helper functions to extract values from metadata
func getStringFromMetadata(metadata map[string]interface{}, key string) string {
	if val, ok := metadata[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getFloatFromMetadata(metadata map[string]interface{}, key string) float64 {
	if val, ok := metadata[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

func getIntFromMetadata(metadata map[string]interface{}, key string) int {
	if val, ok := metadata[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func getBoolFromMetadata(metadata map[string]interface{}, key string) bool {
	if val, ok := metadata[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
