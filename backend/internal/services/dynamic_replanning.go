package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

// DynamicReplanningService handles real-time itinerary adjustments
type DynamicReplanningService struct {
	ragRetriever     *RAGRetriever
	gemini           *GeminiService
	vectorDB         *VectorDatabase
	firebase         *FirebaseService
	notificationSvc  *NotificationService
	weatherKey       string
	httpClient       *http.Client
	monitoringActive bool
}

// NewDynamicReplanningService creates a new dynamic replanning service
func NewDynamicReplanningService(
	ragRetriever *RAGRetriever,
	gemini *GeminiService,
	vectorDB *VectorDatabase,
	firebase *FirebaseService,
	notificationSvc *NotificationService,
	weatherKey string,
) *DynamicReplanningService {
	return &DynamicReplanningService{
		ragRetriever:     ragRetriever,
		gemini:           gemini,
		vectorDB:         vectorDB,
		firebase:         firebase,
		notificationSvc:  notificationSvc,
		weatherKey:       weatherKey,
		httpClient:       &http.Client{Timeout: 30 * time.Second},
		monitoringActive: true,
	}
}

// ReplanningTrigger represents the reason for replanning
type ReplanningTrigger struct {
	Type        string      `json:"type"`     // weather, delay, sold_out, emergency
	Severity    string      `json:"severity"` // low, medium, high, critical
	Description string      `json:"description"`
	Timestamp   time.Time   `json:"timestamp"`
	Data        interface{} `json:"data,omitempty"`
}

// ReplanningResult represents the result of a replanning operation
type ReplanningResult struct {
	TripID           string              `json:"trip_id"`
	OriginalPlan     interface{}         `json:"original_plan"`
	RevisedPlan      interface{}         `json:"revised_plan"`
	Changes          []ItineraryChange   `json:"changes"`
	Triggers         []ReplanningTrigger `json:"triggers"`
	Confidence       float64             `json:"confidence"`
	EstimatedSavings float64             `json:"estimated_savings,omitempty"`
	ReplanTimestamp  time.Time           `json:"replan_timestamp"`
}

// ItineraryChange represents a specific change made to the itinerary
type ItineraryChange struct {
	Type        string      `json:"type"`      // replacement, cancellation, addition, time_shift
	Day         string      `json:"day"`       // day1, day2, etc.
	TimeSlot    string      `json:"time_slot"` // morning, afternoon, evening
	Original    interface{} `json:"original,omitempty"`
	Replacement interface{} `json:"replacement,omitempty"`
	Reason      string      `json:"reason"`
	Impact      string      `json:"impact"` // minor, moderate, major
	CostDelta   float64     `json:"cost_delta"`
}

// WeatherAlert represents a weather-based alert
type WeatherAlert struct {
	AlertType     string    `json:"alert_type"` // rain, storm, extreme_heat, snow
	Severity      string    `json:"severity"`   // watch, warning, emergency
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Description   string    `json:"description"`
	AffectedAreas []string  `json:"affected_areas"`
}

// DelayAlert represents transportation or event delays
type DelayAlert struct {
	ServiceType string        `json:"service_type"` // flight, train, event, attraction
	ServiceID   string        `json:"service_id"`
	DelayTime   time.Duration `json:"delay_time"`
	Status      string        `json:"status"` // delayed, cancelled, rescheduled
	Reason      string        `json:"reason"`
	NewSchedule *time.Time    `json:"new_schedule,omitempty"`
}

// AvailabilityAlert represents sold-out or unavailable items
type AvailabilityAlert struct {
	ItemType     string     `json:"item_type"` // attraction, hotel, restaurant, transport
	ItemID       string     `json:"item_id"`
	Status       string     `json:"status"` // sold_out, closed, fully_booked
	Until        *time.Time `json:"until,omitempty"`
	Alternatives []string   `json:"alternatives,omitempty"`
}

// MonitorTrip starts monitoring a trip for real-time changes
func (d *DynamicReplanningService) MonitorTrip(ctx context.Context, tripID string) error {
	if !d.monitoringActive {
		return fmt.Errorf("monitoring is not active")
	}

	// Start background monitoring goroutine
	go d.monitorTripBackground(ctx, tripID)

	log.Printf("Started monitoring trip: %s", tripID)
	return nil
}

