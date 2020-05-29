package comm

import (
	"fmt"
	"time"
)

// SpeedUnit represents speed unit.
type SpeedUnit int

// Speed units.
const (
	B SpeedUnit = iota + 1
	KB
	MB
	GB
)

const speedFormat = "%5.2f"

// CalculateIOSpeed calculates IO speed with byte count and elapsed time.
func CalculateIOSpeed(n int64, elapsed time.Duration) int64 {
	if elapsed == 0 {
		return 0
	}
	return n * int64(time.Second) / int64(elapsed)
}

// FormatIOSpeed format IO speed in specified unit.
func FormatIOSpeed(speed int64, unit SpeedUnit) string {
	switch unit {
	case B:
		return fmt.Sprintf(speedFormat+"  B/s", float64(speed))
	case KB:
		return fmt.Sprintf(speedFormat+" KB/s", float64(speed)/1e3)
	case MB:
		return fmt.Sprintf(speedFormat+" MB/s", float64(speed)/1e6)
	case GB:
		return fmt.Sprintf(speedFormat+" GB/s", float64(speed)/1e9)
	default:
		if speed < 1e2 {
			return FormatIOSpeed(speed, B)
		}
		if speed < 1e5 {
			return FormatIOSpeed(speed, KB)
		}
		if speed < 1e8 {
			return FormatIOSpeed(speed, MB)
		}
		return FormatIOSpeed(speed, GB)
	}
}
