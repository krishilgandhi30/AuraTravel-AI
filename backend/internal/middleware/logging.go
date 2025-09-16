package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware provides request logging
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ErrorHandlingMiddleware handles errors and panics
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
			c.JSON(500, gin.H{
				"error": "Internal server error",
			})
		}
		c.AbortWithStatus(500)
	})
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// TODO: Implement proper rate limiting with Redis or in-memory store
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()
	})
}
