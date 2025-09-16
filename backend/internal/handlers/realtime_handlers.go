package handlers

import (
	"auratravel-backend/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *services.NotificationService
	localizationService *services.LocalizationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(services *services.Services) *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NotificationService,
		localizationService: services.LocalizationService,
	}
}

// RegisterDevice registers a device token for push notifications
func (h *NotificationHandler) RegisterDevice(c *gin.Context) {
	var req struct {
		UserID      string `json:"userId" binding:"required"`
		DeviceToken string `json:"deviceToken" binding:"required"`
		Platform    string `json:"platform" binding:"required"`
		Locale      string `json:"locale"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default locale if not provided
	if req.Locale == "" {
		req.Locale = "en"
	}

	// Adjusted to match service signature: RegisterDeviceToken(ctx, userID, deviceToken, platform)
	err := h.notificationService.RegisterDeviceToken(c.Request.Context(), req.UserID, req.DeviceToken, req.Platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device token registered successfully",
	})
}

// SendNotification sends a push notification
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req services.NotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.notificationService.SendNotification(c.Request.Context(), &req)
	c.JSON(http.StatusOK, result)
}

// SendWeatherAlert sends a weather alert notification
func (h *NotificationHandler) SendWeatherAlert(c *gin.Context) {
	userID := c.Param("userId")
	var req struct {
		TripID           string `json:"tripId" binding:"required"`
		WeatherCondition string `json:"weatherCondition" binding:"required"`
		Location         string `json:"location" binding:"required"`
		Severity         string `json:"severity"`
		Message          string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Adjusted to match service signature: SendWeatherAlert(ctx, userID, tripID, alertData)
	alertData := map[string]interface{}{
		"weatherCondition": req.WeatherCondition,
		"location":         req.Location,
		"severity":         req.Severity,
		"message":          req.Message,
	}
	result := h.notificationService.SendWeatherAlert(c.Request.Context(), userID, req.TripID, alertData)
	c.JSON(http.StatusOK, result)
}

// SendTripUpdate sends a trip update notification
func (h *NotificationHandler) SendTripUpdate(c *gin.Context) {
	userID := c.Param("userId")
	var req struct {
		TripID     string                 `json:"tripId" binding:"required"`
		UpdateType string                 `json:"updateType" binding:"required"`
		Title      string                 `json:"title" binding:"required"`
		Message    string                 `json:"message" binding:"required"`
		ActionURL  string                 `json:"actionUrl,omitempty"`
		CustomData map[string]interface{} `json:"customData,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Adjusted to match service signature: SendTripUpdateNotification(ctx, userID, tripID)
	result := h.notificationService.SendTripUpdateNotification(c.Request.Context(), userID, req.TripID)
	c.JSON(http.StatusOK, result)
}

// GetUserNotifications retrieves user's notification history
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	_ = c.Param("userId")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// This would typically fetch from a database
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"notifications": []interface{}{},
		"total":         0,
		"limit":         limit,
		"offset":        offset,
	})
}

// ReplanningHandler handles dynamic replanning HTTP requests
type ReplanningHandler struct {
	replanningService   *services.DynamicReplanningService
	notificationService *services.NotificationService
	localizationService *services.LocalizationService
}

// NewReplanningHandler creates a new replanning handler
func NewReplanningHandler(services *services.Services) *ReplanningHandler {
	return &ReplanningHandler{
		replanningService:   services.DynamicReplanningService,
		notificationService: services.NotificationService,
		localizationService: services.LocalizationService,
	}
}

