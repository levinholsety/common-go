package comm

// Counter provides simple methods for counting.
type Counter struct {
	n     int
	count int
}

// Count returns current count.
func (p *Counter) Count() int {
	return p.count
}

// N returns the latest number added to counter.
func (p *Counter) N() int {
	return p.n
}

// Add adds number to counter and returns the input error.
func (p *Counter) Add(n int, err error) error {
	p.n = n
	p.count += n
	return err
}

// Counter64 provides simple methods for counting.
type Counter64 struct {
	n     int64
	count int64
}

// Count returns current count.
func (p *Counter64) Count() int64 {
	return p.count
}

// N returns the latest number added to counter.
func (p *Counter64) N() int64 {
	return p.n
}

// Add adds number to counter and returns the input error.
func (p *Counter64) Add(n int64, err error) error {
	p.n = n
	p.count += n
	return err
}
