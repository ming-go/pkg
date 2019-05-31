package mrand

import (
	"math/rand"
)

type Rand struct {
	rand *rand.Rand
}

func New(seed int64) *Rand {
	return &Rand{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (r *Rand) Uniform(min float64, max float64) float64 {
	return r.rand.Float64()*(max-min) + min
}
