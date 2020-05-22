package num

import (
	"errors"
	"fmt"

	"github.com/levinholsety/common-go/comm"
)

var (
	errNotApplicable   = errors.New("not applicable")
	errNotSquareMatrix = errors.New("not square matrix")
)

// Tensor represents a tensor.
type Tensor interface {
	fmt.Stringer
	comm.Equalizer
	Duplicate() Tensor
	Add(b Tensor) Tensor
	Sub(b Tensor) Tensor
	Mul(b Tensor) Tensor
	Div(b Tensor) Tensor
	Negative() Tensor
	Reciprocal() Tensor
	Square() Tensor
	Cube() Tensor
	UO(f func(a float64) float64) Tensor
	BO(f func(a, b float64) float64, b Tensor) Tensor
	Sum() Scalar
	Mean() Scalar
	Norm() Scalar
	T() Tensor
	Dot(b Tensor) Tensor
}

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
