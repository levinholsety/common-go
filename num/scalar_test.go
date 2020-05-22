package num_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/num"
)

var (
	a = num.Scalar(1)
	b = num.Scalar(2)
)

func Test_Scalar_Add(t *testing.T) {
	assert.Equal(t, num.Scalar(3), a.Add(b))
}

func Test_Scalar_Subtract(t *testing.T) {
	assert.Equal(t, num.Scalar(-1), a.Sub(b))
}

func Test_Scalar_Multiply(t *testing.T) {
	assert.Equal(t, num.Scalar(2), a.Mul(b))
}

func Test_Scalar_Divide(t *testing.T) {
	assert.Equal(t, num.Scalar(0.5), a.Div(b))
}

func Test_Scalar_AddVector(t *testing.T) {
	assert.Equal(t, num.Vector{3, 5, 9}, a.Add(bV))
}

func Test_Scalar_SubtractVector(t *testing.T) {
	assert.Equal(t, num.Vector{-1, -3, -7}, a.Sub(bV))
}

func Test_Scalar_MultiplyVector(t *testing.T) {
	assert.Equal(t, num.Vector{2, 4, 8}, a.Mul(bV))
}

func Test_Scalar_DivideVector(t *testing.T) {
	assert.Equal(t, num.Vector{0.5, 0.25, 0.125}, a.Div(bV))
}

func Test_Scalar_AddMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{3, 5, 9},
		{17, 33, 65},
	}, a.Add(bM))
}

func Test_Scalar_SubtractMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{-1, -3, -7},
		{-15, -31, -63},
	}, a.Sub(bM))
}

func Test_Scalar_MultiplyMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 4, 8},
		{16, 32, 64},
	}, a.Mul(bM))
}

func Test_Scalar_DivideMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{0.5, 0.25, 0.125},
		{0.0625, 0.03125, 0.015625},
	}, a.Div(bM))
}

func Test_Scalar_Sum(t *testing.T) {
	assert.Equal(t, num.Scalar(1), a.Sum())
}

func Test_Scalar_Mean(t *testing.T) {
	assert.Equal(t, num.Scalar(1), a.Mean())
}

func Test_Scalar_Transpose(t *testing.T) {
	assert.Equal(t, num.Scalar(1), a.T())
}

func Test_Scalar_DotProduct(t *testing.T) {
	assert.Equal(t, num.Scalar(2), a.Dot(b))
}

func Test_Scalar_DotProductVector(t *testing.T) {
	assert.Equal(t, num.Vector{2, 4, 8}, a.Dot(bV))
}

func Test_Scalar_DotProductMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 4, 8},
		{16, 32, 64},
	}, a.Dot(bM))
}
