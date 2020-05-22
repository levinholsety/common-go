package num

import (
	"errors"
	"math"
)

var (
	errDividedByZero = errors.New("divided by zero")
)

// Float64Equal returns true if the difference between a and b is less than 0.000001.
func Float64Equal(a, b float64) bool {
	return a-b < 0.000001
}

// Sum returns sum of array.
func Sum(array ...float64) float64 {
	y := float64(0)
	for _, x := range array {
		y += x
	}
	return y
}

// QuadraticSum returns quadratic sum of array.
func QuadraticSum(array ...float64) float64 {
	y := float64(0)
	for _, x := range array {
		y += x * x
	}
	return y
}

// Mean returns mean of array.
func Mean(array ...float64) float64 {
	return Sum(array...) / float64(len(array))
}

// Var returns variance of array.
func Var(array ...float64) float64 {
	n := float64(len(array))
	a := Sum(array...) / n
	return QuadraticSum(array...)/n - a*a
}

// Std returns standard deviation of array.
func Std(array ...float64) float64 {
	return math.Sqrt(Var(array...))
}
