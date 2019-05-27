package dlock

import (
	"errors"
	"strconv"
	"strings"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
)

const defaultKeyPrefix = "dLock/"

type DLock struct {
	client *consulAPI.Client
	//key    string
	//ttl    time.Duration
	cLock  *consulAPI.Lock
	isHeld bool
}

type Config struct {
	ConsulAddress string
	Key           string
	TTLSecond     int
	LockWaitTime  time.Duration
}

func NewDefaultConfig(consulAddress string, key string) *Config {
	return &Config{
		ConsulAddress: consulAddress,
		Key:           defaultKeyPrefix + key,
		TTLSecond:     10,
		LockWaitTime:  1 * time.Second,
	}
}

func NewDLock(config *Config) (*DLock, error) {
	if strings.Trim(config.Key, " ") == "" {
		return nil, errors.New("The Key is required")
	}

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = config.ConsulAddress

	client, err := consulAPI.NewClient(consulDefaultConfig)
	if err != nil {
		return nil, err
	}

	if config.TTLSecond < 10 {
		config.TTLSecond = 10
	}

	if config.LockWaitTime <= 10*time.Millisecond {
		config.LockWaitTime = 1 * time.Second
	}

	lock, err := client.LockOpts(
		&consulAPI.LockOptions{
			Key:          config.Key,
			SessionTTL:   strconv.Itoa(config.TTLSecond) + "s",
			LockWaitTime: config.LockWaitTime,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DLock{
		client: client,
		cLock:  lock,
		//ttl:    time.Duration(config.TTLSecond) * time.Second,
	}, nil
}

func (dl *DLock) lock() (bool, error) {
	stopCh := make(<-chan struct{})

	leaderCh, err := dl.cLock.Lock(stopCh)
	if err != nil {
		return false, err
	}

	if leaderCh == nil {
		return false, nil
	}

	dl.isHeld = true

	go func() {
		for {
			select {
			case <-leaderCh:
				dl.isHeld = false
			default:
			}

			<-time.After(1 * time.Millisecond)
		}
	}()

	return dl.isHeld, nil
}

func (dl *DLock) Lock() (bool, error) {
	return dl.lock()
}

func (dl *DLock) Unlock() error {
	return dl.cLock.Unlock()
}

func (dl *DLock) IsHeld() bool {
	return dl.isHeld
}
