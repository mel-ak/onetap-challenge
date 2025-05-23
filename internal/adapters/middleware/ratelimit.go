package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(redisClient *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := "rate_limit:" + r.RemoteAddr

			// Get current count
			count, err := redisClient.Incr(r.Context(), key).Result()
			if err != nil {
				http.Error(w, "Rate limit error", http.StatusInternalServerError)
				return
			}

			// Set expiry on first request
			if count == 1 {
				redisClient.Expire(r.Context(), key, window)
			}

			// Check if limit exceeded
			if count > int64(limit) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
