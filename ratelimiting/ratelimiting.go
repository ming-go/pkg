package ratelimiting

type RateLimiting interface {
	Take()
}
