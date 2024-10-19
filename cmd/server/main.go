package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/lucasdpg/rate-limiter/internal/store"
	"github.com/lucasdpg/rate-limiter/pkg/limiter"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", //passar para ler de env?.
	})

	redisStore := store.NewRedisStore(rdb)

	rl := limiter.NewRateLimiter(redisStore, 2, 4, time.Minute*1) //passar parametro para var

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