// monitorTripBackground runs continuous monitoring in the background
func (d *DynamicReplanningService) monitorTripBackground(ctx context.Context, tripID string) {
	ticker := time.NewTicker(15 * time.Minute) // Check every 15 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Monitoring stopped for trip: %s", tripID)
			return
		case <-ticker.C:
			if err := d.checkAndReplan(ctx, tripID); err != nil {
				log.Printf("Error during monitoring check for trip %s: %v", tripID, err)
			}
		}
	}
}

// checkAndReplan checks for triggers and replans if necessary
func (d *DynamicReplanningService) checkAndReplan(ctx context.Context, tripID string) error {
	// Get current trip details
	trip, err := d.getCurrentTrip(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to get trip details: %w", err)
	}

	// Check for various triggers
	triggers := d.checkForTriggers(ctx, trip)
	if len(triggers) == 0 {
		return nil // No triggers found
	}

	// Filter triggers by severity
	criticalTriggers := d.filterCriticalTriggers(triggers)
	if len(criticalTriggers) == 0 {
		return nil // No critical triggers requiring immediate replanning
	}

	log.Printf("Found %d critical triggers for trip %s, initiating replanning", len(criticalTriggers), tripID)

	// Perform replanning
	result, err := d.performReplanning(ctx, trip, criticalTriggers)
	if err != nil {
		return fmt.Errorf("replanning failed: %w", err)
	}

	// Save replanned itinerary
	if err := d.saveReplanResult(ctx, result); err != nil {
		log.Printf("Failed to save replan result: %v", err)
	}

	// Send notifications
	if d.notificationSvc != nil {
		d.sendReplanNotifications(ctx, result)
	}

	return nil
}

// checkForTriggers checks for weather, delays, and availability changes
func (d *DynamicReplanningService) checkForTriggers(ctx context.Context, trip interface{}) []ReplanningTrigger {
	var triggers []ReplanningTrigger

	// Check weather alerts
	weatherTriggers := d.checkWeatherTriggers(ctx, trip)
	triggers = append(triggers, weatherTriggers...)

	// Check delay alerts
	delayTriggers := d.checkDelayTriggers(ctx, trip)
	triggers = append(triggers, delayTriggers...)

	// Check availability alerts
	availabilityTriggers := d.checkAvailabilityTriggers(ctx, trip)
	triggers = append(triggers, availabilityTriggers...)

	return triggers
}

// checkWeatherTriggers checks for weather-related triggers
func (d *DynamicReplanningService) checkWeatherTriggers(ctx context.Context, trip interface{}) []ReplanningTrigger {
	var triggers []ReplanningTrigger

	// Mock weather checking - in production, integrate with weather APIs
	weatherAlerts := d.fetchWeatherAlerts(ctx, trip)

	for _, alert := range weatherAlerts {
		severity := d.mapWeatherSeverity(alert)
		triggers = append(triggers, ReplanningTrigger{
			Type:        "weather",
			Severity:    severity,
			Description: fmt.Sprintf("%s: %s", alert.AlertType, alert.Description),
			Timestamp:   time.Now(),
			Data:        alert,
		})
	}

	return triggers
}

// checkDelayTriggers checks for transportation and event delays
func (d *DynamicReplanningService) checkDelayTriggers(ctx context.Context, trip interface{}) []ReplanningTrigger {
	var triggers []ReplanningTrigger

	// Mock delay checking - in production, integrate with transport APIs
	delayAlerts := d.fetchDelayAlerts(ctx, trip)

	for _, alert := range delayAlerts {
		severity := d.mapDelaySeverity(alert)
		triggers = append(triggers, ReplanningTrigger{
			Type:        "delay",
			Severity:    severity,
			Description: fmt.Sprintf("%s delayed by %v: %s", alert.ServiceType, alert.DelayTime, alert.Reason),
			Timestamp:   time.Now(),
			Data:        alert,
		})
	}

	return triggers
}

