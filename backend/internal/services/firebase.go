package services

import (
	"context"
	"fmt"
	"log"

	"auratravel-backend/internal/config"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseService handles Firebase operations
type FirebaseService struct {
	app       *firebase.App
	auth      *auth.Client
	firestore *firestore.Client
	cfg       *config.Config
}

// NewFirebaseService creates a new Firebase service
func NewFirebaseService() (*FirebaseService, error) {
	cfg := config.GetConfig()

	ctx := context.Background()

	var opts []option.ClientOption
	if cfg.GoogleApplicationCredentials != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.GoogleApplicationCredentials))
	}

	// Initialize Firebase app
	firebaseConfig := &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	// Initialize Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth: %v", err)
	}

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore: %v", err)
	}

	return &FirebaseService{
		app:       app,
		auth:      authClient,
		firestore: firestoreClient,
		cfg:       cfg,
	}, nil
}

// UserProfile represents a user profile in Firestore
type UserProfile struct {
	UID               string                 `firestore:"uid"`
	Email             string                 `firestore:"email"`
	DisplayName       string                 `firestore:"display_name"`
	PhotoURL          string                 `firestore:"photo_url"`
	TravelPreferences map[string]interface{} `firestore:"travel_preferences"`
	TripHistory       []string               `firestore:"trip_history"`
	Recommendations   []string               `firestore:"recommendations"`
	CreatedAt         interface{}            `firestore:"created_at"`
	UpdatedAt         interface{}            `firestore:"updated_at"`
	LastLogin         interface{}            `firestore:"last_login"`
}

// TripData represents trip data stored in Firestore
type TripData struct {
	ID          string                 `firestore:"id"`
	UserID      string                 `firestore:"user_id"`
	Title       string                 `firestore:"title"`
	Destination string                 `firestore:"destination"`
	StartDate   interface{}            `firestore:"start_date"`
	EndDate     interface{}            `firestore:"end_date"`
	Status      string                 `firestore:"status"`
	Itinerary   map[string]interface{} `firestore:"itinerary"`
	Budget      float64                `firestore:"budget"`
	Travelers   int                    `firestore:"travelers"`
	CreatedAt   interface{}            `firestore:"created_at"`
	UpdatedAt   interface{}            `firestore:"updated_at"`
}

// VerifyIDToken verifies Firebase ID token
func (f *FirebaseService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := f.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}
	return token, nil
}

// GetUser gets user information from Firebase Auth
func (f *FirebaseService) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := f.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

// CreateCustomToken creates a custom token for user authentication
func (f *FirebaseService) CreateCustomToken(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	token, err := f.auth.CustomTokenWithClaims(ctx, uid, claims)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token: %v", err)
	}
	return token, nil
}

// SaveUserProfile saves user profile to Firestore
func (f *FirebaseService) SaveUserProfile(ctx context.Context, profile UserProfile) error {
	_, err := f.firestore.Collection("users").Doc(profile.UID).Set(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to save user profile: %v", err)
	}
	log.Printf("Saved user profile for UID: %s", profile.UID)
	return nil
}

// GetUserProfile retrieves user profile from Firestore
func (f *FirebaseService) GetUserProfile(ctx context.Context, uid string) (*UserProfile, error) {
	doc, err := f.firestore.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %v", err)
	}

	var profile UserProfile
	if err := doc.DataTo(&profile); err != nil {
		return nil, fmt.Errorf("failed to convert user profile: %v", err)
	}

	return &profile, nil
}

// UpdateUserPreferences updates user travel preferences
func (f *FirebaseService) UpdateUserPreferences(ctx context.Context, uid string, preferences map[string]interface{}) error {
	_, err := f.firestore.Collection("users").Doc(uid).Update(ctx, []firestore.Update{
		{Path: "travel_preferences", Value: preferences},
		{Path: "updated_at", Value: firestore.ServerTimestamp},
	})
	if err != nil {
		return fmt.Errorf("failed to update user preferences: %v", err)
	}
	return nil
}

// SaveTrip saves trip data to Firestore
func (f *FirebaseService) SaveTrip(ctx context.Context, trip TripData) error {
	_, err := f.firestore.Collection("trips").Doc(trip.ID).Set(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to save trip: %v", err)
	}
	log.Printf("Saved trip: %s for user: %s", trip.ID, trip.UserID)
	return nil
}

// GetUserTrips retrieves all trips for a user
func (f *FirebaseService) GetUserTrips(ctx context.Context, userID string) ([]TripData, error) {
	iter := f.firestore.Collection("trips").Where("user_id", "==", userID).Documents(ctx)
	defer iter.Stop()

	var trips []TripData
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var trip TripData
		if err := doc.DataTo(&trip); err != nil {
			log.Printf("Error converting trip data: %v", err)
			continue
		}
		trips = append(trips, trip)
	}

	return trips, nil
}

