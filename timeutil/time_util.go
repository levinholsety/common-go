package timeutil

import (
	"fmt"
	"time"
)

//CurrentTimeMillis returns current time in milliseconds from 1970-01-01 00:00:00 UTC.
func CurrentTimeMillis() int64 {
	now := time.Now()
	return now.Unix()*1000 + int64(now.Nanosecond())/time.Millisecond.Nanoseconds()
}

// FormatDuration formats time duration.
func FormatDuration(value time.Duration) string {
	ms := value.Milliseconds()
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
