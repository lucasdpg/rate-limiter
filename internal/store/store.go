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
	SetRequestTimestamp(ctx context.Context, key string) error
	GetRequestTimestamp(ctx context.Context, key string) (int64, error)
	GetRequestTimestamps(ctx context.Context, key string) ([]int64, error)
	AddRequestTimestamp(ctx context.Context, key string, timestamp int64) error
}
