package middleware

import (
	"net/http"
	"sync"
	"time"

	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds rate limiters for different clients
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(float64(requestsPerMinute) / 60.0), // Convert to per-second rate
		burst:    requestsPerMinute,
	}
}

// getLimiter returns a rate limiter for a client
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// Cleanup removes old limiters periodically
func (rl *RateLimiter) Cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In a production system, track last access time and remove inactive limiters
			// For simplicity, we'll keep them for now
			rl.mu.Unlock()
		}
	}()
}

// RateLimit creates a rate limiting middleware
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use API key ID or IP address as the rate limit key
		key := c.ClientIP()

		// If API key is present, use it for rate limiting
		if apiKey, exists := c.Get("api_key"); exists {
			if ak, ok := apiKey.(*models.APIKey); ok {
				key = ak.ID.String()
			}
		}

		// If username is present (admin), use it
		if username, exists := c.Get("username"); exists {
			if u, ok := username.(string); ok {
				key = "admin:" + u
			}
		}

		limiter := limiter.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error: "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