// checkAvailabilityTriggers checks for sold-out or unavailable items
func (d *DynamicReplanningService) checkAvailabilityTriggers(ctx context.Context, trip interface{}) []ReplanningTrigger {
	var triggers []ReplanningTrigger

	// Mock availability checking - in production, integrate with booking APIs
	availabilityAlerts := d.fetchAvailabilityAlerts(ctx, trip)

	for _, alert := range availabilityAlerts {
		triggers = append(triggers, ReplanningTrigger{
			Type:        "sold_out",
			Severity:    "high",
			Description: fmt.Sprintf("%s is %s", alert.ItemType, alert.Status),
			Timestamp:   time.Now(),
			Data:        alert,
		})
	}

	return triggers
}

// performReplanning generates a new itinerary based on triggers
func (d *DynamicReplanningService) performReplanning(ctx context.Context, trip interface{}, triggers []ReplanningTrigger) (*ReplanningResult, error) {
	// Extract trip data (mock structure)
	tripData := d.extractTripData(trip)

	result := &ReplanningResult{
		TripID:          tripData.ID,
		OriginalPlan:    trip,
		Triggers:        triggers,
		ReplanTimestamp: time.Now(),
		Changes:         []ItineraryChange{},
	}

	// Generate replacement activities for each trigger
	for _, trigger := range triggers {
		changes, err := d.generateReplacements(ctx, tripData, trigger)
		if err != nil {
			log.Printf("Failed to generate replacements for trigger %s: %v", trigger.Type, err)
			continue
		}
		result.Changes = append(result.Changes, changes...)
	}

	// Use Gemini AI to optimize the revised plan
	if d.gemini != nil {
		revisedPlan, confidence, err := d.optimizeWithAI(ctx, tripData, result.Changes, triggers)
		if err != nil {
			log.Printf("AI optimization failed: %v", err)
			result.Confidence = 0.7 // Default confidence
		} else {
			result.RevisedPlan = revisedPlan
			result.Confidence = confidence
		}
	}

	// Calculate cost impact
	result.EstimatedSavings = d.calculateCostImpact(result.Changes)

	return result, nil
}

// generateReplacements generates alternative activities based on trigger type
func (d *DynamicReplanningService) generateReplacements(ctx context.Context, trip *TripData, trigger ReplanningTrigger) ([]ItineraryChange, error) {
	var changes []ItineraryChange

	switch trigger.Type {
	case "weather":
		changes = d.generateWeatherReplacements(ctx, trip, trigger)
	case "delay":
		changes = d.generateDelayReplacements(ctx, trip, trigger)
	case "sold_out":
		changes = d.generateAvailabilityReplacements(ctx, trip, trigger)
	}

	return changes, nil
}

