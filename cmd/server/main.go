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
	// Inicializando o cliente Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Criando uma instância do RedisStore
	redisStore := store.NewRedisStore(rdb)

	// Criando uma instância do Rate Limiter
	rl := limiter.NewRateLimiter(redisStore, 10, 20, time.Minute*1)

	// Inicializando o Chi Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Aplicando o middleware de Rate Limiting
	r.Use(limiter.RateLimitMiddleware(rl))

	// Servindo arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Servindo a página HTML
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Iniciando o servidor HTTP
	http.ListenAndServe(":8080", r)
}
