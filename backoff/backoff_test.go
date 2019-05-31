package backoff

import (
	"fmt"
	"testing"
	"time"
)

func TestGetExponentialBackOff(t *testing.T) {
	b := NewBackoff(1, 2000, ExpoBackoffFullJitter, time.Second)
	for i := 0; i < 30; i++ {
		fmt.Println(b.BackoffDuration(i))
	}
}
