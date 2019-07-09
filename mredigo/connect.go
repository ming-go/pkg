package mredigo

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	errors "golang.org/x/xerrors"
)

var (
	poolStore sync.Map
)

type config struct {
	Host     string
	Database string
	Password string
}

func NewConfig() *config {
	return &config{}
}

func newPool(config *config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     128,
		MaxActive:   0,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			var conn redis.Conn
			var err error
			if len(strings.TrimSpace(config.Password)) == 0 {
				conn, err = redis.Dial("tcp", config.Host)
			} else {
				conn, err = redis.Dial("tcp", config.Host, redis.DialPassword(config.Password))
			}

			if err != nil {
				return nil, err
			}

			if len(strings.TrimSpace(config.Database)) == 0 {
				return conn, err
			}

			dbNum, err := strconv.ParseInt(config.Database, 10, 32)
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("SELECT", dbNum)
			if err != nil {
				conn.Close()
				return nil, err
			}

			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
			_, err = c.Do("PING")
			return err
		},
	}
}

func Connect(key string) (redis.Conn, error) {
	v, ok := poolStore.Load(key)
	if !ok {
		return nil, errors.New("key is not exists")
	}

	conn := v.(*redis.Pool).Get()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}

	return conn, err
}

func GetPool(key string) (*redis.Pool, error) {
	v, ok := poolStore.Load(key)
	if !ok {
		return nil, errors.New("key is not exists")
	}

	return v.(*redis.Pool), nil
}

func CreatePool(key string, override bool, config *config) (*redis.Pool, error) {
	if _, ok := poolStore.Load(key); !ok && !override {
		return errors.New("key is exists")
	}

	p := newPool(config)

	conn := p.Get()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, err
	}
	conn.Close()

	if p != nil {
		poolStore.Store(key, p)
	}

	return p, nil
}
