package ratelimiting

import (
	"sync"
	"time"
)

// Note: This file is inspired by https://github.com/uber-go/ratelimit

type LeakyBucket struct {
	perRequest time.Duration
	waitTime   time.Duration
	last       time.Time
	mutex      sync.Mutex
}

func NewLeakyBucket(rate uint64, duration time.Duration) *LeakyBucket {
	return &LeakyBucket{
		perRequest: duration / time.Duration(rate),
	}
}

func (l *LeakyBucket) Take() time.Time {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()

	if l.last.IsZero() {
		l.last = now
		return l.last
	}

	l.waitTime = l.perRequest - now.Sub(l.last)

	if l.waitTime > 0 {
		<-time.After(l.waitTime)
		l.last = now.Add(l.waitTime)
	} else {
		l.last = now
	}
	return l.last
}
