package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

// RAGRetriever handles retrieval of contextual data for AI generation
type RAGRetriever struct {
	firebase      *FirebaseService
	gemini        *GeminiService
	vision        *VisionService
	dataConnector *DataSourceConnector
	validator     *DataValidator
	mapsAPIKey    string
	weatherKey    string
	httpClient    *http.Client
}

// NewRAGRetriever creates a new RAG retriever instance
func NewRAGRetriever(firebase *FirebaseService, gemini *GeminiService, vision *VisionService, mapsAPIKey, weatherKey string) *RAGRetriever {
	dataConnector := NewDataSourceConnector(mapsAPIKey, weatherKey, "")
	
	retriever := &RAGRetriever{
		firebase:      firebase,
		gemini:        gemini,
		vision:        vision,
		dataConnector: dataConnector,
		mapsAPIKey:    mapsAPIKey,
		weatherKey:    weatherKey,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
	}
	
	// Initialize validator with self-reference
	retriever.validator = NewDataValidator(retriever, dataConnector)
	
	return retriever
}

// TripContext represents the context retrieved for trip planning
type TripContext struct {
	Destination    string            `json:"destination"`
	UserProfile    *UserProfile      `json:"user_profile,omitempty"`
	Attractions    []Attraction      `json:"attractions"`
	Hotels         []Hotel           `json:"hotels"`
	Weather        WeatherForecast   `json:"weather"`
	LocalEvents    []LocalEvent      `json:"local_events"`
	Transportation []TransportOption `json:"transportation"`
	SimilarTrips   []TripData        `json:"similar_trips"`
	EMTInventory   []EMTItem         `json:"emt_inventory"`
}

// Attraction represents a tourist attraction
type Attraction struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"` // museum, restaurant, park, etc.
	Location     Location `json:"location"`
	Rating       float64  `json:"rating"`
	PriceLevel   int      `json:"price_level"` // 0-4
	OpeningHours []string `json:"opening_hours"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	Available    bool     `json:"available"`
}

// Hotel represents accommodation options
type Hotel struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Location      Location `json:"location"`
	Rating        float64  `json:"rating"`
	PricePerNight float64  `json:"price_per_night"`
	Amenities     []string `json:"amenities"`
	Available     bool     `json:"available"`
	BookingURL    string   `json:"booking_url,omitempty"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// WeatherForecast represents weather data
type WeatherForecast struct {
	Current  WeatherCondition   `json:"current"`
	Forecast []WeatherCondition `json:"forecast"` // 7-day forecast
}

// WeatherCondition represents weather at a specific time
type WeatherCondition struct {
	Date        time.Time `json:"date"`
	Temperature float64   `json:"temperature"` // Celsius
	Description string    `json:"description"`
	Humidity    int       `json:"humidity"`
	WindSpeed   float64   `json:"wind_speed"`
	Icon        string    `json:"icon"`
}

// LocalEvent represents local events and activities
type LocalEvent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Location    Location  `json:"location"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Available   bool      `json:"available"`
}

// TransportOption represents transportation choices
type TransportOption struct {
	Type       string  `json:"type"` // flight, train, bus, car_rental
	From       string  `json:"from"`
	To         string  `json:"to"`
	Duration   string  `json:"duration"`
	Price      float64 `json:"price"`
	Available  bool    `json:"available"`
	BookingURL string  `json:"booking_url,omitempty"`
	Provider   string  `json:"provider"`
}

// EMTItem represents Emergency Medical Tourism inventory
type EMTItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // service, equipment, facility
	Location    Location `json:"location"`
	Available   bool     `json:"available"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
}

