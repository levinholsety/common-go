package dbutil

// RangeFirst returns a query range for first row.
func RangeFirst() *Range {
	return &Range{
		Offset: 0,
		Length: 1,
	}
}

// RangePage returns a query range for rows in specified page.
func RangePage(page, capacity int) *Range {
	return &Range{
		Offset: page * capacity,
		Length: capacity,
	}
}

// Range represents query range.
type Range struct {
	Offset int
	Length int
}
