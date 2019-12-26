package num

// Abs returns absolute value of an integer.
func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// Sort returns min and max of a and b.
func Sort(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}
