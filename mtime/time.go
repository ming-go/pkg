package mtime

import (
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

// UnixTimeToString -
func UnixTimeToString(t int64) string {
	return strconv.FormatInt(t, 10)
}

type MTime struct {
	time time.Time
}

func New(t time.Time) *MTime {
	return &MTime{time: t}
}

// LastNDay -
func (m *MTime) LastNDay(n time.Duration) time.Time {
	return m.time.Add(-1 * SecondsPerDay * n)
}

func ParseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func UnixMilli(t time.Time) int64 {
	return (t.UnixNano() / int64(time.Millisecond))
}

//LastNDayNormalizationByWeekday -
//func LastNDayNormalizationByWeekday(n int) time.Time {
//	d := time.Duration((SecondsPerDay * time.Duration(((t.Weekday() + 7) % 8))))
//	return LastNDay(n).Add(d)
//}
