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

// FormatDuration formats time duration.
func FormatDuration(value time.Duration) string {
	ms := value.Nanoseconds() / 1e6
	s := ms / 1000
	m := s / 60
	h := m / 60
	d := h / 24
	ms = ms % 1000
	s = s % 60
	m = m % 60
	h = h % 24
	if d > 0 {
		return fmt.Sprintf("%dd %02d:%02d:%02d.%03d", d, h, m, s, ms)
	}
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
