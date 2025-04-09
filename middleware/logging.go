package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Log request details
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Headers: %v", c.Request.Header)

		// Process request
		c.Next()

		// Log response details
		duration := time.Since(start)
		log.Printf("Response: %d %s (%s)", c.Writer.Status(), c.Request.URL.Path, duration)
	}
}
