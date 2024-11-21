package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRatelimiter struct {
	redisCli          *redis.Client
	tokenBucketScript *redis.Script
	capacity          int
	refreshInterval   time.Duration
}

func NewRedisRateLimiter(
	redisCli *redis.Client,
	capacity int,
	refreshInterval time.Duration,
) Ratelimiter {
	return &redisRatelimiter{
		redisCli:          redisCli,
		tokenBucketScript: redis.NewScript(tokenBucketLuaScript),
		capacity:          capacity,
		refreshInterval:   refreshInterval,
	}
}

const (
	tokenBucketLuaScript = `
		local tokenKey          = KEYS[1]
		local lastAccessTimeKey = KEYS[2]
		local capacity          = tonumber(ARGV[1])
		local intervalMS        = tonumber(ARGV[2])
		local cost              = math.max(tonumber(ARGV[3]), 0)

		redis.replicate_commands()

		local time = redis.call('TIME')
		local nowMS = math.floor((time[1] * 1000) + (time[2] / 1000))
		local initialTokens = redis.call('GET', tokenKey)
		local lastUpdateMS

		if initialTokens == false then
			initialTokens = 0
			lastUpdateMS = nowMS - intervalMS
		else
			initialTokens = tonumber(initialTokens)
			lastUpdateMS = tonumber(redis.call('GET', lastAccessTimeKey) or nowMS)
		end

		local addTokens = math.max(((nowMS - lastUpdateMS) / intervalMS) * capacity, 0)
		local grossTokens = math.min(initialTokens + addTokens, capacity)
		local remainderTokens = grossTokens - cost

		if remainderTokens < 0 then
			return {false, math.ceil(((cost - remainderTokens) / capacity) * intervalMS)}
		else
			redis.call('PSETEX', tokenKey, intervalMS, remainderTokens)
			redis.call('PSETEX', lastAccessTimeKey, intervalMS, nowMS)
			return {true, 0}
		end
	`
)

func (rl *redisRatelimiter) AllowByTokenBucket(
	ctx context.Context,
	userID string,
	cost int,
) (bool, int, error) {
	tokenKey := "ratelimit:" + userID + ":token"
	lastAccessTimeKey := "ratelimit:" + userID + ":ts"

	result, err := rl.tokenBucketScript.Run(ctx, rl.redisCli,
		[]string{tokenKey, lastAccessTimeKey},
		rl.capacity,
		rl.refreshInterval.Milliseconds(),
		cost,
		"false",
	).Result()
	if err != nil {
		return false, 0, err
	}

	data := result.([]interface{})
	allowed := data[0].(int) == 1
	retryAfter := data[1].(int)

	return allowed, retryAfter, nil
}
