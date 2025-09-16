package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"firebase.google.com/go/v4/messaging"
)

// NotificationService handles push notifications via Firebase Cloud Messaging
type NotificationService struct {
	messagingClient *messaging.Client
	firebase        *FirebaseService
	enabled         bool
}

// NewNotificationService creates a new notification service
func NewNotificationService(firebase *FirebaseService) (*NotificationService, error) {
	if firebase == nil {
		log.Println("Warning: Firebase service not available, notifications disabled")
		return &NotificationService{enabled: false}, nil
	}

	messagingClient, err := firebase.GetMessagingClient()
	if err != nil {
		log.Printf("Warning: Failed to initialize FCM client: %v", err)
		return &NotificationService{
			firebase: firebase,
			enabled:  false,
		}, nil
	}

	return &NotificationService{
		messagingClient: messagingClient,
		firebase:        firebase,
		enabled:         true,
	}, nil
}

// NotificationType represents different types of notifications
type NotificationType string

const (
	WeatherAlertType NotificationType = "weather_alert"
	ItineraryUpdate  NotificationType = "itinerary_update"
	TripReminder     NotificationType = "trip_reminder"
	DelayAlertType   NotificationType = "delay_alert"
	BookingConfirm   NotificationType = "booking_confirmation"
	GeneralUpdate    NotificationType = "general_update"
	EmergencyAlert   NotificationType = "emergency_alert"
)

// NotificationPriority represents notification priority levels
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// NotificationRequest represents a notification to be sent
type NotificationRequest struct {
	UserID       string               `json:"user_id"`
	TripID       string               `json:"trip_id,omitempty"`
	Type         NotificationType     `json:"type"`
	Priority     NotificationPriority `json:"priority"`
	Title        string               `json:"title"`
	Body         string               `json:"body"`
	Data         map[string]string    `json:"data,omitempty"`
	ImageURL     string               `json:"image_url,omitempty"`
	ActionURL    string               `json:"action_url,omitempty"`
	ScheduleTime *time.Time           `json:"schedule_time,omitempty"`
	Language     string               `json:"language,omitempty"`
}

// UserDeviceToken represents a user's FCM device token
type UserDeviceToken struct {
	UserID      string    `json:"user_id"`
	DeviceToken string    `json:"device_token"`
	DeviceType  string    `json:"device_type"` // ios, android, web
	Language    string    `json:"language"`
	Timezone    string    `json:"timezone"`
	Active      bool      `json:"active"`
	LastUsed    time.Time `json:"last_used"`
	CreatedAt   time.Time `json:"created_at"`
}

// NotificationTemplate represents localized notification templates
type NotificationTemplate struct {
	Type      NotificationType  `json:"type"`
	Language  string            `json:"language"`
	TitleTmpl string            `json:"title_template"`
	BodyTmpl  string            `json:"body_template"`
	Variables map[string]string `json:"variables,omitempty"`
}

// RegisterDeviceToken registers a user's device for notifications
func (n *NotificationService) RegisterDeviceToken(ctx context.Context, userID, deviceToken, deviceType string) error {
	if !n.enabled {
		return fmt.Errorf("notification service not enabled")
	}

	token := &UserDeviceToken{
		UserID:      userID,
		DeviceToken: deviceToken,
		DeviceType:  deviceType,
		Language:    "en", // Default language
		Timezone:    "UTC",
		Active:      true,
		LastUsed:    time.Now(),
		CreatedAt:   time.Now(),
	}

	// Store token in Firestore
	_, err := n.firebase.GetFirestoreClient().
		Collection("user_device_tokens").
		Doc(fmt.Sprintf("%s_%s", userID, deviceType)).
		Set(ctx, token)

	if err != nil {
		return fmt.Errorf("failed to store device token: %w", err)
	}

	log.Printf("Registered device token for user %s", userID)
	return nil
}

