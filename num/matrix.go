package num

import (
	"bytes"
	"math"
	"math/rand"
	"time"

	"github.com/levinholsety/common-go/comm"
)

// NewMatrix creates a matrix in size and returns it.
func NewMatrix(size MatrixSize) (result Matrix) {
	result = make([]Vector, size.RowCount)
	for i := 0; i < size.RowCount; i++ {
		result[i] = make([]Scalar, size.ColumnCount)
	}
	return
}

// NewMatrixWithData creates a matrix with data and returns it.
func NewMatrixWithData(data [][]float64) (result Matrix) {
	var size MatrixSize
	size.RowCount = len(data)
	if size.RowCount > 0 {
		size.ColumnCount = len(data[0])
	}
	result = NewMatrix(size)
	for i, row := range data {
		result[i] = NewVectorWithData(row...)
	}
	return
}

// Matrix represents a matrix.
type Matrix []Vector

func (p Matrix) String() string {
	buf := &bytes.Buffer{}
	size := p.Size()
	buf.WriteString("[")
	if size.RowCount > 0 {
		buf.WriteString(p[0].String())
		for i := 1; i < size.RowCount; i++ {
			buf.WriteString("\n")
			buf.WriteString(p[i].String())
		}
	}
	buf.WriteString("]")
	return buf.String()
}

// Size returns the size of self.
func (p Matrix) Size() (size MatrixSize) {
	size.RowCount = len(p)
	if size.RowCount > 0 {
		size.ColumnCount = len(p[0])
	}
	return
}

// ForEach traverses all the elements of self.
func (p Matrix) ForEach(f func(i, j int, v Scalar)) {
	for i, row := range p {
		for j, v := range row {
			f(i, j, v)
		}
	}
}

// Init initialize each elements with the specified function.
func (p Matrix) Init(f func() float64) Matrix {
	p.ForEach(func(i, j int, v Scalar) { p[i][j] = Scalar(f()) })
	return p
}

// InitRandN initialize each elements with random normally distributed float64 number.
func (p Matrix) InitRandN() Matrix {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return p.Init(func() float64 { return r.NormFloat64() })
}

// Duplicate duplicates self and returns the duplication.
func (p Matrix) Duplicate() Tensor {
	result := NewMatrix(p.Size())
	p.ForEach(func(i, j int, v Scalar) { result[i][j] = v })
	return result
}

// Equal returns true if self and the argument are equal.
func (p Matrix) Equal(a comm.Equalizer) bool {
	matrix, ok := a.(Matrix)
	if !ok {
		return false
	}
	size := p.Size()
	if size != matrix.Size() {
		return false
	}
	for i := 0; i < size.RowCount; i++ {
		for j := 0; j < size.ColumnCount; j++ {
			if !p[i][j].Equal(matrix[i][j]) {
				return false
			}
		}
	}
	return true
}

// UO executes unary operation on self with specified unary operation function.
func (p Matrix) UO(f func(a float64) float64) Tensor {
	result := NewMatrix(p.Size())
	p.ForEach(func(i, j int, v Scalar) { result[i][j] = v.UO(f).(Scalar) })
	return result
}

// BO executes binary operation on self and another tensor with specified binary operation function.
func (p Matrix) BO(f func(a, b float64) float64, tensorA Tensor) Tensor {
	ff := func(a, b Scalar) Scalar {
		return Scalar(f(float64(a), float64(b)))
	}
	if a, ok := tensorA.(Scalar); ok {
		result := NewMatrix(p.Size())
		p.ForEach(func(i, j int, v Scalar) { result[i][j] = ff(v, a) })
		return result
	}
	if a, ok := tensorA.(Vector); ok {
		size := p.Size()
		if size.ColumnCount == len(a) {
			result := NewMatrix(size)
			p.ForEach(func(i, j int, v Scalar) { result[i][j] = ff(v, a[j]) })
			return result
		}
	}
	if a, ok := tensorA.(Matrix); ok {
		size := p.Size()
		if size == a.Size() {
			result := NewMatrix(size)
			p.ForEach(func(i, j int, v Scalar) { result[i][j] = ff(v, a[i][j]) })
			return result
		}
	}
	panic(errNotApplicable)
}

// Add adds another tensor to self and returns the result.
func (p Matrix) Add(a Tensor) Tensor {
	return p.BO(add, a)
}

// Sub subtracts another tensor from self and returns the result.
func (p Matrix) Sub(a Tensor) Tensor {
	return p.BO(sub, a)
}

// Mul multiplies self by another tensor and returns the result.
func (p Matrix) Mul(a Tensor) Tensor {
	return p.BO(mul, a)
}

