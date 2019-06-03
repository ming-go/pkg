package ratelimiting

import "time"

type RateLimiting interface {
	Take() time.Time
}
