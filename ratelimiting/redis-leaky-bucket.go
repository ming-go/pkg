package ratelimiting

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisLeakyBucketImpl struct {
	rPool  *redis.Pool
	limit  int
	period time.Duration
}

type RedisLeakyBucketImplConfig struct {
	RPool  *redis.Pool
	Limit  int
	Period time.Duration
}

func NewRedisLeakyBucketImpl(cfg *RedisLeakyBucketImplConfig) *redisLeakyBucketImpl {
	return &redisLeakyBucketImpl{
		rPool:  cfg.RPool,
		limit:  cfg.Limit,
		period: cfg.Period,
	}
}

const (
	lua_script_leaky_bucket = `
		local current = redis.call("INCRBY", KEYS[1], "1")
		if current == 1 then
			redis.call("PEXPIRE", KEYS[1], KEYS[2])
		end
		local pttl = redis.call("PTTL", KEYS[1])
		return {current, pttl}
	`
)

func (this *redisLeakyBucketImpl) Take(token string) (*Result, error) {
	redisConn := this.rPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("PING")
	if err != nil {
		return nil, err
	}

	luaScript := redis.NewScript(2, lua_script_leaky_bucket)
	ret, err := redis.Ints(luaScript.Do(redisConn, token, strconv.Itoa(int(this.period)/int(time.Millisecond))))
	if err != nil {
		return nil, err
	}

	res := &Result{
		Consumed: ret[0],
		PTTL:     ret[1],
	}

	if ret[0] > this.limit {
		return res, errors.New("API rate limit exceeded")
	}

	return res, nil
}

func (this *redisLeakyBucketImpl) GetLimit() int {
	return this.limit
}

func (this *redisLeakyBucketImpl) GetPeriod() time.Duration {
	return this.period
}
