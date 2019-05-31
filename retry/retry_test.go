package retry

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	errors "golang.org/x/xerrors"
)

func TestRetry(t *testing.T) {

	rs := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(rs)

	rI := &Retry{
		Attempts:    10,
		Backoff:     NewDefaultFixedBackoff(),
		LastErrOnly: false,
	}

	err := rI.Do(func(i int) (bool, error) {
		fmt.Println("Hello, world!")

		if rand.Intn(2) == 0 {
			return false, errors.New(strconv.Itoa(i) + "random error")
		} else {
			return false, nil
		}
	})

	fmt.Printf("%+v\n", err)
}
