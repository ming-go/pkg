package mstring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty(" "))
	assert.True(t, IsEmpty("   "))
	assert.False(t, IsEmpty("%"))
	assert.False(t, IsEmpty(" % "))
	assert.False(t, IsEmpty("   %   "))
}

func TestIsURI(t *testing.T) {
	assert.True(t, IsURI("https://www.google.com.tw/"))
	assert.True(t, IsURI("http://www.google.com.tw/"))
	assert.True(t, IsURI("https://alone.tw/"))
	assert.True(t, IsURI("http://alone.tw/"))
	assert.False(t, IsURI("alone.tw/"))
	assert.False(t, IsURI("httpalonetw"))
}