// RetrievalRequest represents a request for contextual data
type RetrievalRequest struct {
	UserID      string                 `json:"user_id"`
	Destination string                 `json:"destination"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     time.Time              `json:"end_date"`
	Budget      float64                `json:"budget"`
	Travelers   int                    `json:"travelers"`
	Interests   []string               `json:"interests"`
	Preferences map[string]interface{} `json:"preferences"`
}

// RetrieveContext fetches comprehensive context for trip planning
func (r *RAGRetriever) RetrieveContext(ctx context.Context, req RetrievalRequest) (*TripContext, error) {
	log.Printf("Retrieving context for destination: %s", req.Destination)

	tripContext := &TripContext{
		Destination: req.Destination,
	}

	// Retrieve user profile
	if req.UserID != "" && r.firebase != nil {
		profile, err := r.firebase.GetUserProfile(ctx, req.UserID)
		if err == nil {
			tripContext.UserProfile = profile
		}
	}

	// Fetch attractions using data connector
	attractions, err := r.dataConnector.FetchAttractions(ctx, req.Destination, req.Interests)
	if err != nil {
		log.Printf("Error fetching attractions: %v", err)
		attractions = r.getMockAttractions(req.Destination)
	}

	// Fetch hotels using data connector
	hotels, err := r.dataConnector.FetchHotels(ctx, req.Destination, req.StartDate, req.EndDate, req.Budget)
	if err != nil {
		log.Printf("Error fetching hotels: %v", err)
		hotels = r.getMockHotels(req.Destination)
	}

	// Fetch weather forecast
	weather, err := r.fetchWeather(ctx, req.Destination, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		weather = WeatherForecast{} // Empty weather
	}

	// Fetch local events
	events, err := r.fetchLocalEvents(ctx, req.Destination, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Error fetching local events: %v", err)
	} else {
		tripContext.LocalEvents = events
	}

	// Apply validation and ranking
	if r.validator != nil {
		criteria := ValidationCriteria{
			Budget:            req.Budget,
			RequiredRating:    3.0, // Minimum rating
			PreferredTypes:    req.Interests,
			Preferences:       req.Preferences,
			AvailabilityCheck: true,
		}
		
		weights := DefaultRankingWeights()

		// Validate and rank attractions
		validatedAttractions, err := r.validator.ValidateAndRankAttractions(ctx, attractions, criteria, weights)
		if err == nil {
			tripContext.Attractions = validatedAttractions
		} else {
			tripContext.Attractions = attractions
		}

		// Validate and rank hotels
		validatedHotels, err := r.validator.ValidateAndRankHotels(ctx, hotels, criteria, weights)
		if err == nil {
			tripContext.Hotels = validatedHotels
		} else {
			tripContext.Hotels = hotels
		}

		// Apply budget constraints to the entire context
		if err := r.validator.ApplyBudgetConstraints(ctx, tripContext, req.Budget); err != nil {
			log.Printf("Error applying budget constraints: %v", err)
		}
	} else {
		tripContext.Attractions = attractions
		tripContext.Hotels = hotels
	}

	tripContext.Weather = weather

	// Fetch transportation options
	transport, err := r.fetchTransportation(ctx, req.Destination, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Error fetching transportation: %v", err)
	} else {
		tripContext.Transportation = transport
	}

	// Retrieve similar trips from Firebase
	if r.firebase != nil {
		similarTrips, err := r.fetchSimilarTrips(ctx, req)
		if err == nil {
			tripContext.SimilarTrips = similarTrips
		}
	}

	// Fetch EMT inventory
	emtItems, err := r.fetchEMTInventory(ctx, req.Destination)
	if err != nil {
		log.Printf("Error fetching EMT inventory: %v", err)
	} else {
		tripContext.EMTInventory = emtItems
	}

	return tripContext, nil
}

// fetchAttractions retrieves attractions using Google Places API
func (r *RAGRetriever) fetchAttractions(ctx context.Context, destination string, interests []string) ([]Attraction, error) {
	if r.mapsAPIKey == "" {
		return r.getMockAttractions(destination), nil
	}

	// Google Places API implementation would go here
	// For now, return mock data
	return r.getMockAttractions(destination), nil
}

// fetchHotels retrieves hotel options
func (r *RAGRetriever) fetchHotels(ctx context.Context, destination string, budget float64) ([]Hotel, error) {
	// Mock hotel data - in production, integrate with booking APIs
	return []Hotel{
		{
			ID:            "hotel_1",
			Name:          fmt.Sprintf("Grand Hotel %s", destination),
			Location:      Location{Address: fmt.Sprintf("Central %s", destination)},
			Rating:        4.5,
			PricePerNight: budget * 0.3, // 30% of daily budget
			Amenities:     []string{"WiFi", "Breakfast", "Pool"},
			Available:     true,
		},
		{
			ID:            "hotel_2",
			Name:          fmt.Sprintf("Budget Inn %s", destination),
			Location:      Location{Address: fmt.Sprintf("Downtown %s", destination)},
			Rating:        3.8,
			PricePerNight: budget * 0.15, // 15% of daily budget
			Amenities:     []string{"WiFi", "Parking"},
			Available:     true,
		},
	}, nil
}

// fetchWeather retrieves weather forecast
func (r *RAGRetriever) fetchWeather(ctx context.Context, destination string, startDate, endDate time.Time) (WeatherForecast, error) {
	// Mock weather data - in production, integrate with weather API
	forecast := WeatherForecast{
		Current: WeatherCondition{
			Date:        time.Now(),
			Temperature: 22.0,
			Description: "Partly cloudy",
			Humidity:    65,
			WindSpeed:   10.5,
			Icon:        "partly-cloudy",
		},
	}

	// Generate forecast for trip duration
	for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		forecast.Forecast = append(forecast.Forecast, WeatherCondition{
			Date:        d,
			Temperature: 20.0 + float64(len(forecast.Forecast)%5), // Vary temperature
			Description: "Clear skies",
			Humidity:    60,
			WindSpeed:   8.0,
			Icon:        "sunny",
		})
	}

	return forecast, nil
}

// fetchLocalEvents retrieves local events and activities
func (r *RAGRetriever) fetchLocalEvents(ctx context.Context, destination string, startDate, endDate time.Time) ([]LocalEvent, error) {
	// Mock events data
	return []LocalEvent{
		{
			ID:          "event_1",
			Name:        fmt.Sprintf("%s Music Festival", destination),
			Location:    Location{Address: fmt.Sprintf("Central Park, %s", destination)},
			Date:        startDate.AddDate(0, 0, 1),
			Category:    "Music",
			Price:       50.0,
			Description: "Annual music festival featuring local and international artists",
			Available:   true,
		},
		{
			ID:          "event_2",
			Name:        fmt.Sprintf("%s Food Market", destination),
			Location:    Location{Address: fmt.Sprintf("Market Square, %s", destination)},
			Date:        startDate.AddDate(0, 0, 2),
			Category:    "Food",
			Price:       0.0,
			Description: "Weekly food market with local delicacies",
			Available:   true,
		},
	}, nil
}

// fetchTransportation retrieves transportation options
func (r *RAGRetriever) fetchTransportation(ctx context.Context, destination string, startDate, endDate time.Time) ([]TransportOption, error) {
	// Mock transportation data
	return []TransportOption{
		{
			Type:      "flight",
			From:      "Current Location",
			To:        destination,
			Duration:  "2h 30m",
			Price:     299.99,
			Available: true,
			Provider:  "AirlineX",
		},
		{
			Type:      "train",
			From:      "Current Location",
			To:        destination,
			Duration:  "4h 15m",
			Price:     89.99,
			Available: true,
			Provider:  "RailService",
		},
	}, nil
}

// fetchSimilarTrips retrieves similar trips from Firebase
func (r *RAGRetriever) fetchSimilarTrips(ctx context.Context, req RetrievalRequest) ([]TripData, error) {
	if req.UserID == "" {
		return []TripData{}, nil
	}

	// Get user's trip history
	trips, err := r.firebase.GetUserTrips(ctx, req.UserID)
	if err != nil {
		return []TripData{}, err
	}

	// Filter and rank similar trips
	var similarTrips []TripData
	for _, trip := range trips {
		// Simple similarity based on destination proximity
		if trip.Destination != "" {
			similarTrips = append(similarTrips, trip)
		}
	}

	// Limit to top 5 similar trips
	if len(similarTrips) > 5 {
		similarTrips = similarTrips[:5]
	}

	return similarTrips, nil
}

// fetchEMTInventory retrieves Emergency Medical Tourism inventory
func (r *RAGRetriever) fetchEMTInventory(ctx context.Context, destination string) ([]EMTItem, error) {
	// Mock EMT inventory data
	return []EMTItem{
		{
			ID:          "emt_1",
			Name:        fmt.Sprintf("%s International Hospital", destination),
			Type:        "facility",
			Location:    Location{Address: fmt.Sprintf("Medical District, %s", destination)},
			Available:   true,
			Description: "24/7 emergency medical services for international travelers",
			Contact:     "+1-800-MEDICAL",
		},
		{
			ID:          "emt_2",
			Name:        "Travel Medical Kit",
			Type:        "equipment",
			Location:    Location{Address: fmt.Sprintf("Pharmacy Chain, %s", destination)},
			Available:   true,
			Description: "Essential medical supplies for travelers",
			Contact:     "+1-800-PHARMACY",
		},
	}, nil
}

// getMockAttractions returns mock attraction data
func (r *RAGRetriever) getMockAttractions(destination string) []Attraction {
	return []Attraction{
		{
			ID:           "attr_1",
			Name:         fmt.Sprintf("%s Museum of History", destination),
			Type:         "museum",
			Location:     Location{Address: fmt.Sprintf("Museum Quarter, %s", destination)},
			Rating:       4.3,
			PriceLevel:   2,
			OpeningHours: []string{"9:00-17:00", "Mon-Sun"},
			Description:  "Explore the rich history and culture",
			Tags:         []string{"history", "culture", "educational"},
			Available:    true,
		},
		{
			ID:           "attr_2",
			Name:         fmt.Sprintf("%s Central Park", destination),
			Type:         "park",
			Location:     Location{Address: fmt.Sprintf("City Center, %s", destination)},
			Rating:       4.6,
			PriceLevel:   0,
			OpeningHours: []string{"6:00-22:00", "Daily"},
			Description:  "Beautiful park perfect for relaxation",
			Tags:         []string{"nature", "outdoor", "relaxation"},
			Available:    true,
		},
		{
			ID:           "attr_3",
			Name:         fmt.Sprintf("Local Cuisine Restaurant in %s", destination),
			Type:         "restaurant",
			Location:     Location{Address: fmt.Sprintf("Culinary District, %s", destination)},
			Rating:       4.4,
			PriceLevel:   3,
			OpeningHours: []string{"11:00-23:00", "Daily"},
			Description:  "Authentic local dishes and flavors",
			Tags:         []string{"food", "local", "dining"},
			Available:    true,
		},
	}
}

func (r *RAGRetriever) getMockHotels(destination string) []Hotel {
	return []Hotel{
		{
			ID:            "hotel_1",
			Name:          fmt.Sprintf("Grand %s Hotel", destination),
			Location:      Location{Address: fmt.Sprintf("Downtown, %s", destination)},
			Rating:        4.2,
			PricePerNight: 120,
			Amenities:     []string{"WiFi", "Pool", "Restaurant", "Gym"},
			Available:     true,
			BookingURL:    "https://example.com/book",
		},
		{
			ID:            "hotel_2",
			Name:          fmt.Sprintf("%s Budget Inn", destination),
			Location:      Location{Address: fmt.Sprintf("City Center, %s", destination)},
			Rating:        3.8,
			PricePerNight: 80,
			Amenities:     []string{"WiFi", "Parking"},
			Available:     true,
			BookingURL:    "https://example.com/book",
		},
	}
}

// ValidateAvailability checks and updates availability of retrieved items
func (r *RAGRetriever) ValidateAvailability(ctx context.Context, tripContext *TripContext) error {
	// Validate attraction availability
	for i := range tripContext.Attractions {
		// Mock validation - in production, check real-time availability
		tripContext.Attractions[i].Available = true
	}

	// Validate hotel availability
	for i := range tripContext.Hotels {
		// Mock validation - in production, check booking APIs
		tripContext.Hotels[i].Available = true
	}

	// Validate event availability
	for i := range tripContext.LocalEvents {
		// Mock validation - in production, check event APIs
		tripContext.LocalEvents[i].Available = true
	}

	return nil
}

// RankByRelevance sorts retrieved items by relevance to user preferences
func (r *RAGRetriever) RankByRelevance(tripContext *TripContext, interests []string) {
	// Rank attractions by relevance
	sort.Slice(tripContext.Attractions, func(i, j int) bool {
		scoreI := r.calculateRelevanceScore(tripContext.Attractions[i].Tags, interests)
		scoreJ := r.calculateRelevanceScore(tripContext.Attractions[j].Tags, interests)
		return scoreI > scoreJ
	})

	// Rank hotels by rating and price
	sort.Slice(tripContext.Hotels, func(i, j int) bool {
		return tripContext.Hotels[i].Rating > tripContext.Hotels[j].Rating
	})
}

// calculateRelevanceScore calculates how relevant an item is to user interests
func (r *RAGRetriever) calculateRelevanceScore(tags, interests []string) int {
	score := 0
	for _, tag := range tags {
		for _, interest := range interests {
			if tag == interest {
				score++
			}
		}
	}
	return score
}
