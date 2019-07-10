package ratelimiting

import "time"

type Result struct {
	Remaining int
	Reset     int64
}

type RateLimitingInf interface {
	Take(string) (*Result, error)
	GetLimit() int
	GetPeriod() time.Duration
}

//type RateLimiting struct {
//	Impl RateLimitingInf
//	//Limit   int
//	//Average int // TODO
//	//Burst   int // TODO
//	//Period  time.Duration
