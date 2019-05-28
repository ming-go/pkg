package dlock

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
)

var httpClientOnce sync.Once
var httpClient *http.Client

const defaultKeyPrefix = "dLock/"

type DLock struct {
	client *consulAPI.Client
	cfg    *Config
	cLock  *consulAPI.Lock
	isHeld bool
}

type Config struct {
	ConsulAddress    string
	Key              string
	TTLSecond        int
	LockWaitTime     time.Duration
	LockTryOnce      bool
	MonitorRetries   int
	UnlockAndDestroy bool
}

func NewDefaultConfig(consulAddress string, key string) *Config {
	return &Config{
		ConsulAddress:    consulAddress,
		Key:              defaultKeyPrefix + key,
		TTLSecond:        10,
		LockWaitTime:     1 * time.Second,
		LockTryOnce:      true,
		MonitorRetries:   3,
		UnlockAndDestroy: true,
	}
}

func getHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		cfg := consulAPI.DefaultConfig()
		//cfg.Transport.DisableKeepAlives = false
		//cfg.Transport.MaxIdleConns = 1024
		//cfg.Transport.MaxIdleConnsPerHost = -1
		//cfg.Transport.IdleConnTimeout = 1 * time.Hour
		consulAPI.NewClient(cfg)

		httpClient = cfg.HttpClient
	})

	return httpClient
}

func NewDLock(config *Config) (*DLock, error) {
	if strings.Trim(config.Key, " ") == "" {
		return nil, errors.New("The Key is required")
	}

	consulDefaultConfig := consulAPI.DefaultConfig()
	consulDefaultConfig.Address = config.ConsulAddress
	//consulDefaultConfig.HttpClient = getHTTPClient()

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
			LockTryOnce:  config.LockTryOnce,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DLock{
		client: client,
		cLock:  lock,
		cfg:    config,
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
	err := dl.cLock.Unlock()
	if err != nil {
		return err
	}

	if dl.cfg.UnlockAndDestroy {
		return dl.cLock.Destroy()
	}

	return nil
}

func (dl *DLock) IsHeld() bool {
	return dl.isHeld
}
