package mtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnixTimeToString(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2019-02-11T15:42:18+08:00")
	assert.Equal(t, "1549870938", UnixTimeToString(tt.Unix()))
	assert.Equal(t, "1549870938000000000", UnixTimeToString(tt.UnixNano()))
}

func TestLastNDay(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2019-02-11T15:42:18+08:00")
	assert.Equal(t, "2019-02-10T15:42:18+08:00", New(tt).LastNDay(1).Format(time.RFC3339))
	assert.Equal(t, "2019-02-04T15:42:18+08:00", New(tt).LastNDay(7).Format(time.RFC3339))
	assert.Equal(t, "2019-01-12T15:42:18+08:00", New(tt).LastNDay(30).Format(time.RFC3339))
	assert.Equal(t, "2018-02-11T15:42:18+08:00", New(tt).LastNDay(365).Format(time.RFC3339))
}
