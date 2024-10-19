package limiter

import (
	"testing"
	"time"

	"github.com/lucasdpg/rate-limiter/internal/store"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_CheckRateLimitIP(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	redisStore := store.NewRedisStore(client)
	rl := NewRateLimiter(redisStore, 5, 0, 1*time.Minute)
	ip := "192.168.1.1"

	mr.Set(ip+":blocked", "false")
	mr.ZAdd(ip+":timestamps", float64(time.Now().UnixNano()), "some_value")

	isBlocked, err := rl.CheckRateLimitIP(ip)
	assert.NoError(t, err)
	assert.False(t, isBlocked)
}

func TestRateLimiter_CheckRateLimitToken(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	redisStore := store.NewRedisStore(client)
	rl := NewRateLimiter(redisStore, 0, 5, 1*time.Minute)
	token := "token123"

	mr.Set(token+":blocked", "false")
	mr.ZAdd(token+":timestamps", float64(time.Now().UnixNano()), "some_value")

	isBlocked, err := rl.CheckRateLimitToken(token)
	assert.NoError(t, err)
	assert.False(t, isBlocked)
}