// Div divides self by another tensor and returns the result.
func (p Matrix) Div(a Tensor) Tensor {
	return p.BO(div, a)
}

// Negative returns the negative value of self.
func (p Matrix) Negative() Tensor {
	return p.UO(negative)
}

// Reciprocal returns the reciprocal of self.
func (p Matrix) Reciprocal() Tensor {
	return p.UO(reciprocal)
}

// Square returns the square of self.
func (p Matrix) Square() Tensor {
	return p.UO(square)
}

// Cube returns the cube of self.
func (p Matrix) Cube() Tensor {
	return p.UO(cube)
}

// Sum returns the sum of self.
func (p Matrix) Sum() (result Scalar) {
	p.ForEach(func(i, j int, v Scalar) { result += v })
	return
}

// Mean returns the mean value of self.
func (p Matrix) Mean() Scalar {
	size := p.Size()
	return p.Sum() / Scalar(size.RowCount*size.ColumnCount)
}

// Norm returns the norm value of self.
func (p Matrix) Norm() (result Scalar) {
	for _, row := range p {
		result += row.Norm().Square().(Scalar)
	}
	result = result.UO(math.Sqrt).(Scalar)
	return
}

// T returns the transpose of self.
func (p Matrix) T() Tensor {
	size := p.Size()
	size = MatrixSize{size.ColumnCount, size.RowCount}
	result := NewMatrix(size)
	p.ForEach(func(i, j int, v Scalar) { result[j][i] = v })
	return result
}

// Dot returns the dot product between self and another tensor and returns the result.
func (p Matrix) Dot(tensorA Tensor) Tensor {
	if a, ok := tensorA.(Scalar); ok {
		return p.BO(mul, a)
	}
	if a, ok := tensorA.(Vector); ok {
		size := p.Size()
		if size.ColumnCount == len(a) {
			result := NewVector(size.RowCount)
			p.ForEach(func(i, j int, v Scalar) { result[i] += v * a[j] })
			return result
		}
	}
	if a, ok := tensorA.(Matrix); ok {
		sizeA := p.Size()
		sizeB := a.Size()
		if sizeA.ColumnCount == sizeB.RowCount {
			size := MatrixSize{sizeA.RowCount, sizeB.ColumnCount}
			result := NewMatrix(size)
			p.ForEach(func(i, j int, v Scalar) {
				for k := 0; k < sizeB.ColumnCount; k++ {
					result[i][k] += v * a[j][k]
				}
			})
			return result
		}
	}
	panic(errNotApplicable)
}

// Determinant returns the determinant of self.
func (p Matrix) Determinant() (result Scalar) {
	size := p.Size()
	if size.RowCount != size.ColumnCount {
		panic(errNotSquareMatrix)
	}
	return p.determinant()
}

func (p Matrix) determinant() (result Scalar) {
	n := len(p)
	if n == 1 {
		return p[0][0]
	}
	for i := 0; i < n; i += 2 {
		minor := p.minor(0, i)
		result += p[0][i] * minor.determinant()
	}
	for i := 1; i < n; i += 2 {
		minor := p.minor(0, i)
		result -= p[0][i] * minor.determinant()
	}
	return
}

func (p Matrix) minor(rowIndex, columnIndex int) (result Matrix) {
	size := p.Size()
	size = MatrixSize{size.RowCount - 1, size.ColumnCount - 1}
	result = NewMatrix(size)
	for i := 0; i < rowIndex; i++ {
		for j := 0; j < columnIndex; j++ {
			result[i][j] = p[i][j]
		}
	}
	for i := 0; i < rowIndex; i++ {
		for j := columnIndex; j < size.ColumnCount; j++ {
			result[i][j] = p[i][j+1]
		}
	}
	for i := rowIndex; i < size.RowCount; i++ {
		for j := 0; j < columnIndex; j++ {
			result[i][j] = p[i+1][j]
		}
	}
	for i := rowIndex; i < size.RowCount; i++ {
		for j := columnIndex; j < size.ColumnCount; j++ {
			result[i][j] = p[i+1][j+1]
		}
	}
	return
}

// Inv returns the inverse matrix of self.
func (p Matrix) Inv() (result Matrix) {
	size := p.Size()
	if size.RowCount != size.ColumnCount {
		panic(errNotSquareMatrix)
	}
	result = NewMatrix(size)
	p.ForEach(func(i, j int, v Scalar) {
		if (i+j)%2 == 0 {
			result[i][j] = p.minor(i, j).determinant()
		} else {
			result[i][j] = -p.minor(i, j).determinant()
		}
	})
	result = result.T().(Matrix)
	result = result.Div(Scalar(p.Determinant())).(Matrix)
	return
}
