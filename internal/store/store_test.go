package store

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore_IncrementRequestCount(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"
	duration := time.Minute

	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, duration).SetVal(true)

	count, err := store.IncrementRequestCount(ctx, key, duration)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_GetRequestCount(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"

	mock.ExpectGet(key).SetVal("5")

	count, err := store.GetRequestCount(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_BlockKey(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"
	duration := time.Minute

	mock.ExpectSet(key+":blocked", true, duration).SetVal("OK")

	err := store.BlockKey(ctx, key, duration)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_IsBlocked(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"

	mock.ExpectGet(key + ":blocked").SetVal("1")

	blocked, err := store.IsBlocked(ctx, key)

	assert.NoError(t, err)
	assert.True(t, blocked)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_GetRequestTimestamp(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"

	timestamp := time.Now().UnixNano()
	mock.ExpectGet(key + ":timestamp").SetVal(fmt.Sprintf("%d", timestamp))

	ts, err := store.GetRequestTimestamp(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, timestamp, ts)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_AddRequestTimestamp(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"
	timestamp := time.Now().UnixNano()

	mock.ExpectZAdd(key+":timestamps", &redis.Z{
		Score:  float64(timestamp),
		Member: timestamp,
	}).SetVal(1)

	err := store.AddRequestTimestamp(ctx, key, timestamp)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_GetRequestTimestamps(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"

	timestamp := float64(time.Now().UnixNano())
	mock.ExpectZRangeWithScores(key+":timestamps", 0, -1).SetVal([]redis.Z{
		{Score: timestamp, Member: timestamp},
	})

	timestamps, err := store.GetRequestTimestamps(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, []int64{int64(timestamp)}, timestamps)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_CleanupOldTimestamps(t *testing.T) {
	db, mock := redismock.NewClientMock()
	store := NewRedisStore(db)
	ctx := context.Background()

	key := "user:123"
	threshold := time.Now().UnixNano() - int64(time.Minute)

	mock.ExpectZRemRangeByScore(key+":timestamps", "0", strconv.FormatInt(threshold, 10)).SetVal(1)

	err := store.CleanupOldTimestamps(ctx, key, threshold)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
