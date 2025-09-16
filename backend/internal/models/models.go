package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	DisplayName       string     `json:"display_name"`
	PhotoURL          string     `json:"photo_url"`
	EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
	PhoneNumber       string     `json:"phone_number"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	DateOfBirth       *time.Time `json:"date_of_birth"`
	Nationality       string     `json:"nationality"`
	PreferredCurrency string     `json:"preferred_currency" gorm:"default:'USD'"`
	PreferredLanguage string     `json:"preferred_language" gorm:"default:'en'"`
	IsActive          bool       `json:"is_active" gorm:"default:true"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TravelPreferences stores user's travel preferences
type TravelPreferences struct {
	ID                   uint                 `json:"id" gorm:"primaryKey"`
	UserID               string               `json:"user_id" gorm:"not null"`
	BudgetRange          string               `json:"budget_range"`
	TravelStyles         []string             `json:"travel_styles" gorm:"type:text[]"`
	Interests            []string             `json:"interests" gorm:"type:text[]"`
	AccommodationTypes   []string             `json:"accommodation_types" gorm:"type:text[]"`
	TransportationModes  []string             `json:"transportation_modes" gorm:"type:text[]"`
	FoodPreferences      []string             `json:"food_preferences" gorm:"type:text[]"`
	AccessibilityNeeds   []string             `json:"accessibility_needs" gorm:"type:text[]"`
	NotificationSettings NotificationSettings `json:"notification_settings" gorm:"embedded"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

// NotificationSettings for user preferences
type NotificationSettings struct {
	Email           bool `json:"email" gorm:"default:true"`
	SMS             bool `json:"sms" gorm:"default:false"`
	Push            bool `json:"push" gorm:"default:true"`
	TripReminders   bool `json:"trip_reminders" gorm:"default:true"`
	PriceAlerts     bool `json:"price_alerts" gorm:"default:true"`
	Recommendations bool `json:"recommendations" gorm:"default:true"`
}

// EmergencyContact stores emergency contact information
type EmergencyContact struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null"`
	Relationship string    `json:"relationship"`
	Phone        string    `json:"phone" gorm:"not null"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Trip represents a travel trip
type Trip struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Destination string    `json:"destination" gorm:"not null"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status" gorm:"default:'draft'"` // draft, planned, ongoing, completed, cancelled
	Travelers   int       `json:"travelers" gorm:"default:1"`
	TotalBudget float64   `json:"total_budget"`
	Currency    string    `json:"currency" gorm:"default:'USD'"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	ShareCode   string    `json:"share_code" gorm:"uniqueIndex"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TripPreferences stores preferences specific to a trip
type TripPreferences struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	TripID             string    `json:"trip_id" gorm:"not null"`
	TravelStyles       []string  `json:"travel_styles" gorm:"type:text[]"`
	Interests          []string  `json:"interests" gorm:"type:text[]"`
	AccommodationType  string    `json:"accommodation_type"`
	TransportationType string    `json:"transportation_type"`
	FoodPreferences    []string  `json:"food_preferences" gorm:"type:text[]"`
	AccessibilityNeeds []string  `json:"accessibility_needs" gorm:"type:text[]"`
	AdditionalRequests string    `json:"additional_requests"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TripCollaborator represents users who can collaborate on a trip
type TripCollaborator struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	TripID     string     `json:"trip_id" gorm:"not null"`
	UserID     string     `json:"user_id" gorm:"not null"`
	Role       string     `json:"role" gorm:"default:'viewer'"` // owner, editor, viewer
	InvitedBy  string     `json:"invited_by"`
	InvitedAt  time.Time  `json:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Itinerary represents a trip's detailed itinerary
type Itinerary struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	TripID            string     `json:"trip_id" gorm:"not null;uniqueIndex"`
	TotalActivities   int        `json:"total_activities"`
	EstimatedCost     float64    `json:"estimated_cost"`
	Currency          string     `json:"currency" gorm:"default:'USD'"`
	GeneratedBy       string     `json:"generated_by"` // ai, manual
	GeneratedAt       time.Time  `json:"generated_at"`
	LastOptimizedAt   *time.Time `json:"last_optimized_at"`
	OptimizationScore float64    `json:"optimization_score"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// DayPlan represents a single day's plan in an itinerary
type DayPlan struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	ItineraryID   string    `json:"itinerary_id" gorm:"not null"`
	Date          time.Time `json:"date" gorm:"not null"`
	DayNumber     int       `json:"day_number" gorm:"not null"`
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
	ID              string     `json:"id" gorm:"primaryKey"`
	TripID          string     `json:"trip_id"`
	DayPlanID       *string    `json:"day_plan_id"`
	Name            string     `json:"name" gorm:"not null"`
	Description     string     `json:"description"`
	Type            string     `json:"type"` // sightseeing, museum, adventure, cultural, etc.
	Location        Location   `json:"location" gorm:"embedded;embeddedPrefix:location_"`
	Duration        int        `json:"duration"` // in minutes
	Cost            float64    `json:"cost"`
	Currency        string     `json:"currency" gorm:"default:'USD'"`
	Rating          float64    `json:"rating"`
	BookingRequired bool       `json:"booking_required"`
	BookingURL      string     `json:"booking_url"`
	OpeningHours    string     `json:"opening_hours"` // JSON string
	Tips            []string   `json:"tips" gorm:"type:text[]"`
	Images          []string   `json:"images" gorm:"type:text[]"`
	Tags            []string   `json:"tags" gorm:"type:text[]"`
	ScheduledTime   *time.Time `json:"scheduled_time"`
	Priority        int        `json:"priority" gorm:"default:0"`
	Status          string     `json:"status" gorm:"default:'planned'"` // planned, booked, completed, skipped
	Source          string     `json:"source"`                          // ai, manual, api
	ExternalID      string     `json:"external_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Meal represents a dining experience
type Meal struct {
	ID                  string     `json:"id" gorm:"primaryKey"`
	DayPlanID           string     `json:"day_plan_id" gorm:"not null"`
	Name                string     `json:"name" gorm:"not null"`
	Type                string     `json:"type"` // breakfast, lunch, dinner, snack
	Location            Location   `json:"location" gorm:"embedded;embeddedPrefix:location_"`
	Cuisine             string     `json:"cuisine"`
	Cost                float64    `json:"cost"`
	Currency            string     `json:"currency" gorm:"default:'USD'"`
	Rating              float64    `json:"rating"`
	ReservationRequired bool       `json:"reservation_required"`
	DietaryOptions      []string   `json:"dietary_options" gorm:"type:text[]"`
	ScheduledTime       *time.Time `json:"scheduled_time"`
	BookingURL          string     `json:"booking_url"`
	Images              []string   `json:"images" gorm:"type:text[]"`
	Source              string     `json:"source"` // ai, manual, api
	ExternalID          string     `json:"external_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Accommodation represents lodging
type Accommodation struct {
	ID                 string     `json:"id" gorm:"primaryKey"`
	TripID             string     `json:"trip_id"`
	DayPlanID          *string    `json:"day_plan_id"`
	Name               string     `json:"name" gorm:"not null"`
	Type               string     `json:"type"` // hotel, resort, apartment, hostel, villa, boutique
	Location           Location   `json:"location" gorm:"embedded;embeddedPrefix:location_"`
	PricePerNight      float64    `json:"price_per_night"`
	Currency           string     `json:"currency" gorm:"default:'USD'"`
	TotalNights        int        `json:"total_nights"`
	TotalCost          float64    `json:"total_cost"`
	Rating             float64    `json:"rating"`
	Amenities          []string   `json:"amenities" gorm:"type:text[]"`
	Images             []string   `json:"images" gorm:"type:text[]"`
	BookingURL         string     `json:"booking_url"`
	CheckInTime        string     `json:"check_in_time"`
	CheckOutTime       string     `json:"check_out_time"`
	CheckInDate        *time.Time `json:"check_in_date"`
	CheckOutDate       *time.Time `json:"check_out_date"`
	CancellationPolicy string     `json:"cancellation_policy"`
	Status             string     `json:"status" gorm:"default:'planned'"` // planned, booked, checked_in, checked_out
	Source             string     `json:"source"`                          // ai, manual, api
	ExternalID         string     `json:"external_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Transportation represents travel transportation
type Transportation struct {
	ID            string     `json:"id" gorm:"primaryKey"`
	TripID        string     `json:"trip_id"`
	DayPlanID     *string    `json:"day_plan_id"`
	Type          string     `json:"type"` // flight, train, bus, car, taxi, metro, walk
	FromLocation  Location   `json:"from_location" gorm:"embedded;embeddedPrefix:from_"`
	ToLocation    Location   `json:"to_location" gorm:"embedded;embeddedPrefix:to_"`
	Provider      string     `json:"provider"`
	DepartureTime *time.Time `json:"departure_time"`
	ArrivalTime   *time.Time `json:"arrival_time"`
	Duration      int        `json:"duration"` // in minutes
	Cost          float64    `json:"cost"`
	Currency      string     `json:"currency" gorm:"default:'USD'"`
	BookingURL    string     `json:"booking_url"`
	BookingRef    string     `json:"booking_ref"`
	SeatNumber    string     `json:"seat_number"`
	Class         string     `json:"class"`
	Notes         string     `json:"notes"`
	Status        string     `json:"status" gorm:"default:'planned'"` // planned, booked, completed
	Source        string     `json:"source"`                          // ai, manual, api
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
	ID           string     `json:"id" gorm:"primaryKey"`
	UserID       string     `json:"user_id"`
	TripID       string     `json:"trip_id"`
	Type         string     `json:"type"` // destination, activity, restaurant, accommodation
	Title        string     `json:"title" gorm:"not null"`
	Description  string     `json:"description"`
	Location     Location   `json:"location" gorm:"embedded;embeddedPrefix:location_"`
	Images       []string   `json:"images" gorm:"type:text[]"`
	Rating       float64    `json:"rating"`
	Price        float64    `json:"price"`
	Currency     string     `json:"currency" gorm:"default:'USD'"`
	Tags         []string   `json:"tags" gorm:"type:text[]"`
	Reasons      []string   `json:"reasons" gorm:"type:text[]"`
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	Query     string    `json:"query"`
	Type      string    `json:"type"` // destination, activity, accommodation
	Filters   string    `json:"filters"` // JSON string of applied filters
	Results   int       `json:"results"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsEvent stores events for analytics
type AnalyticsEvent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	EventData string    `json:"event_data"` // JSON string
	SessionID string    `json:"session_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Referrer  string    `json:"referrer"`
	CreatedAt time.Time `json:"created_at"`
}

// GetTableName returns the table name for GORM
func (User) TableName() string             { return "users" }
func (Trip) TableName() string             { return "trips" }
func (Itinerary) TableName() string        { return "itineraries" }
func (Activity) TableName() string         { return "activities" }
func (Recommendation) TableName() string   { return "recommendations" }
func (TravelPreferences) TableName() string { return "travel_preferences" }
func (EmergencyContact) TableName() string { return "emergency_contacts" }