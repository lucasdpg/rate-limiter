package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/lucasdpg/rate-limiter/internal/store"
)

type RateLimiter struct {
	store            store.RateLimiterStore
	maxRequestsIP    int
	maxRequestsToken int
	blockDuration    time.Duration
}

func NewRateLimiter(s store.RateLimiterStore, maxReqIP, maxReqToken int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		store:            s,
		maxRequestsIP:    maxReqIP,
		maxRequestsToken: maxReqToken,
		blockDuration:    blockDuration,
	}
}

func (rl *RateLimiter) CheckRateLimitIP(ip string) (bool, error) {
	count, err := rl.store.GetRequestCount(context.Background(), ip)
	if err != nil {
		return false, err
	}

	if count >= rl.maxRequestsIP {
		return true, nil
	}

	_, err = rl.store.IncrementRequestCount(context.Background(), ip, rl.blockDuration)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (rl *RateLimiter) CheckRateLimitToken(token string) (bool, error) {
	count, err := rl.store.GetRequestCount(context.Background(), token)
	if err != nil {
		return false, err
	}

	if count >= rl.maxRequestsToken {
		return true, nil
	}

	_, err = rl.store.IncrementRequestCount(context.Background(), token, rl.blockDuration)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (rl *RateLimiter) MiddlewareRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		token := r.Header.Get("API_KEY")

		if token != "" {
			exceeded, err := rl.CheckRateLimitToken(token)
			if err != nil {
				http.Error(w, "Error checking the limit", http.StatusInternalServerError)
				return
			}
			if exceeded {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			exceeded, err := rl.CheckRateLimitIP(ip)
			if err != nil {
				http.Error(w, "Error checking the limit", http.StatusInternalServerError)
				return
			}
			if exceeded {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
