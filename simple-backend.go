package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Simple in-memory data structures
type Destination struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	Country     string   `json:"country"`
	Vibes       []string `json:"vibes"`
	CostLevel   string   `json:"cost_level"`
	Description string   `json:"description"`
	MinDays     int      `json:"min_days"`
	MaxDays     int      `json:"max_days"`
}

type TripRequest struct {
	UserID          string   `json:"user_id"`
	Destination     string   `json:"destination"`
	DestinationType string   `json:"destination_type"`
	Origin          string   `json:"origin"`
	Duration        int      `json:"duration"`
	Budget          float64  `json:"budget"`
	Currency        string   `json:"currency"`
	Vibes           []string `json:"vibes"`
	Preferences     string   `json:"preferences"`
}

type Trip struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Destination Destination `json:"destination"`
	Duration    int         `json:"duration"`
	Budget      float64     `json:"budget"`
	Status      string      `json:"status"`
}

type ChatRequest struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Sample data
var destinations = []Destination{
	{
		ID:          "1",
		Name:        "Rishikesh",
		City:        "Rishikesh",
		State:       "Uttarakhand",
		Country:     "India",
		Vibes:       []string{"serene", "spiritual", "adventure", "mountains"},
		CostLevel:   "medium",
		Description: "Spiritual capital of the world with adventure activities",
		MinDays:     3,
		MaxDays:     7,
	},
	{
		ID:          "2",
		Name:        "Goa Beaches",
		City:        "Panaji",
		State:       "Goa",
		Country:     "India",
		Vibes:       []string{"beach", "party", "relaxing", "nightlife"},
		CostLevel:   "medium",
		Description: "Beach paradise with vibrant nightlife",
		MinDays:     4,
		MaxDays:     10,
	},
	{
		ID:          "3",
		Name:        "Shimla",
		City:        "Shimla",
		State:       "Himachal Pradesh",
		Country:     "India",
		Vibes:       []string{"mountains", "serene", "colonial", "cool climate"},
		CostLevel:   "medium",
		Description: "Hill station with colonial charm",
		MinDays:     3,
		MaxDays:     6,
	},
	{
		ID:          "4",
		Name:        "Manali",
		City:        "Manali",
		State:       "Himachal Pradesh",
		Country:     "India",
		Vibes:       []string{"mountains", "adventure", "snow", "trekking"},
		CostLevel:   "medium",
		Description: "Adventure capital in the Himalayas",
		MinDays:     4,
		MaxDays:     8,
	},
	{
		ID:          "5",
		Name:        "Gujarat Heritage",
		City:        "Ahmedabad",
		State:       "Gujarat",
		Country:     "India",
		Vibes:       []string{"cultural", "heritage", "business", "vibrant"},
		CostLevel:   "medium",
		Description: "Rich cultural heritage with modern business centers, home to Gandhi's legacy",
		MinDays:     5,
		MaxDays:     10,
	},
	{
		ID:          "6",
		Name:        "Rann of Kutch",
		City:        "Bhuj",
		State:       "Gujarat",
		Country:     "India",
		Vibes:       []string{"unique", "desert", "cultural", "festivals"},
		CostLevel:   "medium",
		Description: "White desert landscape with vibrant festivals and handicrafts",
		MinDays:     3,
		MaxDays:     6,
	},
	{
		ID:          "7",
		Name:        "Dwarka",
		City:        "Dwarka",
		State:       "Gujarat",
		Country:     "India",
		Vibes:       []string{"spiritual", "religious", "peaceful", "historic"},
		CostLevel:   "low",
		Description: "Sacred city of Lord Krishna with ancient temples and peaceful atmosphere",
		MinDays:     2,
		MaxDays:     4,
	},
	{
		ID:          "8",
		Name:        "Somnath",
		City:        "Somnath",
		State:       "Gujarat",
		Country:     "India",
		Vibes:       []string{"spiritual", "beach", "religious", "historic"},
		CostLevel:   "low",
		Description: "Famous Somnath Temple with beautiful coastal location",
		MinDays:     2,
		MaxDays:     3,
	},
	{
		ID:          "9",
		Name:        "Rajasthan Golden Triangle",
		City:        "Jaipur",
		State:       "Rajasthan",
		Country:     "India",
		Vibes:       []string{"royal", "cultural", "heritage", "palaces"},
		CostLevel:   "medium",
		Description: "Royal palaces, forts, and vibrant culture of the Pink City",
		MinDays:     6,
		MaxDays:     12,
	},
	{
		ID:          "10",
		Name:        "Kerala Backwaters",
		City:        "Alleppey",
		State:       "Kerala",
		Country:     "India",
		Vibes:       []string{"serene", "nature", "relaxing", "unique"},
		CostLevel:   "medium",
		Description: "Peaceful backwater cruises through lush green landscapes",
		MinDays:     4,
		MaxDays:     8,
	},
}

