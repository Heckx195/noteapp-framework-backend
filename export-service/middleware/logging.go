package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware logs request details with RequestID, Time, and Duration
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		// Add RequestID to the context
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("RequestID", requestID)

		// Process the request
		c.Next()

		// Log the request details
		duration := time.Since(start)
		log.Printf("[%s] [RequestID: %s] %s %s - %d - Duration: %v",
			start.Format(time.RFC3339),
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
