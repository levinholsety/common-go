package comm

// StringArrayEqual returns true if the two input arrays are equal.
func StringArrayEqual(array1, array2 []string) bool {
	if len(array1) != len(array2) {
		return false
	}
	for i, val := range array1 {
		if val != array2[i] {
			return false
		}
	}
	return true
}

// IntArrayEqual returns true if the two input arrays are equal.
func IntArrayEqual(array1, array2 []int) bool {
	if len(array1) != len(array2) {
		return false
	}
	for i, val := range array1 {
		if val != array2[i] {
			return false
		}
	}
	return true
}