// SendNotification sends a single notification to a user
func (n *NotificationService) SendNotification(ctx context.Context, req *NotificationRequest) error {
	if !n.enabled {
		log.Printf("Notification service disabled, skipping: %s", req.Title)
		return nil
	}

	// Get user's device tokens
	tokens, err := n.getUserDeviceTokens(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user device tokens: %w", err)
	}

	if len(tokens) == 0 {
		log.Printf("No device tokens found for user %s", req.UserID)
		return nil
	}

	// Apply localization if needed
	localizedReq := n.localizeNotification(req, tokens[0].Language)

	// Build FCM message
	fcmMessage := n.buildFCMMessage(localizedReq, tokens)

	// Send message
	response, err := n.messagingClient.SendMulticast(ctx, fcmMessage)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	// Log results
	log.Printf("Sent notification to %d devices, %d successful, %d failed",
		len(tokens), response.SuccessCount, response.FailureCount)

	// Handle failed tokens
	n.handleFailedTokens(ctx, response, tokens)

	// Store notification history
	n.storeNotificationHistory(ctx, req, response)

	return nil
}

// SendWeatherAlert sends weather-related notifications
func (n *NotificationService) SendWeatherAlert(ctx context.Context, userID, tripID string, alert interface{}) error {
	weatherAlert, ok := alert.(WeatherAlert)
	if !ok {
		return fmt.Errorf("invalid weather alert type")
	}

	req := &NotificationRequest{
		UserID:   userID,
		TripID:   tripID,
		Type:     WeatherAlertType,
		Priority: n.mapWeatherPriority(weatherAlert.Severity),
		Title:    n.getWeatherAlertTitle(weatherAlert),
		Body:     n.getWeatherAlertBody(weatherAlert),
		Data: map[string]string{
			"trip_id":    tripID,
			"alert_type": weatherAlert.AlertType,
			"severity":   weatherAlert.Severity,
			"start_time": weatherAlert.StartTime.Format(time.RFC3339),
			"end_time":   weatherAlert.EndTime.Format(time.RFC3339),
		},
		ActionURL: fmt.Sprintf("/trips/%s?tab=weather", tripID),
	}

	return n.SendNotification(ctx, req)
}

// SendTripUpdateNotification sends itinerary update notifications
func (n *NotificationService) SendTripUpdateNotification(ctx context.Context, tripID, message string) error {
	// Get all users for this trip
	userIDs, err := n.getTripUserIDs(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to get trip users: %w", err)
	}

	for _, userID := range userIDs {
		req := &NotificationRequest{
			UserID:   userID,
			TripID:   tripID,
			Type:     ItineraryUpdate,
			Priority: PriorityHigh,
			Title:    "Trip Update",
			Body:     message,
			Data: map[string]string{
				"trip_id": tripID,
				"type":    "itinerary_update",
			},
			ActionURL: fmt.Sprintf("/trips/%s", tripID),
		}

		if err := n.SendNotification(ctx, req); err != nil {
			log.Printf("Failed to send trip update to user %s: %v", userID, err)
		}
	}

	return nil
}

// SendDelayAlert sends transportation/event delay notifications
func (n *NotificationService) SendDelayAlert(ctx context.Context, userID, tripID string, alert interface{}) error {
	delayAlert, ok := alert.(DelayAlert)
	if !ok {
		return fmt.Errorf("invalid delay alert type")
	}

	req := &NotificationRequest{
		UserID:   userID,
		TripID:   tripID,
		Type:     DelayAlertType,
		Priority: n.mapDelayPriority(delayAlert),
		Title:    n.getDelayAlertTitle(delayAlert),
		Body:     n.getDelayAlertBody(delayAlert),
		Data: map[string]string{
			"trip_id":      tripID,
			"service_type": delayAlert.ServiceType,
			"service_id":   delayAlert.ServiceID,
			"status":       delayAlert.Status,
			"delay_time":   delayAlert.DelayTime.String(),
		},
		ActionURL: fmt.Sprintf("/trips/%s?tab=transportation", tripID),
	}

	return n.SendNotification(ctx, req)
}

// SendTripReminder sends trip reminders (e.g., "Trip starts tomorrow")
func (n *NotificationService) SendTripReminder(ctx context.Context, userID, tripID string, reminderType string, timeUntil time.Duration) error {
	req := &NotificationRequest{
		UserID:   userID,
		TripID:   tripID,
		Type:     TripReminder,
		Priority: PriorityNormal,
		Title:    n.getTripReminderTitle(reminderType, timeUntil),
		Body:     n.getTripReminderBody(reminderType, timeUntil),
		Data: map[string]string{
			"trip_id":       tripID,
			"reminder_type": reminderType,
			"time_until":    timeUntil.String(),
		},
		ActionURL: fmt.Sprintf("/trips/%s", tripID),
	}

	return n.SendNotification(ctx, req)
}

