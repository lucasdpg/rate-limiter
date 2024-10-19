package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/lucasdpg/rate-limiter/config"
	"github.com/lucasdpg/rate-limiter/internal/store"
	"github.com/lucasdpg/rate-limiter/pkg/limiter"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config %s", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	redisStore := store.NewRedisStore(rdb)

	rl := limiter.NewRateLimiter(redisStore, cfg.MaxRequestsPerIP, cfg.MaxRequestsPerToken, cfg.BlockDuration)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(limiter.RateLimitMiddleware(rl))

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.ListenAndServe(":8080", r)
}
