package ratelimiting

import (
	"fmt"
	"testing"
	"time"
)

func TestLeakyBucketTake(t *testing.T) {
	prev := time.Now()

	lb := NewLeakyBucket(10, 1*time.Second)

	for i := 0; i < 100; i++ {
		now := lb.Take()
		fmt.Println(now.Sub(prev))
		if i < 5 {
			<-time.After(2 * time.Second)
		}
		prev = now
	}
}
