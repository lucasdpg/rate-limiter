package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) IncrementRequestCount(ctx context.Context, key string, duration time.Duration) (int, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	_, err = r.client.Expire(ctx, key, duration).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RedisStore) GetRequestCount(ctx context.Context, key string) (int, error) {
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisStore) BlockKey(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key+":blocked", true, duration).Err()
}

func (r *RedisStore) IsBlocked(ctx context.Context, key string) (bool, error) {
	blocked, err := r.client.Get(ctx, key+":blocked").Bool()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return blocked, nil
}
