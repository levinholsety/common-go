package num_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/num"
)

var (
	aV = num.Vector{1, 2, 3}
	bV = num.Vector{2, 4, 8}
	cV = num.Vector{1, 2}
)

func Test_Vector_Add(t *testing.T) {
	assert.Equal(t, num.Vector{3, 6, 11}, aV.Add(bV))
}

func Test_Vector_Subtract(t *testing.T) {
	assert.Equal(t, num.Vector{-1, -2, -5}, aV.Sub(bV))
}

func Test_Vector_Multiply(t *testing.T) {
	assert.Equal(t, num.Vector{2, 8, 24}, aV.Mul(bV))
}

func Test_Vector_Divide(t *testing.T) {
	assert.Equal(t, num.Vector{0.5, 0.5, 0.375}, aV.Div(bV))
}

func Test_Vector_AddVector(t *testing.T) {
	assert.Equal(t, num.Vector{3, 4, 5}, aV.Add(b))
}

func Test_Vector_SubtractVector(t *testing.T) {
	assert.Equal(t, num.Vector{-1, 0, 1}, aV.Sub(b))
}

func Test_Vector_MultiplyVector(t *testing.T) {
	assert.Equal(t, num.Vector{2, 4, 6}, aV.Mul(b))
}

func Test_Vector_DivideVector(t *testing.T) {
	assert.Equal(t, num.Vector{0.5, 1, 1.5}, aV.Div(b))
}

func Test_Vector_AddMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{3, 6, 11},
		{17, 34, 61},
	}, aV.Add(bM))
}

func Test_Vector_SubtractMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{-1, -2, -5},
		{-15, -30, -61},
	}, aV.Sub(bM))
}

func Test_Vector_MultiplyMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 8, 24},
		{16, 64, 192},
	}, aV.Mul(bM))
}

func Test_Vector_DivideMatrix(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{0.5, 0.5, 0.375},
		{0.0625, 0.0625, 0.046875},
	}, aV.Div(bM))
}

func Test_Vector_Sum(t *testing.T) {
	assert.Equal(t, num.Scalar(6), aV.Sum())
}

func Test_Vector_Mean(t *testing.T) {
	assert.Equal(t, num.Scalar(2), aV.Mean())
}

func Test_Vector_Transpose(t *testing.T) {
	assert.Equal(t, num.Vector{1, 2, 3}, aV.T())
}

func Test_Vector_DotProduct(t *testing.T) {
	assert.Equal(t, num.Scalar(34), aV.Dot(bV))
}

func Test_Vector_DotProductScalar(t *testing.T) {
	assert.Equal(t, num.Vector{2, 4, 6}, aV.Dot(b))
}

func Test_Vector_DotProductMatrix(t *testing.T) {
	assert.Equal(t, num.Vector{34, 68, 136}, cV.Dot(bM))
}

func Test_Vector_Norm(t *testing.T) {
	assert.Equal(t, num.Scalar(3.741657), aV.Norm())
}

func Test_Vector_CrossProduct(t *testing.T) {
	a := num.Vector{1, 2}
	b := num.Vector{3, 4}
	assert.Equal(t, num.Scalar(-2), a.Cross(b))
	assert.Equal(t, num.Scalar(2), b.Cross(a))
	c := num.Vector{3, 4, 5}
	assert.Equal(t, num.Vector{10, -5, -2}, a.Cross(c))
	assert.Equal(t, num.Vector{-20, 15, 0}, c.Cross(b))
	d := num.Vector{1, 2, 4}
	assert.Equal(t, num.Vector{6, -7, 2}, c.Cross(d))
	assert.Equal(t, num.Vector{-6, 7, -2}, d.Cross(c))
}
