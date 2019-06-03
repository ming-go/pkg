package ratelimiting

import (
	"sync"
	"time"
)

// Note: This file is inspired by https://github.com/uber-go/ratelimit

type LeakyBucket struct {
	perRequst time.Duration
	waitTime  time.Duration
	last      time.Time
	mutex     sync.Mutex
}

func New(rate uint64, duration time.Duration) *LeakyBucket {
	return &LeakyBucket{
		perRequst: duration / time.Duration(rate),
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

	l.waitTime = l.perRequst - now.Sub(l.last)

	<-time.After(l.waitTime)

	if l.waitTime > 0 {
		<-time.After(l.waitTime)
		l.last = now.Add(l.waitTime)
	} else {
		l.last = now
	}
	return l.last
}
