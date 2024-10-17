package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MaxRequestsPerIP    int
	MaxRequestsPerToken int
	BlockDuration       time.Duration
	RedisURL            string
	RedisTTL            time.Duration
	ServerPort          string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables.")
	}

	maxRequestsPerIP, err := strconv.Atoi(getEnv("MAX_REQUESTS_PER_IP", "10"))
	if err != nil {
		return nil, err
	}

	maxRequestsPerToken, err := strconv.Atoi(getEnv("MAX_REQUESTS_PER_TOKEN", "100"))
	if err != nil {
		return nil, err
	}

	blockDuration, err := time.ParseDuration(getEnv("BLOCK_DURATION", "5m"))
	if err != nil {
		return nil, err
	}

	redisTTL, err := time.ParseDuration(getEnv("REDIS_TTL", "1h"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		MaxRequestsPerIP:    maxRequestsPerIP,
		MaxRequestsPerToken: maxRequestsPerToken,
		BlockDuration:       blockDuration,
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisTTL:            redisTTL,
		ServerPort:          getEnv("SERVER_PORT", "8080"),
	}

	return config, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
