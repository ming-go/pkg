package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowFlake(t *testing.T) {
}

func TestSnowFlakeRetrievalWorkerId(t *testing.T) {
	s, _ := New(27, 27)
	id, _ := s.NextId()
	assert.Equal(t, int64(27), s.RetrievalWorkerId(id))
}

func TestSnowFlakeRetrievalDatacenterId(t *testing.T) {
	s, _ := New(27, 27)
	id, _ := s.NextId()
	assert.Equal(t, int64(27), s.RetrievalDatacenterId(id))
}
