package ratelimiting

import "time"

type RateLimitingInf interface {
	Take(string) error
	GetLimit() int
	GetPeriod() time.Duration
}

//type RateLimiting struct {
//	Impl RateLimitingInf
//	//Limit   int
//	//Average int // TODO
//	//Burst   int // TODO
//	//Period  time.Duration
