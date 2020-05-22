package num

import (
	"bytes"
	"math"
	"math/rand"
	"time"

	"github.com/levinholsety/common-go/comm"
)

// NewVector creates a vector in size and returns it.
func NewVector(size int) Vector {
	return make([]Scalar, size)
}

// NewVectorWithData creates a vector with data and returns it.
func NewVectorWithData(data ...float64) (result Vector) {
	result = NewVector(len(data))
	for i, v := range data {
		result[i] = Scalar(v)
	}
	return
}

// Vector represents a vector.
type Vector []Scalar

func (p Vector) String() string {
	buf := &bytes.Buffer{}
	size := len(p)
	buf.WriteString("[")
	if size > 0 {
		buf.WriteString(p[0].String())
		for i := 1; i < size; i++ {
			buf.WriteString(" ")
			buf.WriteString(p[i].String())
		}
	}
	buf.WriteString("]")
	return buf.String()
}

// ForEach traverses all the elements of self.
func (p Vector) ForEach(f func(i int, v Scalar)) {
	for i, v := range p {
		f(i, v)
	}
}

// Init initialize each elements with the specified function.
func (p Vector) Init(f func() float64) Vector {
	p.ForEach(func(i int, v Scalar) { p[i] = Scalar(f()) })
	return p
}

// InitRandN initialize each elements with random normally distributed float64 number.
func (p Vector) InitRandN() Vector {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return p.Init(func() float64 { return r.NormFloat64() })
}

// Equal returns true if self and the argument are equal.
func (p Vector) Equal(a comm.Equalizer) bool {
	vector, ok := a.(Vector)
	if !ok {
		return false
	}
	length := len(p)
	if length != len(vector) {
		return false
	}
	for i := 0; i < length; i++ {
		if !p[i].Equal(vector[i]) {
			return false
		}
	}
	return true
}

// Duplicate duplicates self and returns the duplication.
func (p Vector) Duplicate() Tensor {
	result := NewVector(len(p))
	p.ForEach(func(i int, v Scalar) { result[i] = v })
	return result
}

// UO executes unary operation on self with specified unary operation function.
func (p Vector) UO(f func(a float64) float64) Tensor {
	result := NewVector(len(p))
	p.ForEach(func(i int, v Scalar) { result[i] = v.UO(f).(Scalar) })
	return result
}

// BO executes binary operation on self and another tensor with specified binary operation function.
func (p Vector) BO(f func(a, b float64) float64, tensorA Tensor) Tensor {
	ff := func(a, b Scalar) Scalar {
		return Scalar(f(float64(a), float64(b)))
	}
	if a, ok := tensorA.(Scalar); ok {
		result := NewVector(len(p))
		p.ForEach(func(i int, v Scalar) { result[i] = ff(v, a) })
		return result
	}
	if a, ok := tensorA.(Vector); ok {
		size := len(p)
		if size == len(a) {
			result := NewVector(size)
			p.ForEach(func(i int, v Scalar) { result[i] = ff(v, a[i]) })
			return result
		}
	}
	if a, ok := tensorA.(Matrix); ok {
		size := a.Size()
		if len(p) == size.ColumnCount {
			result := NewMatrix(size)
			a.ForEach(func(i, j int, v Scalar) { result[i][j] = ff(p[j], v) })
			return result
		}
	}
	panic(errNotApplicable)
}

// Add adds another tensor to self and returns the result.
func (p Vector) Add(a Tensor) Tensor {
	return p.BO(add, a)
}

// Sub subtracts another tensor from self and returns the result.
func (p Vector) Sub(a Tensor) Tensor {
	return p.BO(sub, a)
}

// Mul multiplies self by another tensor and returns the result.
func (p Vector) Mul(a Tensor) Tensor {
	return p.BO(mul, a)
}

// Div divides self by another tensor and returns the result.
func (p Vector) Div(a Tensor) Tensor {
	return p.BO(div, a)
}

// Negative returns the negative value of self.
func (p Vector) Negative() Tensor {
	return p.UO(negative)
}

// Reciprocal returns the reciprocal of self.
func (p Vector) Reciprocal() Tensor {
	return p.UO(reciprocal)
}

// Square returns the square of self.
func (p Vector) Square() Tensor {
	return p.UO(square)
}

// Cube returns the cube of self.
func (p Vector) Cube() Tensor {
	return p.UO(cube)
}

// Sum returns the sum of self.
func (p Vector) Sum() (result Scalar) {
	p.ForEach(func(i int, v Scalar) { result += v })
	return
}

// Mean returns the mean value of self.
func (p Vector) Mean() Scalar {
	return p.Sum() / Scalar(len(p))
}

// Norm returns the norm value of self.
func (p Vector) Norm() (result Scalar) {
	for _, v := range p {
		result += v.Square().(Scalar)
	}
	result = result.UO(math.Sqrt).(Scalar)
	return
}

// T returns the transpose of self.
func (p Vector) T() Tensor {
	return p.Duplicate()
}

// Dot returns the dot product between self and another tensor and returns the result.
func (p Vector) Dot(tensorA Tensor) Tensor {
	if a, ok := tensorA.(Scalar); ok {
		return p.BO(mul, a)
	}
	if a, ok := tensorA.(Vector); ok {
		if len(p) == len(a) {
			var result Scalar
			p.ForEach(func(i int, v Scalar) { result += v * a[i] })
			return result
		}
	}
	if a, ok := tensorA.(Matrix); ok {
		size := a.Size()
		if len(p) == size.RowCount {
			result := NewVector(size.ColumnCount)
			a.ForEach(func(i, j int, v Scalar) { result[j] += p[i] * v })
			return result
		}
	}
	panic(errNotApplicable)
}

// Cross returns the cross product between self and another vector and returns the result.
func (p Vector) Cross(a Vector) Tensor {
	pSize := len(p)
	aSize := len(a)
	if pSize == 2 && aSize == 2 {
		return p[0]*a[1] - p[1]*a[0]
	}
	if pSize == 2 && aSize == 3 {
		return Vector{p[1] * a[2], -p[0] * a[2], p[0]*a[1] - p[1]*a[0]}
	}
	if pSize == 3 && aSize == 2 {
		return Vector{-p[2] * a[1], p[2] * a[0], p[0]*a[1] - p[1]*a[0]}
	}
	if pSize == 3 && aSize == 3 {
		return Vector{p[1]*a[2] - p[2]*a[1], p[2]*a[0] - p[0]*a[2], p[0]*a[1] - p[1]*a[0]}
	}
	panic(errNotApplicable)
}