var trips []Trip
var tripCounter = 1

func main() {
	// Add CORS middleware
	http.HandleFunc("/", corsMiddleware(routeHandler))

	port := getEnv("PORT", "8080")
	fmt.Printf("🚀 AuraTravel AI Backend starting on port %s\n", port)
	fmt.Printf("🌐 API available at: http://localhost:%s/api/v1/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.URL.Path == "/health":
		handleHealth(w, r)
	case r.URL.Path == "/api/v1/trips/plan" && r.Method == "POST":
		handlePlanTrip(w, r)
	case strings.HasPrefix(r.URL.Path, "/api/v1/trips/") && r.Method == "GET":
		handleGetTrip(w, r)
	case r.URL.Path == "/api/v1/trips/user/demo" && r.Method == "GET":
		handleGetUserTrips(w, r)
	case r.URL.Path == "/api/v1/destinations" && r.Method == "GET":
		handleGetDestinations(w, r)
	case r.URL.Path == "/api/v1/destinations/search" && r.Method == "GET":
		handleSearchDestinations(w, r)
	case r.URL.Path == "/api/v1/ai/chat" && r.Method == "POST":
		handleChatWithAI(w, r)
	default:
		writeResponse(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Endpoint not found",
		})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"status": "healthy"},
	})
}

func handlePlanTrip(w http.ResponseWriter, r *http.Request) {
	var req TripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Smart destination matching
	var selectedDest Destination
	found := false

	// First try exact name, city, or state match
	for _, dest := range destinations {
		if strings.EqualFold(dest.Name, req.Destination) ||
			strings.EqualFold(dest.City, req.Destination) ||
			strings.EqualFold(dest.State, req.Destination) {
			selectedDest = dest
			found = true
			break
		}
	}

	// Second, try partial matches in name, city, state, or description
	if !found {
		reqDestLower := strings.ToLower(req.Destination)
		for _, dest := range destinations {
			if strings.Contains(strings.ToLower(dest.Name), reqDestLower) ||
				strings.Contains(strings.ToLower(dest.City), reqDestLower) ||
				strings.Contains(strings.ToLower(dest.State), reqDestLower) ||
				strings.Contains(strings.ToLower(dest.Description), reqDestLower) {
				selectedDest = dest
				found = true
				break
			}
		}
	}

	// Third, try vibe matching
	if !found && len(req.Vibes) > 0 {
		for _, dest := range destinations {
			for _, destVibe := range dest.Vibes {
				for _, reqVibe := range req.Vibes {
					if strings.EqualFold(destVibe, reqVibe) {
						selectedDest = dest
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
	}

	// If still not found, return an error instead of defaulting
	if !found {
		writeResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Sorry, I couldn't find a destination matching '%s'. Try: Gujarat, Rishikesh, Goa, Shimla, Manali, Rajasthan, Kerala, Dwarka, or Somnath", req.Destination),
		})
		return
	}

	// Create trip
	tripID := generateID()
	trip := Trip{
		ID:          tripID,
		Title:       fmt.Sprintf("%d Days in %s", req.Duration, selectedDest.Name),
		Destination: selectedDest,
		Duration:    req.Duration,
		Budget:      req.Budget,
		Status:      "planned",
	}

	trips = append(trips, trip)

	writeResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    trip,
	})
}

func handleGetTrip(w http.ResponseWriter, r *http.Request) {
	tripID := strings.TrimPrefix(r.URL.Path, "/api/v1/trips/")
	for _, trip := range trips {
		if trip.ID == tripID {
			writeResponse(w, http.StatusOK, APIResponse{
				Success: true,
				Data:    trip,
			})
			return
		}
	}
	writeResponse(w, http.StatusNotFound, APIResponse{
		Success: false,
		Error:   "Trip not found",
	})
}

