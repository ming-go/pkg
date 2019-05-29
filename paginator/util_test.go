package paginator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 7, maxInt(0, 7))
	assert.Equal(t, 7, maxInt(7, 0))
	assert.Equal(t, 0, maxInt(0, 0))
	assert.Equal(t, 7, maxInt(7, 7))
}

func TestIsSlice(t *testing.T) {
	assert.True(t, isSlice([]string{"one"}))
	assert.True(t, isSlice([]int{1}))
	assert.True(t, isSlice([]struct{}{}))
	assert.False(t, isSlice(1))
	assert.False(t, isSlice("one"))
	assert.False(t, isSlice(struct{}{}))
}

// func TestCopySlice(t *testing.T)
