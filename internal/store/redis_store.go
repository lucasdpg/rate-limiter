package store

import (
	"context"
	"fmt"
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

func (r *RedisStore) SetRequestTimestamp(ctx context.Context, key string) error {
	timestamp := time.Now().UnixNano()
	return r.client.Set(ctx, key+":timestamp", timestamp, 0).Err()
}

func (r *RedisStore) GetRequestTimestamp(ctx context.Context, key string) (int64, error) {
	timestamp, err := r.client.Get(ctx, key+":timestamp").Int64()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return timestamp, nil
}

func (r *RedisStore) AddRequestTimestamp(ctx context.Context, key string, timestamp int64) error {
	return r.client.ZAdd(ctx, key+":timestamps", &redis.Z{
		Score:  float64(timestamp),
		Member: timestamp,
	}).Err()
}

func (r *RedisStore) GetRequestTimestamps(ctx context.Context, key string) ([]int64, error) {
	timestamps, err := r.client.ZRangeWithScores(ctx, key+":timestamps", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, ts := range timestamps {
		result = append(result, int64(ts.Score))
	}
	return result, nil
}

func (r *RedisStore) CleanupOldTimestamps(ctx context.Context, key string, threshold int64) error {
	return r.client.ZRemRangeByScore(ctx, key+":timestamps", "0", fmt.Sprintf("%d", threshold)).Err()
}
