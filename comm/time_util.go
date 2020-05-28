package comm

import (
	"fmt"
	"time"
)

// TimeMillis returns unix timestamp in milliseconds of specified time.
func TimeMillis(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/time.Millisecond.Nanoseconds()
}

// CurrentTimeMillis returns current unix timestamp in milliseconds.
func CurrentTimeMillis() int64 {
	return TimeMillis(time.Now())
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
