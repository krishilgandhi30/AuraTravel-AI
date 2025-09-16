package handlers

import (
	"net/http"
	"strconv"
	"time"

	"auratravel-backend/internal/database"
	"auratravel-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related endpoints
type UserHandler struct{}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// RegisterUserRequest represents the user registration request
type RegisterUserRequest struct {
	Email             string `json:"email" binding:"required,email"`
	DisplayName       string `json:"display_name"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	PreferredCurrency string `json:"preferred_currency"`
	PreferredLanguage string `json:"preferred_language"`
}

// UpdateUserRequest represents the user update request
type UpdateUserRequest struct {
	DisplayName       string     `json:"display_name"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	PhoneNumber       string     `json:"phone_number"`
	DateOfBirth       *time.Time `json:"date_of_birth"`
	Nationality       string     `json:"nationality"`
	PreferredCurrency string     `json:"preferred_currency"`
	PreferredLanguage string     `json:"preferred_language"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Create new user
	user := models.User{
		ID:                uuid.New().String(),
		Email:             req.Email,
		DisplayName:       req.DisplayName,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PhoneNumber:       req.PhoneNumber,
		PreferredCurrency: req.PreferredCurrency,
		PreferredLanguage: req.PreferredLanguage,
		EmailVerified:     false,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Set defaults
	if user.PreferredCurrency == "" {
		user.PreferredCurrency = "USD"
	}
	if user.PreferredLanguage == "" {
		user.PreferredLanguage = "en"
	}

	// Save to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var user models.User
	if err := database.DB.Preload("TravelPreferences").Preload("EmergencyContact").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update last login
	user.LastLoginAt = &[]time.Time{time.Now()}[0]
	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update fields
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.Nationality != "" {
		user.Nationality = req.Nationality
	}
	if req.PreferredCurrency != "" {
		user.PreferredCurrency = req.PreferredCurrency
	}
	if req.PreferredLanguage != "" {
		user.PreferredLanguage = req.PreferredLanguage
	}

	user.UpdatedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Delete the authenticated user's account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/profile [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Soft delete by setting IsActive to false
	user.IsActive = false
	user.UpdatedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// GetUsers godoc
// @Summary Get all users (Admin only)
// @Description Get a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/admin/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	// Get total count
	database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&total)

	// Get users with pagination
	if err := database.DB.Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch users",
			"details": err.Error(),
		})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
