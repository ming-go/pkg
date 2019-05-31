package retry

import (
	"fmt"
	"time"

	"github.com/ming-go/pkg/backoff"
)

type RetryableFunc func(int) (bool, error)

type Retry struct {
	Attempts    int
	Backoff     *backoff.Backoff
	LastErrOnly bool
}

func NewDefaultFixedBackoff() *backoff.Backoff {
	return backoff.NewBackoff(
		1,
		600,
		backoff.FixedBackoff,
		time.Second,
	)
}

func NewDefaultBackoffWithFullJitter() *backoff.Backoff {
	return backoff.NewBackoff(
		1,
		600,
		backoff.ExpoBackoffFullJitter,
		time.Second,
	)
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

		<-time.After(retry.Backoff.BackoffDuration(n))
	}

	return
}
