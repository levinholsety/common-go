package num

import (
	"errors"
	"fmt"

	"github.com/levinholsety/common-go/comm"
)

// Tensor represents a tensor.
type Tensor interface {
	fmt.Stringer
	comm.Equalizer
	Duplicate() Tensor
	Add(a Tensor) Tensor
	Sub(a Tensor) Tensor
	Mul(a Tensor) Tensor
	Div(a Tensor) Tensor
	Negative() Tensor
	Reciprocal() Tensor
	Square() Tensor
	Cube() Tensor
	Pow(a Tensor) Tensor
	Sqrt() Tensor
	Cbrt() Tensor
	Abs() Tensor
	Exp() Tensor
	Log() Tensor
	Sum() Scalar
	Mean() Scalar
	Norm() Scalar
	T() Tensor
	Dot(a Tensor) Tensor
}

var (
	errNotApplicable   = errors.New("not applicable")
	errNotSquareMatrix = errors.New("not square matrix")
)

func add(a, b float64) float64 {
	return a + b
}

func sub(a, b float64) float64 {
	return a - b
}

func mul(a, b float64) float64 {
	return a * b
}

func div(a, b float64) float64 {
	return a / b
}

func negative(a float64) float64 {
	return -a
}

func reciprocal(a float64) float64 {
	return 1 / a
}

func square(a float64) float64 {
	return a * a
}

func cube(a float64) float64 {
	return a * a * a
}
