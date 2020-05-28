package comm

// FillBytes fills buffer with specified byte.
func FillBytes(buf []byte, b byte) {
	l := len(buf)
	if l == 0 {
		return
	}
	buf[0] = b
	n := 1
	for n < l {
		copy(buf[n:], buf[:n])
		n *= 2
	}
}
