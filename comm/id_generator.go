package comm

import "time"

// IDGenerator can generates unique ID.
type IDGenerator interface {
	NewID() int64
}

var defaultIDGenerator = NewSnowflake(10, 12, time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC), 0)

// GenerateID generates an unique ID.
func GenerateID() int64 {
	return defaultIDGenerator.NewID()
}
