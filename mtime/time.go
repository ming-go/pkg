package mtime

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	SecondsPerMinute = time.Minute
	SecondsPerHour   = time.Hour
	SecondsPerDay    = 24 * SecondsPerHour
	SecondsPerWeek   = 7 * SecondsPerDay
	SecondsPer30Days = 30 * SecondsPerDay
	SecondsPer31Days = 31 * SecondsPerDay
	DaysPer400Years  = 365*400 + 97
	DaysPer100Years  = 365*100 + 24
	DaysPer4Years    = 365*4 + 1
)

// UnixTimestampNanoToTime unix timestamp (in nanosecond) to time.Time
func UnixTimestampNanoToTime(utns string) (time.Time, error) {
	nsec, err := strconv.ParseInt(utns, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, nsec), nil
}

// UnixTimestampToTime unix timestamp (in second) to time.Time
func UnixTimestampToTime(uts string) (time.Time, error) {
	sec, err := strconv.ParseInt(uts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(sec, 0), nil
}

// UnixTimeToString Unit timestamp (in second) to string
func UnixTimeToString(t int64) string {
	return strconv.FormatInt(t, 10)
}

// MTime provides a time handler
type MTime struct {
	time time.Time
}

// New returns a new MTime
func New(t time.Time) *MTime {
	return &MTime{time: t}
}

// LastNDay get n days ago
func (m *MTime) LastNDay(n time.Duration) time.Time {
	return m.time.Add(-1 * SecondsPerDay * n)
}

// LastNHour get n hours ago
func (m *MTime) LastNHour(n int) time.Time {
	return m.time.Add(-1 * SecondsPerHour * time.Duration(n))
}

// LastNMinute get n minutes ago
func (m *MTime) LastNMinute(n int) time.Time {
	return m.time.Add(-1 * SecondsPerMinute * time.Duration(n))
}

// LastNSecond get n seconds ago
func (m *MTime) LastNSecond(n int) time.Time {
	return m.time.Add(-1 * time.Second * time.Duration(n))
}

// ParseRFC3339 Converting RFC3339 string to time.Time eg. 2019-11-25T23:00:41-04:00
func ParseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// ToRFC3339 Converting time.Time to RFC3339 string
func (m *MTime) ToRFC3339() string {
	return m.time.Format(time.RFC3339)
}

// UnixMilli time.Time to Unix Timestamp (in millisecond)
func (m *MTime) UnixMilli() int64 {
	return (m.time.UnixNano() / int64(time.Millisecond))
}

// UTCOffsetByHours Converting time.Time to UTC+-N
func (m *MTime) UTCOffsetByHours(nHours int) time.Time {
	return m.UTCOffsetByMinutes(nHours * 60)
}

func getUTCOffsetLocationName(nMinutes int) string {
	sign := "+"
	if nMinutes < 0 {
		sign = "-"
	}

	return "UTC" + sign + fmt.Sprintf("%02d:%02d", int(math.Abs(float64(nMinutes/60))), int(math.Abs(float64(nMinutes%60))))
}

// UTCOffsetByMinutes Converting time.Time to UTC+-N
func (m *MTime) UTCOffsetByMinutes(nMinutes int) time.Time {
	return m.time.In(
		time.FixedZone(
			getUTCOffsetLocationName(nMinutes),
			nMinutes*60,
		),
	)
}
