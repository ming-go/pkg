package retry

import (
	"fmt"
	"time"
)

type RetryableFunc func(int) (bool, error)

type Retry struct {
	Attempts    int
	Delay       time.Duration
	LastErrOnly bool
}

func (retry *Retry) Do(retryableFunc RetryableFunc) (lastErr error) {
	n := 0

	for n < retry.Attempts {
		stop, err := retryableFunc(n + 1)
		if stop {
			return err
		}

		if err != nil {
			if retry.LastErrOnly {
				lastErr = err
			} else {
				if lastErr == nil {
					lastErr = fmt.Errorf("retry.Attempts: %d, Error: %v", n+1, err)
				} else {
					lastErr = fmt.Errorf("%v\nretry.Attempts: %d, Error: %v", lastErr, n+1, err)
				}
			}
		}

		n++

		<-time.After(retry.Delay)
	}

	return
}
