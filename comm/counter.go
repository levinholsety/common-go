package comm

// Counter provides simple methods for counting.
type Counter struct {
	count int
}

// Count returns current count.
func (p *Counter) Count() int {
	return p.count
}

// Add adds number to total count and returns the input error.
func (p *Counter) Add(n int, err error) error {
	p.count += n
	return err
}
