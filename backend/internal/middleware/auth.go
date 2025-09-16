package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens from Firebase
func AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract Bearer token
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := bearerToken[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Validate token with Firebase (we'll implement this when we set up Firebase)
		userID, err := validateFirebaseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user ID in context for use in handlers
		c.Set("userID", userID)
		c.Next()
	})
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
				token := bearerToken[1]
				if userID, err := validateFirebaseToken(token); err == nil {
					c.Set("userID", userID)
				}
			}
		}
		c.Next()
	})
}

// AdminMiddleware checks if user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user is admin (implement admin check logic)
		isAdmin, err := checkAdminStatus(userID.(string))
		if err != nil || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// validateFirebaseToken validates JWT token with Firebase
func validateFirebaseToken(token string) (string, error) {
	// TODO: Implement Firebase token validation
	// For now, return a dummy user ID for testing
	// In production, this should validate with Firebase Auth
	return "test-user-id", nil
}

// checkAdminStatus checks if user has admin privileges
func checkAdminStatus(userID string) (bool, error) {
	// TODO: Implement admin check logic
	// This should check user role in database
	return false, nil
}