// SendBookingConfirmation sends booking confirmation notifications
func (n *NotificationService) SendBookingConfirmation(ctx context.Context, userID, tripID, bookingType, confirmationNumber string) error {
	req := &NotificationRequest{
		UserID:   userID,
		TripID:   tripID,
		Type:     BookingConfirm,
		Priority: PriorityHigh,
		Title:    fmt.Sprintf("%s Booking Confirmed", strings.Title(bookingType)),
		Body:     fmt.Sprintf("Your %s booking has been confirmed. Confirmation #%s", bookingType, confirmationNumber),
		Data: map[string]string{
			"trip_id":             tripID,
			"booking_type":        bookingType,
			"confirmation_number": confirmationNumber,
		},
		ActionURL: fmt.Sprintf("/trips/%s/bookings", tripID),
	}

	return n.SendNotification(ctx, req)
}

// ScheduleNotification schedules a notification for future delivery
func (n *NotificationService) ScheduleNotification(ctx context.Context, req *NotificationRequest) error {
	if req.ScheduleTime == nil {
		return n.SendNotification(ctx, req)
	}

	// Store scheduled notification in Firestore
	_, err := n.firebase.GetFirestoreClient().
		Collection("scheduled_notifications").
		Doc(fmt.Sprintf("%s_%d", req.UserID, req.ScheduleTime.Unix())).
		Set(ctx, req)

	if err != nil {
		return fmt.Errorf("failed to schedule notification: %w", err)
	}

	log.Printf("Scheduled notification for user %s at %v", req.UserID, req.ScheduleTime)
	return nil
}

// ProcessScheduledNotifications processes notifications that are due
func (n *NotificationService) ProcessScheduledNotifications(ctx context.Context) error {
	if !n.enabled {
		return nil
	}

	now := time.Now()

	// Query scheduled notifications that are due
	docs, err := n.firebase.GetFirestoreClient().
		Collection("scheduled_notifications").
		Where("schedule_time", "<=", now).
		Documents(ctx).
		GetAll()

	if err != nil {
		return fmt.Errorf("failed to query scheduled notifications: %w", err)
	}

	for _, doc := range docs {
		var req NotificationRequest
		if err := doc.DataTo(&req); err != nil {
			log.Printf("Failed to parse scheduled notification: %v", err)
			continue
		}

		// Send the notification
		if err := n.SendNotification(ctx, &req); err != nil {
			log.Printf("Failed to send scheduled notification: %v", err)
			continue
		}

		// Delete the scheduled notification
		doc.Ref.Delete(ctx)
	}

	if len(docs) > 0 {
		log.Printf("Processed %d scheduled notifications", len(docs))
	}

	return nil
}

// Helper methods

