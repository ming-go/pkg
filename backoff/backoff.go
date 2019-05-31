// Inspire by https://www.awsarchitectureblog.com/2015/03/backoff.html and
// https://github.com/aws-samples/aws-arch-backoff-simulator

package backoff

import (
	"math"
	rand "math/rand"

	"sync"
	"time"

	"github.com/ming-go/pkg/mrand"
)

var grand = rand.New(rand.NewSource(time.Now().UnixNano()))
var mrnd = mrand.New(time.Now().UnixNano())

type BackoffAlgorithm string

const (
	NoBackoff              = "NoBackoff"
	FixedBackoff           = "FixedBackoff"
	ExpoBackoff            = "ExpoBackoff"
	ExpoBackoffEqualJitter = "ExpoBackoffEqualJitter"
	ExpoBackoffFullJitter  = "ExpoBackoffFullJitter"
	ExpoBackoffDecorr      = "ExpoBackoffDecorr"
)

type Backoff struct {
	base      float64
	duration  time.Duration
	cap       float64
	algorithm BackoffAlgorithm
	decorr    *Decorr
}

type Decorr struct {
	sleep float64
	cap   float64
	base  float64
	mutex sync.RWMutex
}

func (d *Decorr) getSleep() float64 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.sleep
}

func (d *Decorr) setSleep(v float64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.sleep = v
}

func (d *Decorr) Decorr() float64 {
	// TODO: Can be optimized
	d.setSleep(math.Min(d.cap, mrnd.Uniform(d.base, d.getSleep()*3.0)))
	return d.getSleep()
}

func NewBackoff(factor float64, maximum float64, algorithm BackoffAlgorithm, duration time.Duration) *Backoff {
	var decorr *Decorr = nil

	switch algorithm {
	case NoBackoff, ExpoBackoff, ExpoBackoffEqualJitter, ExpoBackoffFullJitter, FixedBackoff:
	case ExpoBackoffDecorr:
		decorr = &Decorr{
			sleep: factor,
			cap:   maximum,
			base:  factor,
		}
	default:
		algorithm = NoBackoff
	}

	return &Backoff{
		base:      factor,
		cap:       maximum,
		algorithm: algorithm,
		decorr:    decorr,
		duration:  duration,
	}
}

func (b *Backoff) backoff(n int) float64 {
	switch b.algorithm {
	case NoBackoff:
		return 0.0
	case FixedBackoff:
		return b.base
	case ExpoBackoff:
		return Exponential(b.base, n, b.cap)
	case ExpoBackoffEqualJitter:
		return EqualJitter(Exponential(b.base, n, b.cap))
	case ExpoBackoffFullJitter:
		return FullJitter(Exponential(b.base, n, b.cap))
	case ExpoBackoffDecorr:
		return b.decorr.Decorr()
	default:
		return 0.0
	}
}

func (b *Backoff) Backoff(n int) float64 {
	return b.backoff(n)
}

func (b *Backoff) BackoffDuration(n int) time.Duration {
	return time.Duration(b.backoff(n)) * b.duration
}

func Exponential(base float64, n int, cap float64) float64 {
	// Will be zero if factor equals 0
	countdown := math.Min(cap, base*math.Pow(2.0, float64(n)))

	// Adjust according to maximum wait time and account for negative values.
	return math.Max(0.0, countdown)
}

func EqualJitter(expo float64) float64 {
	return (expo / 2) + mrnd.Uniform(0, (expo/2))
}

func FullJitter(expo float64) float64 {
	return mrnd.Uniform(0, expo)
}
