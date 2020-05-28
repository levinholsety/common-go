package comm

import "strconv"

// ParseInt parses input string to int. If the length of input string is zero or any error occurs during parsing, it returns default value.
func ParseInt(s string, defaultValue int) int {
	if len(s) == 0 {
		return defaultValue
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return result
}
