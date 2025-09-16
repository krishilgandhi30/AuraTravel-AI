package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	DisplayName       string     `json:"display_name"`
	PhotoURL          string     `json:"photo_url"`
	EmailVerified     bool       `json:"email_verified"`
	PhoneNumber       string     `json:"phone_number"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	DateOfBirth       *time.Time `json:"date_of_birth"`
	Nationality       string     `json:"nationality"`
	PreferredCurrency string     `json:"preferred_currency"`
	PreferredLanguage string     `json:"preferred_language"`
	IsActive          bool       `json:"is_active"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TravelPreferences stores user's travel preferences
type TravelPreferences struct {
	ID                   uint                 `json:"id"`
	UserID               string               `json:"user_id"`
	BudgetRange          string               `json:"budget_range"`
	TravelStyles         []string             `json:"travel_styles"`
	Interests            []string             `json:"interests"`
	AccommodationTypes   []string             `json:"accommodation_types"`
	TransportationModes  []string             `json:"transportation_modes"`
	FoodPreferences      []string             `json:"food_preferences"`
	AccessibilityNeeds   []string             `json:"accessibility_needs"`
	NotificationSettings NotificationSettings `json:"notification_settings"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

// NotificationSettings for user preferences
type NotificationSettings struct {
	Email           bool `json:"email"`
	SMS             bool `json:"sms"`
	Push            bool `json:"push"`
	TripReminders   bool `json:"trip_reminders"`
	PriceAlerts     bool `json:"price_alerts"`
	Recommendations bool `json:"recommendations"`
}

// EmergencyContact stores emergency contact information
type EmergencyContact struct {
	ID           uint      `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Trip represents a travel trip
type Trip struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Destination string    `json:"destination"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"` // draft, planned, ongoing, completed, cancelled
	Travelers   int       `json:"travelers"`
	TotalBudget float64   `json:"total_budget"`
	Currency    string    `json:"currency"`
	IsPublic    bool      `json:"is_public"`
	ShareCode   string    `json:"share_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TripPreferences stores preferences specific to a trip
type TripPreferences struct {
	ID                 uint      `json:"id"`
	TripID             string    `json:"trip_id"`
	TravelStyles       []string  `json:"travel_styles"`
	Interests          []string  `json:"interests"`
	AccommodationType  string    `json:"accommodation_type"`
	TransportationType string    `json:"transportation_type"`
	FoodPreferences    []string  `json:"food_preferences"`
	AccessibilityNeeds []string  `json:"accessibility_needs"`
	AdditionalRequests string    `json:"additional_requests"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TripCollaborator represents users who can collaborate on a trip
type TripCollaborator struct {
	ID         uint       `json:"id"`
	TripID     string     `json:"trip_id"`
	UserID     string     `json:"user_id"`
	Role       string     `json:"role"` // owner, editor, viewer
	InvitedBy  string     `json:"invited_by"`
	InvitedAt  time.Time  `json:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Itinerary represents a trip's detailed itinerary
type Itinerary struct {
	ID                string     `json:"id"`
	TripID            string     `json:"trip_id"`
	TotalActivities   int        `json:"total_activities"`
	EstimatedCost     float64    `json:"estimated_cost"`
	Currency          string     `json:"currency"`
	GeneratedBy       string     `json:"generated_by"` // ai, manual
	GeneratedAt       time.Time  `json:"generated_at"`
	LastOptimizedAt   *time.Time `json:"last_optimized_at"`
	OptimizationScore float64    `json:"optimization_score"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// DayPlan represents a single day's plan in an itinerary
type DayPlan struct {
	ID            string    `json:"id"`
	ItineraryID   string    `json:"itinerary_id"`
	Date          time.Time `json:"date"`
	DayNumber     int       `json:"day_number"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	TotalCost     float64   `json:"total_cost"`
	EstimatedTime int       `json:"estimated_time"` // in minutes
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Activity represents a travel activity
type Activity struct {
	ID              string     `json:"id"`
	TripID          string     `json:"trip_id"`
	DayPlanID       *string    `json:"day_plan_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Type            string     `json:"type"` // sightseeing, museum, adventure, cultural, etc.
	Location        Location   `json:"location"`
	Duration        int        `json:"duration"` // in minutes
	Cost            float64    `json:"cost"`
	Currency        string     `json:"currency"`
	Rating          float64    `json:"rating"`
	BookingRequired bool       `json:"booking_required"`
	BookingURL      string     `json:"booking_url"`
	OpeningHours    string     `json:"opening_hours"` // JSON string
	Tips            []string   `json:"tips"`
	Images          []string   `json:"images"`
	Tags            []string   `json:"tags"`
	ScheduledTime   *time.Time `json:"scheduled_time"`
	Priority        int        `json:"priority"`
	Status          string     `json:"status"` // planned, booked, completed, skipped
	Source          string     `json:"source"` // ai, manual, api
	ExternalID      string     `json:"external_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Meal represents a dining experience
type Meal struct {
	ID                  string     `json:"id"`
	DayPlanID           string     `json:"day_plan_id"`
	Name                string     `json:"name"`
	Type                string     `json:"type"` // breakfast, lunch, dinner, snack
	Location            Location   `json:"location"`
	Cuisine             string     `json:"cuisine"`
	Cost                float64    `json:"cost"`
	Currency            string     `json:"currency"`
	Rating              float64    `json:"rating"`
	ReservationRequired bool       `json:"reservation_required"`
	DietaryOptions      []string   `json:"dietary_options"`
	ScheduledTime       *time.Time `json:"scheduled_time"`
	BookingURL          string     `json:"booking_url"`
	Images              []string   `json:"images"`
	Source              string     `json:"source"` // ai, manual, api
	ExternalID          string     `json:"external_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Accommodation represents lodging
type Accommodation struct {
	ID                 string     `json:"id"`
	TripID             string     `json:"trip_id"`
	DayPlanID          *string    `json:"day_plan_id"`
	Name               string     `json:"name"`
	Type               string     `json:"type"` // hotel, resort, apartment, hostel, villa, boutique
	Location           Location   `json:"location"`
	PricePerNight      float64    `json:"price_per_night"`
	Currency           string     `json:"currency"`
	TotalNights        int        `json:"total_nights"`
	TotalCost          float64    `json:"total_cost"`
	Rating             float64    `json:"rating"`
	Amenities          []string   `json:"amenities"`
	Images             []string   `json:"images"`
	BookingURL         string     `json:"booking_url"`
	CheckInTime        string     `json:"check_in_time"`
	CheckOutTime       string     `json:"check_out_time"`
	CheckInDate        *time.Time `json:"check_in_date"`
	CheckOutDate       *time.Time `json:"check_out_date"`
	CancellationPolicy string     `json:"cancellation_policy"`
	Status             string     `json:"status"` // planned, booked, checked_in, checked_out
	Source             string     `json:"source"` // ai, manual, api
	ExternalID         string     `json:"external_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Transportation represents travel transportation
type Transportation struct {
	ID            string     `json:"id"`
	TripID        string     `json:"trip_id"`
	DayPlanID     *string    `json:"day_plan_id"`
	Type          string     `json:"type"` // flight, train, bus, car, taxi, metro, walk
	FromLocation  Location   `json:"from_location"`
	ToLocation    Location   `json:"to_location"`
	Provider      string     `json:"provider"`
	DepartureTime *time.Time `json:"departure_time"`
	ArrivalTime   *time.Time `json:"arrival_time"`
	Duration      int        `json:"duration"` // in minutes
	Cost          float64    `json:"cost"`
	Currency      string     `json:"currency"`
	BookingURL    string     `json:"booking_url"`
	BookingRef    string     `json:"booking_ref"`
	SeatNumber    string     `json:"seat_number"`
	Class         string     `json:"class"`
	Notes         string     `json:"notes"`
	Status        string     `json:"status"` // planned, booked, completed
	Source        string     `json:"source"` // ai, manual, api
	ExternalID    string     `json:"external_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Location represents a geographical location
type Location struct {
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Timezone   string  `json:"timezone"`
	PlaceID    string  `json:"place_id"`
}

// Recommendation represents AI-generated recommendations
type Recommendation struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	TripID       string     `json:"trip_id"`
	Type         string     `json:"type"` // destination, activity, restaurant, accommodation
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Location     Location   `json:"location"`
	Images       []string   `json:"images"`
	Rating       float64    `json:"rating"`
	Price        float64    `json:"price"`
	Currency     string     `json:"currency"`
	Tags         []string   `json:"tags"`
	Reasons      []string   `json:"reasons"`
	Confidence   float64    `json:"confidence"`
	Source       string     `json:"source"` // gemini, vertex_ai, manual
	ModelVersion string     `json:"model_version"`
	IsAccepted   bool       `json:"is_accepted"`
	AcceptedAt   *time.Time `json:"accepted_at"`
	ExternalID   string     `json:"external_id"`
	ExternalURL  string     `json:"external_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// SearchHistory stores user search queries for analytics
type SearchHistory struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Query     string    `json:"query"`
	Type      string    `json:"type"`    // destination, activity, accommodation
	Filters   string    `json:"filters"` // JSON string of applied filters
	Results   int       `json:"results"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsEvent stores events for analytics
type AnalyticsEvent struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	EventData string    `json:"event_data"` // JSON string
	SessionID string    `json:"session_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	CreatedAt time.Time `json:"created_at"`
}

func (TravelPreferences) TableName() string { return "travel_preferences" }
func (EmergencyContact) TableName() string  { return "emergency_contacts" }
