package ratelimiting

import "time"

type RateLimitingInf interface {
	Take(string) error
	GetLimit() int
	GetPeriod() time.Duration
}

type RateLimiting struct {
	Impl RateLimitingInf
	//Limit   int
	//Average int // TODO
	//Burst   int // TODO
	//Period  time.Duration
}

func (rl *RateLimiting) Take(string) error {
	rl.Impl.Take()
}
