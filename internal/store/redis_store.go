package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore implementa a interface RateLimiterStore usando o Redis como backend
type RedisStore struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration // Tempo de vida (TTL) para as contagens de requisições
}

// NewRedisStore cria uma nova instância de RedisStore
func NewRedisStore(redisURL string, ttl time.Duration) (*RedisStore, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	ctx := context.Background()

	// Verifica a conexão com o Redis
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client: client,
		ctx:    ctx,
		ttl:    ttl,
	}, nil
}

// IncrementRequestCount incrementa o contador de requisições para a chave (IP ou Token) no Redis
func (r *RedisStore) IncrementRequestCount(key string) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Define o TTL na primeira requisição
	if count == 1 {
		err = r.client.Expire(r.ctx, key, r.ttl).Err()
		if err != nil {
			return 0, err
		}
	}

	return int(count), nil
}

// GetRequestCount retorna o número de requisições associadas a uma chave (IP ou Token)
func (r *RedisStore) GetRequestCount(key string) (int, error) {
	count, err := r.client.Get(r.ctx, key).Int()
	if err != nil {
		// Se a chave não existir, retorna 0
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// BlockKey bloqueia uma chave (IP ou Token) por um tempo definido
func (r *RedisStore) BlockKey(key string, duration time.Duration) error {
	return r.client.Set(r.ctx, key, -1, duration).Err()
}

// IsBlocked verifica se uma chave (IP ou Token) está bloqueada
func (r *RedisStore) IsBlocked(key string) (bool, error) {
	count, err := r.client.Get(r.ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// Se o valor for -1, a chave está bloqueada
	if count == -1 {
		return true, nil
	}

	return false, nil
}
