package num

import "fmt"

// NewMatrixSize returns a MatrixSize with specified row count and column count.
func NewMatrixSize(rowCount, columnCount int) MatrixSize {
	return MatrixSize{
		RowCount:    rowCount,
		ColumnCount: columnCount,
	}
}

// MatrixSize represents size of matrix.
type MatrixSize struct {
	RowCount    int
	ColumnCount int
}

func (v MatrixSize) String() string {
	return fmt.Sprintf("Matrix(%d,%d)", v.RowCount, v.ColumnCount)
}
