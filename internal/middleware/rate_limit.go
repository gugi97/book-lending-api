package middleware

import (
	"book-lending-api/internal/domain"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter provides a simple per‑IP token bucket implementation.
// It stores a limiter per client IP and periodically cleans up old
// entries.  This middleware is intended to protect the API against
// bursts of traffic or abuse.
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter constructs a RateLimiter with the provided rate and
// burst capacity.
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if limiter, exists := rl.limiters[key]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[key] = limiter
	return limiter
}

// cleanupOldEntries removes unused limiters periodically.  This is a
// naive implementation; in production you may want to track last
// access time per key.
func (rl *RateLimiter) cleanupOldEntries() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if len(rl.limiters) > 1000 {
		for k := range rl.limiters {
			delete(rl.limiters, k)
		}
	}
}

// RateLimitMiddleware returns a Gin middleware that applies rate
// limiting based on the client's IP address.  If a request exceeds
// the allowed rate a 429 response is returned and the request is
// aborted.
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanupOldEntries()
		}
	}()
	return func(c *gin.Context) {
		key := c.ClientIP()
		limiter := rl.getLimiter(key)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, domain.ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