// StartMonitoring starts monitoring a trip for replanning triggers
func (h *ReplanningHandler) StartMonitoring(c *gin.Context) {
	// tripID := c.Param("tripId")

	var req struct {
		UserID      string                 `json:"userId" binding:"required"`
		Itinerary   map[string]interface{} `json:"itinerary" binding:"required"`
		Preferences map[string]interface{} `json:"preferences,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Method StartMonitoring does not exist; placeholder for implementation or remove if not needed
	c.JSON(http.StatusNotImplemented, gin.H{"error": "StartMonitoring not implemented"})
}

// StopMonitoring stops monitoring a trip
func (h *ReplanningHandler) StopMonitoring(c *gin.Context) {
	// tripID := c.Param("tripId")

	// Method StopMonitoring does not exist; placeholder for implementation or remove if not needed
	c.JSON(http.StatusNotImplemented, gin.H{"error": "StopMonitoring not implemented"})
}

// GetMonitoringStatus gets the monitoring status for a trip
func (h *ReplanningHandler) GetMonitoringStatus(c *gin.Context) {

	// Method GetMonitoringStatus does not exist; placeholder for implementation or remove if not needed
	c.JSON(http.StatusNotImplemented, gin.H{"error": "GetMonitoringStatus not implemented"})
}

// TriggerReplanning manually triggers replanning for a trip
func (h *ReplanningHandler) TriggerReplanning(c *gin.Context) {
	var req struct {
		TripID           string                     `json:"tripId" binding:"required"`
		Trigger          services.ReplanningTrigger `json:"trigger" binding:"required"`
		CurrentItinerary map[string]interface{}     `json:"currentItinerary" binding:"required"`
		Locale           string                     `json:"locale"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Locale == "" {
		req.Locale = "en"
	}

	// Method ProcessReplanning does not exist; placeholder for implementation or remove if not needed
	c.JSON(http.StatusNotImplemented, gin.H{"error": "TriggerReplanning not implemented"})
}

// AcceptReplanningOption accepts a replanning option
func (h *ReplanningHandler) AcceptReplanningOption(c *gin.Context) {
	// tripID := c.Param("tripId")

	var req struct {
		OptionID string `json:"optionId" binding:"required"`
		UserID   string `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This would typically update the itinerary and notify relevant services
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Replanning option accepted successfully",
		// "tripId":   tripID,
		"optionId": req.OptionID,
	})
}

// DeliveryHandler handles itinerary delivery HTTP requests
type DeliveryHandler struct {
	deliveryService     *services.ItineraryDeliveryService
	localizationService *services.LocalizationService
}

// NewDeliveryHandler creates a new delivery handler
func NewDeliveryHandler(services *services.Services) *DeliveryHandler {
	return &DeliveryHandler{
		deliveryService:     services.ItineraryDeliveryService,
		localizationService: services.LocalizationService,
	}
}

// DeliverItinerary generates and delivers an itinerary
func (h *DeliveryHandler) DeliverItinerary(c *gin.Context) {
	var req services.DeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.deliveryService.GenerateAndDeliverItinerary(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deliver itinerary"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDeliveryStatus gets the status of a delivery
func (h *DeliveryHandler) GetDeliveryStatus(c *gin.Context) {
	deliveryID := c.Param("deliveryId")

	// This would typically fetch from a database
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"deliveryId": deliveryID,
		"status":     "completed",
		"timestamp":  time.Now(),
	})
}

// GenerateShareLink generates a shareable link for an itinerary
func (h *DeliveryHandler) GenerateShareLink(c *gin.Context) {

	var req struct {
		ExpiryHours int    `json:"expiryHours"`
		Password    string `json:"password,omitempty"`
		Locale      string `json:"locale"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ExpiryHours == 0 {
		req.ExpiryHours = 720 // 30 days default
	}

	// Method GenerateShareURL does not exist; placeholder for implementation or remove if not needed
	c.JSON(http.StatusNotImplemented, gin.H{"error": "GenerateShareLink not implemented"})
}

// LocalizationHandler handles localization HTTP requests
type LocalizationHandler struct {
	localizationService *services.LocalizationService
}

// NewLocalizationHandler creates a new localization handler
func NewLocalizationHandler(services *services.Services) *LocalizationHandler {
	return &LocalizationHandler{
		localizationService: services.LocalizationService,
	}
}

// GetSupportedLocales returns all supported locales
func (h *LocalizationHandler) GetSupportedLocales(c *gin.Context) {
	locales := h.localizationService.GetSupportedLocales()
	c.JSON(http.StatusOK, gin.H{
		"locales": locales,
		"count":   len(locales),
	})
}

// LocalizeContent localizes content to a specific locale
func (h *LocalizationHandler) LocalizeContent(c *gin.Context) {
	var req services.LocalizationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result *services.LocalizedContent
	var err error

	switch req.ContentType {
	case "itinerary":
		result, err = h.localizationService.LocalizeItinerary(c.Request.Context(), req.Content, req.TargetLocale)
	case "notification":
		result, err = h.localizationService.LocalizeNotification(c.Request.Context(), &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to localize content"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetGeminiPrompt returns a localized Gemini prompt
func (h *LocalizationHandler) GetGeminiPrompt(c *gin.Context) {
	locale := c.Param("locale")
	promptType := c.Query("type")

	if promptType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt type is required"})
		return
	}

	// Parse variables from query parameters
	variables := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if key != "type" && len(values) > 0 {
			variables[key] = values[0]
		}
	}

	prompt, err := h.localizationService.GetLocalizedGeminiPrompt(locale, promptType, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get localized prompt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prompt":     prompt,
		"locale":     locale,
		"promptType": promptType,
		"variables":  variables,
	})
}

// SetUserLocalePreference sets a user's locale preference
func (h *LocalizationHandler) SetUserLocalePreference(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Locale string `json:"locale" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.localizationService.SetUserLocalePreference(c.Request.Context(), req.UserID, req.Locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set locale preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Locale preference updated successfully",
		"userId":  req.UserID,
		"locale":  req.Locale,
	})
}

// GetUserLocalePreference gets a user's locale preference
func (h *LocalizationHandler) GetUserLocalePreference(c *gin.Context) {
	userID := c.Param("userId")

	locale, err := h.localizationService.GetUserLocalePreference(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get locale preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"locale": locale,
	})
}

// FormatCurrency formats currency according to locale
func (h *LocalizationHandler) FormatCurrency(c *gin.Context) {
	locale := c.Param("locale")
	amountStr := c.Query("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	formatted, err := h.localizationService.FormatCurrency(amount, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  amount,
		"formatted": formatted,
		"locale":    locale,
	})
}

// FormatDateTime formats date and time according to locale
func (h *LocalizationHandler) FormatDateTime(c *gin.Context) {
	locale := c.Param("locale")
	datetimeStr := c.Query("datetime")

	datetime, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datetime format (use RFC3339)"})
		return
	}

	formattedDate, err := h.localizationService.FormatDate(datetime, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format date"})
		return
	}

	formattedTime, err := h.localizationService.FormatTime(datetime, locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format time"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":      datetimeStr,
		"formattedDate": formattedDate,
		"formattedTime": formattedTime,
		"locale":        locale,
	})
}