func (n *NotificationService) getUserDeviceTokens(ctx context.Context, userID string) ([]UserDeviceToken, error) {
	docs, err := n.firebase.GetFirestoreClient().
		Collection("user_device_tokens").
		Where("user_id", "==", userID).
		Where("active", "==", true).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	var tokens []UserDeviceToken
	for _, doc := range docs {
		var token UserDeviceToken
		if err := doc.DataTo(&token); err == nil {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (n *NotificationService) buildFCMMessage(req *NotificationRequest, tokens []UserDeviceToken) *messaging.MulticastMessage {
	// Extract device tokens
	var deviceTokens []string
	for _, token := range tokens {
		deviceTokens = append(deviceTokens, token.DeviceToken)
	}

	// Build notification
	notification := &messaging.Notification{
		Title: req.Title,
		Body:  req.Body,
	}

	if req.ImageURL != "" {
		notification.ImageURL = req.ImageURL
	}

	// Build data payload
	data := make(map[string]string)
	if req.Data != nil {
		for k, v := range req.Data {
			data[k] = v
		}
	}

	// Add standard fields
	data["type"] = string(req.Type)
	data["priority"] = string(req.Priority)
	if req.TripID != "" {
		data["trip_id"] = req.TripID
	}
	if req.ActionURL != "" {
		data["action_url"] = req.ActionURL
	}

	// Build Android config
	androidConfig := &messaging.AndroidConfig{
		Priority: n.mapPriorityToAndroid(req.Priority),
		Notification: &messaging.AndroidNotification{
			Title:       req.Title,
			Body:        req.Body,
			ClickAction: req.ActionURL,
			Sound:       n.getSoundForPriority(req.Priority),
			Priority:    n.mapPriorityToAndroidNotification(req.Priority),
			Visibility:  messaging.VisibilityPublic,
		},
	}

	// Build APNS config for iOS
	apnsConfig := &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: req.Title,
					Body:  req.Body,
				},
				Sound: n.getSoundForPriority(req.Priority),
				Badge: n.getBadgeCount(req.UserID),
			},
		},
	}

	return &messaging.MulticastMessage{
		Tokens:       deviceTokens,
		Notification: notification,
		Data:         data,
		Android:      androidConfig,
		APNS:         apnsConfig,
	}
}

func (n *NotificationService) localizeNotification(req *NotificationRequest, language string) *NotificationRequest {
	// Skip localization if already in desired language or default language
	if language == "" || language == "en" {
		return req
	}

	// Get localized template
	template := n.getNotificationTemplate(req.Type, language)
	if template == nil {
		return req // Fallback to original
	}

	// Clone request and apply localization
	localizedReq := *req
	localizedReq.Title = n.applyTemplate(template.TitleTmpl, req.Data)
	localizedReq.Body = n.applyTemplate(template.BodyTmpl, req.Data)
	localizedReq.Language = language

	return &localizedReq
}

func (n *NotificationService) getNotificationTemplate(notifType NotificationType, language string) *NotificationTemplate {
	// In production, load from database or configuration
	templates := map[string]*NotificationTemplate{
		string(WeatherAlertType) + "_hi": {
			Type:      WeatherAlertType,
			Language:  "hi",
			TitleTmpl: "मौसम चेतावनी",
			BodyTmpl:  "आपकी यात्रा के लिए मौसम की चेतावनी: {{description}}",
		},
		string(ItineraryUpdate) + "_hi": {
			Type:      ItineraryUpdate,
			Language:  "hi",
			TitleTmpl: "यात्रा अपडेट",
			BodyTmpl:  "आपका यात्रा कार्यक्रम अपडेट किया गया है",
		},
	}

	key := string(notifType) + "_" + language
	return templates[key]
}