func handleGetUserTrips(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    trips,
	})
}

func handleGetDestinations(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    destinations,
	})
}

func handleSearchDestinations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var results []Destination

	if query == "" {
		writeResponse(w, http.StatusOK, APIResponse{
			Success: true,
			Data:    destinations,
		})
		return
	}

	queryLower := strings.ToLower(query)
	for _, dest := range destinations {
		if strings.Contains(strings.ToLower(dest.Name), queryLower) ||
			strings.Contains(strings.ToLower(dest.City), queryLower) ||
			strings.Contains(strings.ToLower(dest.State), queryLower) ||
			strings.Contains(strings.ToLower(dest.Description), queryLower) {
			results = append(results, dest)
			continue
		}
		// Also check vibes
		for _, vibe := range dest.Vibes {
			if strings.Contains(strings.ToLower(vibe), queryLower) {
				results = append(results, dest)
				break
			}
		}
	}

	writeResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

func handleChatWithAI(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	response := generateAIResponse(req.Message)

	writeResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"response": response,
			"intent":   "general",
			"entities": map[string]interface{}{},
		},
	})
}

// Helper functions
func generateID() string {
	id := fmt.Sprintf("trip_%d", tripCounter)
	tripCounter++
	return id
}

func generateAIResponse(message string) string {
	messageLower := strings.ToLower(message)

	// Check for specific state/place mentions
	if strings.Contains(messageLower, "gujarat") {
		return "Gujarat is fantastic! I recommend exploring Ahmedabad for rich cultural heritage and Gandhi's legacy, the stunning white desert of Rann of Kutch for unique landscapes, spiritual Dwarka for Lord Krishna's temples, and Somnath for beautiful coastal temples. Each offers a different facet of Gujarat's diverse culture!"
	}
	if strings.Contains(messageLower, "rajasthan") {
		return "Rajasthan is the land of kings! The Golden Triangle including Jaipur offers magnificent palaces, colorful markets, and royal heritage. You'll experience majestic forts, vibrant culture, and authentic Rajasthani cuisine. Perfect for 6-12 days of royal exploration!"
	}
	if strings.Contains(messageLower, "kerala") {
		return "Kerala is God's Own Country! The backwaters of Alleppey offer peaceful houseboat cruises through lush landscapes. It's perfect for relaxation, nature lovers, and experiencing unique waterway culture. Ideal for 4-8 days of serene beauty!"
	}

	// Check for activity preferences
	if strings.Contains(messageLower, "mountain") || strings.Contains(messageLower, "hill") {
		return "For mountain experiences, I suggest: Rishikesh for spiritual mountains with adventure, Shimla for colonial hill station charm, or Manali for snow-capped adventure activities. Each offers different mountain vibes - spiritual, heritage, or adventure!"
	}
	if strings.Contains(messageLower, "beach") || strings.Contains(messageLower, "sea") {
		return "For beach destinations: Goa offers vibrant beaches with nightlife and water sports, while Somnath in Gujarat combines spiritual temples with beautiful coastal views. Goa is great for parties, Somnath for peaceful spiritual beach experience!"
	}
	if strings.Contains(messageLower, "adventure") || strings.Contains(messageLower, "thrill") {
		return "For adventure seekers: Rishikesh offers river rafting, bungee jumping, and trekking. Manali provides paragliding, snow activities, and mountain adventures. Both destinations offer heart-pumping activities in stunning natural settings!"
	}
	if strings.Contains(messageLower, "spiritual") || strings.Contains(messageLower, "religious") || strings.Contains(messageLower, "temple") {
		return "For spiritual journeys: Rishikesh is the yoga capital with Ganga Aarti, Dwarka is Lord Krishna's sacred city, and Somnath has the famous Jyotirlinga temple. Each offers deep spiritual experiences and peace for the soul!"
	}
	if strings.Contains(messageLower, "cultural") || strings.Contains(messageLower, "heritage") {
		return "For cultural experiences: Gujarat offers Gandhi's heritage in Ahmedabad and traditional crafts in Rann of Kutch. Rajasthan's Jaipur showcases royal palaces and vibrant markets. Both states offer rich cultural immersion!"
	}
	if strings.Contains(messageLower, "budget") || strings.Contains(messageLower, "cheap") || strings.Contains(messageLower, "affordable") {
		return "For budget-friendly trips: Dwarka and Somnath in Gujarat offer low-cost spiritual experiences. Rishikesh and Shimla provide medium-budget options with great value. You can find affordable accommodation, local food, and many free activities!"
	}
	if strings.Contains(messageLower, "family") {
		return "For family trips: Shimla offers cool climate and easy accessibility for all ages. Gujarat destinations like Ahmedabad and Dwarka are family-friendly with good facilities. Kerala backwaters provide relaxing experiences perfect for family bonding!"
	}
	if strings.Contains(messageLower, "solo") || strings.Contains(messageLower, "alone") {
		return "For solo travelers: Rishikesh is safe and perfect for spiritual solo journeys with yoga and meditation. Goa offers vibrant solo travel scene with easy connections. Manali provides adventure opportunities for independent travelers!"
	}
	if strings.Contains(messageLower, "romantic") || strings.Contains(messageLower, "couple") {
		return "For romantic getaways: Shimla offers romantic hill station vibes, Kerala backwaters provide intimate houseboat experiences, and Manali has beautiful mountain settings. Goa beaches also offer romantic sunset views!"
	}

	// Weather and season queries
	if strings.Contains(messageLower, "winter") || strings.Contains(messageLower, "cold") {
		return "Winter is perfect for: Gujarat destinations (pleasant weather), Goa beaches (cool and comfortable), Rajasthan exploration (avoiding summer heat). Avoid hill stations if you don't like extreme cold!"
	}
	if strings.Contains(messageLower, "summer") || strings.Contains(messageLower, "hot") {
		return "For summer escapes: Head to hill stations like Shimla and Manali for cool climate. Kerala backwaters are also pleasant. Avoid Gujarat and Rajasthan during peak summer - they can be quite hot!"
	}

	// Duration-based suggestions
	if strings.Contains(messageLower, "weekend") || strings.Contains(messageLower, "2 day") || strings.Contains(messageLower, "short") {
		return "For short trips (2-3 days): Dwarka or Somnath for spiritual quick getaways, or Shimla for a brief hill station experience. These can be covered well in a weekend!"
	}
	if strings.Contains(messageLower, "week") || strings.Contains(messageLower, "7 day") {
		return "For a week-long trip: Explore multiple Gujarat destinations (Ahmedabad + Dwarka + Rann of Kutch), experience Rajasthan's Golden Triangle, or combine Rishikesh with nearby hill stations. Perfect duration for deeper exploration!"
	}

	// Help and planning queries
	if strings.Contains(messageLower, "help") || strings.Contains(messageLower, "suggest") || strings.Contains(messageLower, "recommend") {
		return "I'd love to help you plan the perfect trip! Tell me your preferences: Do you prefer mountains, beaches, or cultural cities? Are you looking for adventure, relaxation, spiritual experiences, or heritage exploration? What's your budget range and duration?"
	}
	if strings.Contains(messageLower, "plan") || strings.Contains(messageLower, "itinerary") {
		return "I can help you plan detailed itineraries! First, let me know: Which type of destination interests you - spiritual (Rishikesh, Dwarka), cultural (Gujarat, Rajasthan), adventure (Manali), beaches (Goa), or hill stations (Shimla)? Also share your duration and budget!"
	}

	// Food and cuisine
	if strings.Contains(messageLower, "food") || strings.Contains(messageLower, "cuisine") {
		return "For food lovers: Gujarat offers amazing vegetarian Gujarati thali and street food in Ahmedabad. Rajasthan has royal dal-baati-churma and desert cuisine. Goa provides seafood delicacies, while Kerala offers coconut-based dishes and spices!"
	}

	// Default response - more engaging
	return "Hello! I'm your AI travel assistant for incredible Indian destinations! 🌟 I can help you explore Gujarat's heritage sites, Rajasthan's royal palaces, Goa's beaches, hill stations like Shimla & Manali, spiritual places like Rishikesh, and Kerala's backwaters. What kind of travel experience are you looking for - adventure, culture, spirituality, beaches, or mountains?"
}

func writeResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
