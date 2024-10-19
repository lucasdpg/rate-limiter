package limiter

import (
	"net"
	"net/http"
)

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Unable to parse IP address", http.StatusInternalServerError)
				return
			}
			token := r.Header.Get("API_KEY")

			if token != "" {
				exceeded, err := rl.CheckRateLimitToken(token)
				if err != nil {
					http.Error(w, "Error checking the limit", http.StatusInternalServerError)
					return
				}
				if exceeded {
					http.Error(w, "Too many requests per second", http.StatusTooManyRequests)
					return
				}
			} else {
				exceeded, err := rl.CheckRateLimitIP(ip)
				if err != nil {
					http.Error(w, "Error checking the limit", http.StatusInternalServerError)
					return
				}
				if exceeded {
					http.Error(w, "Too many requests per second", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