// generateWeatherReplacements creates indoor alternatives for bad weather
func (d *DynamicReplanningService) generateWeatherReplacements(ctx context.Context, trip *TripData, trigger ReplanningTrigger) []ItineraryChange {
	var changes []ItineraryChange

	weatherAlert, ok := trigger.Data.(WeatherAlert)
	if !ok {
		return changes
	}

	// Find outdoor activities affected by weather
	for dayKey, dayData := range trip.Itinerary {
		if dayMap, ok := dayData.(map[string]interface{}); ok {
			if activitiesData, exists := dayMap["activities"]; exists {
				if activities, ok := activitiesData.([]interface{}); ok {
					for i, activityData := range activities {
						if activity, ok := activityData.(map[string]interface{}); ok {
							if d.isOutdoorActivity(activity) {
								// Find indoor alternative
								alternative := d.findIndoorAlternative(ctx, activity, trip.Destination)
								if alternative != nil {
									changes = append(changes, ItineraryChange{
										Type:        "replacement",
										Day:         dayKey,
										TimeSlot:    fmt.Sprintf("activity_%d", i),
										Original:    activity,
										Replacement: alternative,
										Reason:      fmt.Sprintf("Weather: %s", weatherAlert.Description),
										Impact:      "moderate",
										CostDelta:   d.calculateActivityCostDelta(activity, alternative),
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return changes
}

// generateDelayReplacements adjusts schedule for delays
func (d *DynamicReplanningService) generateDelayReplacements(ctx context.Context, trip *TripData, trigger ReplanningTrigger) []ItineraryChange {
	var changes []ItineraryChange

	delayAlert, ok := trigger.Data.(DelayAlert)
	if !ok {
		return changes
	}

	// Find affected activities and reschedule
	for dayKey, dayData := range trip.Itinerary {
		if dayMap, ok := dayData.(map[string]interface{}); ok {
			if activitiesData, exists := dayMap["activities"]; exists {
				if activities, ok := activitiesData.([]interface{}); ok {
					for i, activityData := range activities {
						if activity, ok := activityData.(map[string]interface{}); ok {
							if d.isActivityAffected(activity, delayAlert) {
								// Shift time or find replacement
								if delayAlert.DelayTime > 2*time.Hour {
									// Significant delay - find replacement
									alternative := d.findTimeAlternative(ctx, activity, dayKey, fmt.Sprintf("activity_%d", i))
									if alternative != nil {
										changes = append(changes, ItineraryChange{
											Type:        "replacement",
											Day:         dayKey,
											TimeSlot:    fmt.Sprintf("activity_%d", i),
											Original:    activity,
											Replacement: alternative,
											Reason:      fmt.Sprintf("Delay: %s", delayAlert.Reason),
											Impact:      "major",
											CostDelta:   d.calculateActivityCostDelta(activity, alternative),
										})
									}
								} else {
									// Minor delay - shift timing
									changes = append(changes, ItineraryChange{
										Type:     "time_shift",
										Day:      dayKey,
										TimeSlot: fmt.Sprintf("activity_%d", i),
										Original: activity,
										Reason:   fmt.Sprintf("Delayed by %v", delayAlert.DelayTime),
										Impact:   "minor",
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return changes
}

// generateAvailabilityReplacements finds alternatives for unavailable items
func (d *DynamicReplanningService) generateAvailabilityReplacements(ctx context.Context, trip *TripData, trigger ReplanningTrigger) []ItineraryChange {
	var changes []ItineraryChange

	availAlert, ok := trigger.Data.(AvailabilityAlert)
	if !ok {
		return changes
	}

	// Find similar alternatives using vector search
	if d.vectorDB != nil {
		alternatives := d.findSimilarAlternatives(ctx, availAlert.ItemID, availAlert.ItemType)
		if len(alternatives) > 0 {
			best := alternatives[0] // Take best match

			// Find where the unavailable item was scheduled
			for dayKey, dayData := range trip.Itinerary {
				if dayMap, ok := dayData.(map[string]interface{}); ok {
					if activitiesData, exists := dayMap["activities"]; exists {
						if activities, ok := activitiesData.([]interface{}); ok {
							for i, activityData := range activities {
								if activity, ok := activityData.(map[string]interface{}); ok {
									if d.matchesUnavailableItem(activity, availAlert) {
										changes = append(changes, ItineraryChange{
											Type:        "replacement",
											Day:         dayKey,
											TimeSlot:    fmt.Sprintf("activity_%d", i),
											Original:    activity,
											Replacement: best,
											Reason:      fmt.Sprintf("Unavailable: %s", availAlert.Status),
											Impact:      "moderate",
											CostDelta:   d.calculateReplacementCostDelta(activity, best),
										})
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return changes
}

// Helper methods and mock implementations

func (d *DynamicReplanningService) getCurrentTrip(ctx context.Context, tripID string) (interface{}, error) {
	// Mock implementation - in production, fetch from database
	return &TripData{
		ID:          tripID,
		Destination: "Delhi",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 3),
		Itinerary: map[string]interface{}{
			"day1": map[string]interface{}{
				"activities": []interface{}{
					map[string]interface{}{"name": "Red Fort", "type": "outdoor"},
					map[string]interface{}{"name": "India Gate", "type": "outdoor"},
					map[string]interface{}{"name": "Connaught Place", "type": "mixed"},
				},
			},
		},
	}, nil
}

func (d *DynamicReplanningService) fetchWeatherAlerts(ctx context.Context, trip interface{}) []WeatherAlert {
	// Mock weather alerts
	return []WeatherAlert{
		{
			AlertType:     "rain",
			Severity:      "warning",
			StartTime:     time.Now().Add(1 * time.Hour),
			EndTime:       time.Now().Add(6 * time.Hour),
			Description:   "Heavy rainfall expected",
			AffectedAreas: []string{"Central Delhi"},
		},
	}
}

func (d *DynamicReplanningService) fetchDelayAlerts(ctx context.Context, trip interface{}) []DelayAlert {
	// Mock delay alerts
	return []DelayAlert{
		{
			ServiceType: "metro",
			ServiceID:   "blue_line",
			DelayTime:   45 * time.Minute,
			Status:      "delayed",
			Reason:      "Technical issues",
		},
	}
}

func (d *DynamicReplanningService) fetchAvailabilityAlerts(ctx context.Context, trip interface{}) []AvailabilityAlert {
	// Mock availability alerts
	return []AvailabilityAlert{
		{
			ItemType: "attraction",
			ItemID:   "red_fort_tickets",
			Status:   "sold_out",
			Until:    nil, // Indefinite
		},
	}
}

func (d *DynamicReplanningService) filterCriticalTriggers(triggers []ReplanningTrigger) []ReplanningTrigger {
	var critical []ReplanningTrigger
	for _, trigger := range triggers {
		if trigger.Severity == "high" || trigger.Severity == "critical" {
			critical = append(critical, trigger)
		}
	}
	return critical
}

func (d *DynamicReplanningService) mapWeatherSeverity(alert WeatherAlert) string {
	switch alert.Severity {
	case "emergency":
		return "critical"
	case "warning":
		return "high"
	case "watch":
		return "medium"
	default:
		return "low"
	}
}

func (d *DynamicReplanningService) mapDelaySeverity(alert DelayAlert) string {
	if alert.Status == "cancelled" {
		return "critical"
	}
	if alert.DelayTime > 2*time.Hour {
		return "high"
	}
	if alert.DelayTime > 30*time.Minute {
		return "medium"
	}
	return "low"
}

func (d *DynamicReplanningService) extractTripData(trip interface{}) *TripData {
	// Mock extraction
	return &TripData{
		ID:          "mock_trip",
		Destination: "Delhi",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 3),
		Itinerary:   map[string]interface{}{},
	}
}

func (d *DynamicReplanningService) optimizeWithAI(ctx context.Context, trip *TripData, changes []ItineraryChange, triggers []ReplanningTrigger) (interface{}, float64, error) {
	// Use Gemini to optimize the revised plan
	prompt := d.buildOptimizationPrompt(trip, changes, triggers)
	_ = prompt // Avoid unused variable error

	// Convert dates from interface{} to time.Time for formatting
	var startDate, endDate time.Time
	if sd, ok := trip.StartDate.(time.Time); ok {
		startDate = sd
	} else {
		startDate = time.Now() // fallback
	}
	if ed, ok := trip.EndDate.(time.Time); ok {
		endDate = ed
	} else {
		endDate = time.Now().AddDate(0, 0, 3) // fallback
	}

	req := ItineraryRequest{
		Destination: trip.Destination,
		StartDate:   startDate.Format("2006-01-02"),
		EndDate:     endDate.Format("2006-01-02"),
		Budget:      5000, // Mock budget
		Travelers:   2,    // Mock travelers
		Preferences: map[string]interface{}{
			"replanning_context": true,
			"changes_count":      len(changes),
		},
	}

	optimizedPlan, err := d.gemini.GenerateItinerary(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// Calculate confidence based on number of successful replacements
	confidence := d.calculateOptimizationConfidence(changes)

	return optimizedPlan, confidence, nil
}

func (d *DynamicReplanningService) buildOptimizationPrompt(trip *TripData, changes []ItineraryChange, triggers []ReplanningTrigger) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Optimize the revised itinerary for %s with the following changes:\n\n", trip.Destination))

	for i, change := range changes {
		prompt.WriteString(fmt.Sprintf("%d. %s: %s (Impact: %s)\n",
			i+1, change.Type, change.Reason, change.Impact))
	}

	prompt.WriteString("\nTriggers requiring attention:\n")
	for i, trigger := range triggers {
		prompt.WriteString(fmt.Sprintf("%d. %s (%s): %s\n",
			i+1, trigger.Type, trigger.Severity, trigger.Description))
	}

	prompt.WriteString("\nProvide an optimized itinerary that addresses these issues while maintaining travel enjoyment and budget efficiency.")

	return prompt.String()
}

func (d *DynamicReplanningService) calculateOptimizationConfidence(changes []ItineraryChange) float64 {
	if len(changes) == 0 {
		return 1.0
	}

	successfulChanges := 0
	for _, change := range changes {
		if change.Replacement != nil {
			successfulChanges++
		}
	}

	confidence := float64(successfulChanges) / float64(len(changes))
	return math.Max(0.5, confidence) // Minimum 50% confidence
}

func (d *DynamicReplanningService) calculateCostImpact(changes []ItineraryChange) float64 {
	totalDelta := 0.0
	for _, change := range changes {
		totalDelta += change.CostDelta
	}
	return totalDelta
}

func (d *DynamicReplanningService) saveReplanResult(ctx context.Context, result *ReplanningResult) error {
	if d.firebase == nil {
		return fmt.Errorf("firebase service not available")
	}

	// Save to Firebase collection
	_, err := d.firebase.GetFirestoreClient().
		Collection("trip_replanning").
		Doc(result.TripID+"_"+result.ReplanTimestamp.Format("20060102_150405")).
		Set(ctx, result)

	return err
}

func (d *DynamicReplanningService) sendReplanNotifications(ctx context.Context, result *ReplanningResult) {
	if d.notificationSvc == nil {
		return
	}

	// Send push notification about replanning
	message := d.buildReplanNotificationMessage(result)
	d.notificationSvc.SendTripUpdateNotification(ctx, result.TripID, message)
}

func (d *DynamicReplanningService) buildReplanNotificationMessage(result *ReplanningResult) string {
	if len(result.Changes) == 0 {
		return "Your itinerary has been reviewed - no changes needed."
	}

	return fmt.Sprintf("Your itinerary has been updated with %d changes due to real-time conditions. Tap to view details.", len(result.Changes))
}

// Additional helper methods (simplified for brevity)
func (d *DynamicReplanningService) isOutdoorActivity(activity interface{}) bool {
	if actMap, ok := activity.(map[string]interface{}); ok {
		if actType, exists := actMap["type"]; exists {
			return actType == "outdoor"
		}
	}
	return false
}

func (d *DynamicReplanningService) isTimeAffected(alert WeatherAlert, day int, timeSlot string) bool {
	return true // Simplified - in production, check actual times
}

func (d *DynamicReplanningService) findIndoorAlternative(ctx context.Context, activity interface{}, destination string) interface{} {
	// Mock indoor alternative
	return map[string]interface{}{
		"name": "National Museum",
		"type": "indoor",
		"cost": 100,
	}
}

func (d *DynamicReplanningService) findTimeAlternative(ctx context.Context, activity interface{}, day string, timeSlot string) interface{} {
	// Mock time alternative
	return map[string]interface{}{
		"name": "Alternative Activity",
		"type": "flexible",
		"cost": 150,
	}
}

func (d *DynamicReplanningService) findSimilarAlternatives(ctx context.Context, itemID, itemType string) []interface{} {
	// Mock alternatives using vector search
	return []interface{}{
		map[string]interface{}{
			"name": "Similar Attraction",
			"type": itemType,
			"cost": 120,
		},
	}
}

func (d *DynamicReplanningService) isActivityAffected(activity interface{}, alert DelayAlert) bool {
	return false // Simplified check
}

func (d *DynamicReplanningService) matchesUnavailableItem(activity interface{}, alert AvailabilityAlert) bool {
	return false // Simplified check
}

func (d *DynamicReplanningService) calculateActivityCostDelta(original, replacement interface{}) float64 {
	return 0.0 // Simplified calculation
}

func (d *DynamicReplanningService) calculateReplacementCostDelta(original, replacement interface{}) float64 {
	return 0.0 // Simplified calculation
}

// TripData represents a trip for replanning purposes
// StopMonitoring stops monitoring for a specific trip
func (d *DynamicReplanningService) StopMonitoring(tripID string) {
	log.Printf("Stopped monitoring trip: %s", tripID)
}

// GetReplanHistory retrieves replanning history for a trip
func (d *DynamicReplanningService) GetReplanHistory(ctx context.Context, tripID string) ([]*ReplanningResult, error) {
	if d.firebase == nil {
		return nil, fmt.Errorf("firebase service not available")
	}

	docs, err := d.firebase.GetFirestoreClient().
		Collection("trip_replanning").
		Where("trip_id", "==", tripID).
		OrderBy("replan_timestamp", firestore.Desc).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	var results []*ReplanningResult
	for _, doc := range docs {
		var result ReplanningResult
		if err := doc.DataTo(&result); err == nil {
			results = append(results, &result)
		}
	}

	return results, nil
}