func (n *NotificationService) applyTemplate(template string, data map[string]string) string {
	result := template
	if data != nil {
		for key, value := range data {
			placeholder := "{{" + key + "}}"
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}
	return result
}

func (n *NotificationService) handleFailedTokens(ctx context.Context, response *messaging.BatchResponse, tokens []UserDeviceToken) {
	for i, resp := range response.Responses {
		if !resp.Success && i < len(tokens) {
			// Handle invalid tokens
			if messaging.IsRegistrationTokenNotRegistered(resp.Error) ||
				messaging.IsInvalidArgument(resp.Error) {
				// Deactivate invalid token
				n.deactivateDeviceToken(ctx, tokens[i].DeviceToken)
			}
		}
	}
}

func (n *NotificationService) deactivateDeviceToken(ctx context.Context, deviceToken string) {
	// Mark token as inactive
	_, err := n.firebase.GetFirestoreClient().
		Collection("user_device_tokens").
		Where("device_token", "==", deviceToken).
		Documents(ctx).
		GetAll()

	if err != nil {
		log.Printf("Failed to query device token for deactivation: %v", err)
		return
	}
}

func (n *NotificationService) storeNotificationHistory(ctx context.Context, req *NotificationRequest, response *messaging.BatchResponse) {
	history := map[string]interface{}{
		"user_id":       req.UserID,
		"trip_id":       req.TripID,
		"type":          req.Type,
		"title":         req.Title,
		"body":          req.Body,
		"sent_at":       time.Now(),
		"success_count": response.SuccessCount,
		"failure_count": response.FailureCount,
	}

	n.firebase.GetFirestoreClient().
		Collection("notification_history").
		Add(ctx, history)
}

func (n *NotificationService) getTripUserIDs(ctx context.Context, tripID string) ([]string, error) {
	// Mock implementation - get users associated with trip
	return []string{"user_123"}, nil
}

// Priority and alert mapping methods

func (n *NotificationService) mapWeatherPriority(severity string) NotificationPriority {
	switch severity {
	case "emergency":
		return PriorityCritical
	case "warning":
		return PriorityHigh
	case "watch":
		return PriorityNormal
	default:
		return PriorityLow
	}
}

func (n *NotificationService) mapDelayPriority(alert DelayAlert) NotificationPriority {
	if alert.Status == "cancelled" {
		return PriorityCritical
	}
	if alert.DelayTime > time.Hour {
		return PriorityHigh
	}
	return PriorityNormal
}

func (n *NotificationService) mapPriorityToAndroid(priority NotificationPriority) string {
	switch priority {
	case PriorityCritical, PriorityHigh:
		return "high"
	default:
		return "normal"
	}
}

func (n *NotificationService) mapPriorityToAndroidNotification(priority NotificationPriority) messaging.AndroidNotificationPriority {
	switch priority {
	case PriorityCritical, PriorityHigh:
		return messaging.PriorityHigh
	default:
		return messaging.PriorityDefault
	}
}

func (n *NotificationService) getSoundForPriority(priority NotificationPriority) string {
	switch priority {
	case PriorityCritical:
		return "emergency.wav"
	case PriorityHigh:
		return "alert.wav"
	default:
		return "default"
	}
}

func (n *NotificationService) getBadgeCount(userID string) *int {
	// In production, calculate unread notification count
	count := 1
	return &count
}

// Alert message generators

func (n *NotificationService) getWeatherAlertTitle(alert WeatherAlert) string {
	return fmt.Sprintf("Weather Alert: %s", strings.Title(alert.AlertType))
}

func (n *NotificationService) getWeatherAlertBody(alert WeatherAlert) string {
	return fmt.Sprintf("%s expected from %s to %s. %s",
		strings.Title(alert.AlertType),
		alert.StartTime.Format("3:04 PM"),
		alert.EndTime.Format("3:04 PM"),
		alert.Description)
}

func (n *NotificationService) getDelayAlertTitle(alert DelayAlert) string {
	return fmt.Sprintf("%s %s", strings.Title(alert.ServiceType), strings.Title(alert.Status))
}

func (n *NotificationService) getDelayAlertBody(alert DelayAlert) string {
	if alert.Status == "cancelled" {
		return fmt.Sprintf("Your %s (%s) has been cancelled: %s", alert.ServiceType, alert.ServiceID, alert.Reason)
	}
	return fmt.Sprintf("Your %s (%s) is delayed by %v: %s", alert.ServiceType, alert.ServiceID, alert.DelayTime, alert.Reason)
}

func (n *NotificationService) getTripReminderTitle(reminderType string, timeUntil time.Duration) string {
	switch reminderType {
	case "departure":
		return "Trip Departure Reminder"
	case "checkin":
		return "Check-in Reminder"
	case "activity":
		return "Upcoming Activity"
	default:
		return "Trip Reminder"
	}
}

func (n *NotificationService) getTripReminderBody(reminderType string, timeUntil time.Duration) string {
	timeStr := n.formatDuration(timeUntil)
	switch reminderType {
	case "departure":
		return fmt.Sprintf("Your trip starts in %s. Have a great journey!", timeStr)
	case "checkin":
		return fmt.Sprintf("Don't forget to check in for your flight/hotel in %s", timeStr)
	case "activity":
		return fmt.Sprintf("Your next activity starts in %s", timeStr)
	default:
		return fmt.Sprintf("Trip reminder: %s", timeStr)
	}
}

func (n *NotificationService) formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	if hours >= 24 {
		days := hours / 24
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
	if hours >= 1 {
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	minutes := int(d.Minutes())
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}

// Shutdown gracefully shuts down the notification service
func (n *NotificationService) Shutdown(ctx context.Context) error {
	log.Println("Notification service shut down successfully")
	return nil
}
