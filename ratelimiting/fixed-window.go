package ratelimiting

import (
	"sync"
	"time"
)

type FixedWindow struct {
	counter  uint64
	rate     uint64
	last     time.Time
	duration time.Duration
	mutex    sync.Mutex
}

func NewFixedWindow(rate uint64, duration time.Duration) *FixedWindow {
	return &FixedWindow{}
}

func (f *FixedWindow) Teke() time.Time {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	now := time.Now()

	if now.Sub(f.last) > f.duration {
		f.counter = 0
	}

	f.counter++

	if f.counter > f.rate {
		waitTime := f.duration - now.Sub(f.last)
		<-time.After(waitTime)
		f.last = now.Add(waitTime)
	} else {
		f.last = now
	}

	return f.last
}
