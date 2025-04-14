package middleware

import (
	"ops_service/internal/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware creates a middleware for Prometheus metrics
func MetricsMiddleware(metrics *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment requests counter
		metrics.IncRequestsTotal()
		
		// Measure response time
		startTime := time.Now()
		
		// Process request
		c.Next()
		
		// Observe response time
		metrics.ObserveResponseTime(startTime)
	}
} 