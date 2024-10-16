package store

import (
	"time"
)

type RateLimiterStore interface {
	IncrementRequestCount(key string) (int, error)

	GetRequestCount(key string) (int, error)

	BlockKey(key string, duration time.Duration) error

	IsBlocked(key string) (bool, error)
}
