package num

import (
	"fmt"
	"math"
	"strconv"

	"github.com/levinholsety/common-go/comm"
)

// Scalar represents a scalar.
type Scalar float64

func (v Scalar) String() string {
	f, _ := strconv.ParseFloat(fmt.Sprintf("%f", v), 64)
	return fmt.Sprintf("%g", f)
}

// Equal returns true if self and the argument are equal.
func (v Scalar) Equal(a comm.Equalizer) bool {
	scalar, ok := a.(Scalar)
	if !ok {
		return false
	}
	return v-scalar < 0.000001
}

// Duplicate duplicates self and returns the duplication.
func (v Scalar) Duplicate() Tensor {
	return v
}

// UO executes unary operation on self with specified unary operation function.
func (v Scalar) UO(f func(a float64) float64) Tensor {
	return Scalar(f(float64(v)))
}

// BO executes binary operation on self and another tensor with specified binary operation function.
func (v Scalar) BO(f func(a, b float64) float64, tensorA Tensor) Tensor {
	ff := func(a, b Scalar) Scalar {
		return Scalar(f(float64(a), float64(b)))
	}
	if a, ok := tensorA.(Scalar); ok {
		return ff(v, a)
	}
	if a, ok := tensorA.(Vector); ok {
		result := NewVector(len(a))
		a.ForEach(func(i int, s Scalar) { result[i] = ff(v, s) })
		return result
	}
	if a, ok := tensorA.(Matrix); ok {
		result := NewMatrix(a.Size())
		a.ForEach(func(i, j int, s Scalar) { result[i][j] = ff(v, s) })
		return result
	}
	panic(errNotApplicable)
}

// Add adds another tensor to self and returns the result.
func (v Scalar) Add(a Tensor) Tensor {
	return v.BO(add, a)
}

// Sub subtracts another tensor from self and returns the result.
func (v Scalar) Sub(a Tensor) Tensor {
	return v.BO(sub, a)
}

// Mul multiplies self by another tensor and returns the result.
func (v Scalar) Mul(a Tensor) Tensor {
	return v.BO(mul, a)
}

// Div divides self by another tensor and returns the result.
func (v Scalar) Div(a Tensor) Tensor {
	return v.BO(div, a)
}

// Negative returns the negative value of self.
func (v Scalar) Negative() Tensor {
	return v.UO(negative)
}

// Reciprocal returns the reciprocal of self.
func (v Scalar) Reciprocal() Tensor {
	return v.UO(reciprocal)
}

// Square returns the square of self.
func (v Scalar) Square() Tensor {
	return v.UO(square)
}

// Cube returns the cube of self.
func (v Scalar) Cube() Tensor {
	return v.UO(cube)
}

// Sum returns the sum of self.
func (v Scalar) Sum() Scalar {
	return v
}

// Mean returns the mean value of self.
func (v Scalar) Mean() Scalar {
	return v
}

// Norm returns the norm value of self.
func (v Scalar) Norm() Scalar {
	return Scalar(math.Sqrt(float64(v * v)))
}

// T returns the transpose of self.
func (v Scalar) T() Tensor {
	return v
}

// Dot returns the dot product between self and another tensor and returns the result.
func (v Scalar) Dot(a Tensor) Tensor {
	return v.BO(mul, a)
}