// GetTrip retrieves a specific trip
func (f *FirebaseService) GetTrip(ctx context.Context, tripID string) (*TripData, error) {
	doc, err := f.firestore.Collection("trips").Doc(tripID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %v", err)
	}

	var trip TripData
	if err := doc.DataTo(&trip); err != nil {
		return nil, fmt.Errorf("failed to convert trip data: %v", err)
	}

	return &trip, nil
}

// UpdateTrip updates trip data
func (f *FirebaseService) UpdateTrip(ctx context.Context, tripID string, updates map[string]interface{}) error {
	var firestoreUpdates []firestore.Update
	for key, value := range updates {
		firestoreUpdates = append(firestoreUpdates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}
	firestoreUpdates = append(firestoreUpdates, firestore.Update{
		Path:  "updated_at",
		Value: firestore.ServerTimestamp,
	})

	_, err := f.firestore.Collection("trips").Doc(tripID).Update(ctx, firestoreUpdates)
	if err != nil {
		return fmt.Errorf("failed to update trip: %v", err)
	}
	return nil
}

// DeleteTrip deletes a trip (soft delete by updating status)
func (f *FirebaseService) DeleteTrip(ctx context.Context, tripID string) error {
	_, err := f.firestore.Collection("trips").Doc(tripID).Update(ctx, []firestore.Update{
		{Path: "status", Value: "deleted"},
		{Path: "updated_at", Value: firestore.ServerTimestamp},
	})
	if err != nil {
		return fmt.Errorf("failed to delete trip: %v", err)
	}
	return nil
}

// SaveRecommendations saves AI recommendations to Firestore
func (f *FirebaseService) SaveRecommendations(ctx context.Context, userID string, recommendations []map[string]interface{}) error {
	batch := f.firestore.Batch()

	for i, rec := range recommendations {
		docRef := f.firestore.Collection("recommendations").NewDoc()
		rec["user_id"] = userID
		rec["created_at"] = firestore.ServerTimestamp
		batch.Set(docRef, rec)

		if i >= 500 { // Firestore batch limit
			break
		}
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to save recommendations: %v", err)
	}

	log.Printf("Saved %d recommendations for user: %s", len(recommendations), userID)
	return nil
}

// GetRecommendations retrieves recommendations for a user
func (f *FirebaseService) GetRecommendations(ctx context.Context, userID string, limit int) ([]map[string]interface{}, error) {
	query := f.firestore.Collection("recommendations").
		Where("user_id", "==", userID).
		OrderBy("created_at", firestore.Desc).
		Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var recommendations []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		recommendations = append(recommendations, doc.Data())
	}

	return recommendations, nil
}

// SaveAnalyticsEvent saves analytics events to Firestore
func (f *FirebaseService) SaveAnalyticsEvent(ctx context.Context, event map[string]interface{}) error {
	event["timestamp"] = firestore.ServerTimestamp
	_, _, err := f.firestore.Collection("analytics").Add(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to save analytics event: %v", err)
	}
	return nil
}

// GetPopularDestinations gets popular destinations from Firestore analytics
func (f *FirebaseService) GetPopularDestinations(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	// This would typically involve complex aggregation queries
	// For now, return mock data
	destinations := []map[string]interface{}{
		{
			"destination":  "Paris, France",
			"search_count": 1250,
			"book_count":   890,
			"rating":       4.5,
		},
		{
			"destination":  "Tokyo, Japan",
			"search_count": 1100,
			"book_count":   750,
			"rating":       4.7,
		},
		{
			"destination":  "Bali, Indonesia",
			"search_count": 980,
			"book_count":   650,
			"rating":       4.4,
		},
	}

	return destinations, nil
}

// BackupUserData creates a backup of user data
func (f *FirebaseService) BackupUserData(ctx context.Context, userID string) (map[string]interface{}, error) {
	backup := make(map[string]interface{})

	// Get user profile
	profile, err := f.GetUserProfile(ctx, userID)
	if err == nil {
		backup["profile"] = profile
	}

	// Get user trips
	trips, err := f.GetUserTrips(ctx, userID)
	if err == nil {
		backup["trips"] = trips
	}

	// Get user recommendations
	recommendations, err := f.GetRecommendations(ctx, userID, 100)
	if err == nil {
		backup["recommendations"] = recommendations
	}

	backup["exported_at"] = firestore.ServerTimestamp

	return backup, nil
}

// Shutdown closes the Firebase connections
func (f *FirebaseService) Shutdown(ctx context.Context) error {
	if f.firestore != nil {
		if err := f.firestore.Close(); err != nil {
			return fmt.Errorf("failed to close Firestore client: %v", err)
		}
		log.Println("Firebase service shut down successfully")
	}
	return nil
}
