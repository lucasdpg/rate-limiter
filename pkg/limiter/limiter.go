package limiter

import (
	"context"
	"time"

	"github.com/lucasdpg/rate-limiter/internal/store"
)

type RateLimiter struct {
	store            store.RateLimiterStore
	maxRequestsIP    int
	maxRequestsToken int
	blockDuration    time.Duration
}

func NewRateLimiter(s store.RateLimiterStore, maxReqIP, maxReqToken int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		store:            s,
		maxRequestsIP:    maxReqIP,
		maxRequestsToken: maxReqToken,
		blockDuration:    blockDuration,
	}
}

func (rl *RateLimiter) CheckRateLimitIP(ip string) (bool, error) {
	isBlocked, err := rl.store.IsBlocked(context.Background(), ip)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return true, nil
	}

	requestTimes, err := rl.store.GetRequestTimestamps(context.Background(), ip)
	if err != nil {
		return false, err
	}

	currentTimestamp := time.Now().UnixNano()
	oneSecondAgo := currentTimestamp - int64(time.Second)
	validRequests := filterValidRequests(requestTimes, oneSecondAgo)

	if len(validRequests) >= rl.maxRequestsIP {
		err := rl.store.BlockKey(context.Background(), ip, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = rl.store.AddRequestTimestamp(context.Background(), ip, currentTimestamp)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (rl *RateLimiter) CheckRateLimitToken(token string) (bool, error) {
	isBlocked, err := rl.store.IsBlocked(context.Background(), token)
	if err != nil {
		return false, err
	}
	if isBlocked {
		return true, nil
	}

	requestTimes, err := rl.store.GetRequestTimestamps(context.Background(), token)
	if err != nil {
		return false, err
	}

	currentTimestamp := time.Now().UnixNano()
	oneSecondAgo := currentTimestamp - int64(time.Second)
	validRequests := filterValidRequests(requestTimes, oneSecondAgo)

	if len(validRequests) >= rl.maxRequestsToken {
		err := rl.store.BlockKey(context.Background(), token, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	err = rl.store.AddRequestTimestamp(context.Background(), token, currentTimestamp)
	if err != nil {
		return false, err
	}

	return false, nil
}

func filterValidRequests(requestTimes []int64, threshold int64) []int64 {
	var validRequests []int64
	for _, t := range requestTimes {
		if t > threshold {
			validRequests = append(validRequests, t)
		}
	}
	return validRequests
}
