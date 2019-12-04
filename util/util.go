package util

var defaultIDGenerator = NewIDGenerator(0)

// GenerateID generates a new ID with default IDGenerator.
func GenerateID() int64 {
	return defaultIDGenerator.GenerateID()
}
