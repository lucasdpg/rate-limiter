package store

import (
	"context"
	"time"
)

type RateLimiterStore interface {
	IncrementRequestCount(ctx context.Context, key string, duration time.Duration) (int, error)

	GetRequestCount(ctx context.Context, key string) (int, error)

	BlockKey(ctx context.Context, key string, duration time.Duration) error

	IsBlocked(ctx context.Context, key string) (bool, error)
}
