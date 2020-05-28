package comm

import (
	"fmt"
	"time"
)

// TimeMilli returns unix timestamp of specified time in millisecond.
func TimeMilli(t time.Time) int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond()/1e6)
}

// UnixMilli returns current unix timestamp in millisecond.
func UnixMilli() int64 {
	return TimeMilli(time.Now())
}

func div(a, b int64) (int64, int64) {
	return a / b, a % b
}

const (
	timeDurationFormat         = "%02d:%02d:%02d.%03d"
	timeDurationWithDaysFormat = "%dd " + timeDurationFormat
)

// FormatTimeDuration formats time duration.
func FormatTimeDuration(value time.Duration) string {
	ms := value.Nanoseconds() / int64(time.Millisecond)
	s, ms := div(ms, 1000)
	m, s := div(s, 60)
	h, m := div(m, 60)
	d, h := div(h, 24)
	if d > 0 {
		return fmt.Sprintf(timeDurationWithDaysFormat, d, h, m, s, ms)
	}
	return fmt.Sprintf(timeDurationFormat, h, m, s, ms)
}
