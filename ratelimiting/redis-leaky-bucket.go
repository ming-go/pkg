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

type redisLeakyBucketImplConfig struct {
	RPool  *redis.Pool
	Limit  int
	Period time.Duration
}

func NewRedisLeakyBucketImpl(cfg redisLeakyBucketImplConfig) *redisLeakyBucketImpl {
	return &redisLeakyBucketImpl{
		rPool:  cfg.RPool,
		limit:  cfg.Limit,
		period: cfg.Period,
	}
}

const (
	lua_script_leaky_bucket = `
		local ret = redis.call("INCRBY", KEYS[1], "1")
		if ret == 1 then
			redis.call("PEXPIRE", KEYS[1], KEYS[2])
		end
		return ret
	`
)

func (this *RedisLeakyBucketImpl) Take(token string) error {
	redisConn := this.rPool.Get()
	defer redisConn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return err
	}

	luaScript := redis.NewScript(2, lua_script_leaky_bucket)
	current, err := redis.Int(luaScript.Do(redisConn, token, strconv.Itoa(this.period/time.Millisecond)))
	if err != nil {
		return err
	}

	if current > rl.limit {
		return errors.New("API rate limit exceeded")
	}

	return nil
}

func (this *RedisLeakyBucketImpl) GetLimit() int {
	return this.limit
}

func (this *RedisLeakyBucketImpl) GetPeriod() time.Duration {
	return this.period
}
