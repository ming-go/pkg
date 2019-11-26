package mtime

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixTimestampToTime(t *testing.T) {
	t1, err := UnixTimestampToTime("1546300800")
	assert.Nil(t, err)
	assert.Equal(t, int64(1546300800), t1.Unix())

	t2, err := UnixTimestampToTime("0")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), t2.Unix())

	t3, err := UnixTimestampToTime("error")
	assert.NotNil(t, err)
	assert.True(t, t3.IsZero())
}

func TestUnixTimestampNanoToTime(t *testing.T) {
	t1, err := UnixTimestampNanoToTime(strconv.FormatInt(int64(1546300800)*int64(time.Second), 10))
	assert.Nil(t, err)
	assert.Equal(t, int64(1546300800000000000), t1.UnixNano())

	t2, err := UnixTimestampNanoToTime(strconv.FormatInt(int64(0)*int64(time.Second), 10))
	assert.Nil(t, err)
	assert.Equal(t, int64(0), t2.UnixNano())

	t3, err := UnixTimestampNanoToTime("error")
	assert.NotNil(t, err)
	assert.True(t, t3.IsZero())
}

func TestParseRFC3339(t *testing.T) {
	t1, err := ParseRFC3339("1970-01-01T00:00:00Z")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), t1.Unix())

	t2, err := ParseRFC3339("2019-01-01T00:00:00Z")
	assert.Nil(t, err)
	assert.Equal(t, int64(1546300800), t2.Unix())

	t3, err := ParseRFC3339("2019-01-01T00:00:00+08:00")
	assert.Nil(t, err)
	assert.Equal(t, int64(1546300800+(-8*60*60)), t3.Unix())
}

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

	t2, _ := time.Parse(time.RFC3339, "1996-09-21T00:00:00+08:00")
	assert.Equal(t, "1996-04-08T00:00:00+08:00", New(t2).LastNDay(166).Format(time.RFC3339))
}

func TestLastNHour(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2019-02-11T15:42:18+08:00")
	assert.Equal(t, "2019-02-10T14:42:18+08:00", New(tt).LastNHour(25).Format(time.RFC3339))
	assert.Equal(t, "2019-02-12T16:42:18+08:00", New(tt).LastNHour(-25).Format(time.RFC3339))
}

func TestLastNMinute(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2019-02-11T15:42:18+08:00")
	assert.Equal(t, "2019-02-11T14:17:18+08:00", New(tt).LastNMinute(85).Format(time.RFC3339))
	assert.Equal(t, "2019-02-11T17:07:18+08:00", New(tt).LastNMinute(-85).Format(time.RFC3339))

}

func TestLastNSecond(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2019-02-11T15:42:18+08:00")
	assert.Equal(t, "2019-02-11T15:40:53+08:00", New(tt).LastNSecond(85).Format(time.RFC3339))
	assert.Equal(t, "2019-02-11T15:43:43+08:00", New(tt).LastNSecond(-85).Format(time.RFC3339))
}

func TestUTCOffsetByMinutes(t *testing.T) {
	t1 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "2019-01-01T00:00:00Z", New(t1).UTCOffsetByMinutes(0).Format(time.RFC3339))
	assert.Equal(t, "2019-01-01T08:00:00+08:00", New(t1).UTCOffsetByMinutes(8*60).Format(time.RFC3339))
	assert.Equal(t, "2018-12-31T20:00:00-04:00", New(t1).UTCOffsetByMinutes(-4*60).Format(time.RFC3339))
}

func TestUTCOffsetByHours(t *testing.T) {
	t1 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "2019-01-01T00:00:00Z", New(t1).UTCOffsetByHours(0).Format(time.RFC3339))
	assert.Equal(t, "2019-01-01T08:00:00+08:00", New(t1).UTCOffsetByHours(8).Format(time.RFC3339))
	assert.Equal(t, "2018-12-31T20:00:00-04:00", New(t1).UTCOffsetByHours(-4).Format(time.RFC3339))
}

func TestToRFC3339(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Taipei") // UTC+8
	t1 := time.Date(1996, time.September, 21, 0, 0, 0, 0, location)
	assert.Equal(t, "1996-09-21T00:00:00+08:00", New(t1).ToRFC3339())

	t2 := time.Date(1996, time.September, 21, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "1996-09-21T00:00:00Z", New(t2).ToRFC3339())
}

func TestGetUTCOffsetLocationName(t *testing.T) {
	assert.Equal(t, "UTC+08:00", getUTCOffsetLocationName(8*60))
	assert.Equal(t, "UTC-04:00", getUTCOffsetLocationName(-4*60))
	assert.Equal(t, "UTC+09:30", getUTCOffsetLocationName((9*60)+30))
	assert.Equal(t, "UTC-09:30", getUTCOffsetLocationName((-9*60)-30))
}

func TestUnixMilli(t *testing.T) {
	t1 := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, int64(0), New(t1).UnixMilli())

	t2 := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, int64(1546300800000), New(t2).UnixMilli())
}
