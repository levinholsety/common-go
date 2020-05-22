package num_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/num"
)

var (
	aM = num.Matrix{
		{1, 2, 3},
		{4, 5, 6},
	}
	bM = num.Matrix{
		{2, 4, 8},
		{16, 32, 64},
	}
)

func Test_Matrix_Add(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{3, 6, 11},
		{20, 37, 70},
	}, aM.Add(bM))
}

func Test_Matrix_Subtract(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{-1, -2, -5},
		{-12, -27, -58},
	}, aM.Sub(bM))
}

func Test_Matrix_Multiply(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 8, 24},
		{64, 160, 384},
	}, aM.Mul(bM))
}

func Test_Matrix_Divide(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{0.5, 0.5, 0.375},
		{0.25, 0.15625, 0.09375},
	}, aM.Div(bM))
}

func Test_Matrix_AddScalar(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{3, 4, 5},
		{6, 7, 8},
	}, aM.Add(b))
}

func Test_Matrix_SubtractScalar(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{-1, 0, 1},
		{2, 3, 4},
	}, aM.Sub(b))
}

func Test_Matrix_MultiplyScalar(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 4, 6},
		{8, 10, 12},
	}, aM.Mul(b))
}

func Test_Matrix_DivideScalar(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{0.5, 1, 1.5},
		{2, 2.5, 3},
	}, aM.Div(b))
}

func Test_Matrix_AddVector(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{3, 6, 11},
		{6, 9, 14},
	}, aM.Add(bV))
}

func Test_Matrix_SubtractVector(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{-1, -2, -5},
		{2, 1, -2},
	}, aM.Sub(bV))
}

func Test_Matrix_MultiplyVector(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 8, 24},
		{8, 20, 48},
	}, aM.Mul(bV))
}

func Test_Matrix_DivideVector(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{0.5, 0.5, 0.375},
		{2, 1.25, 0.75},
	}, aM.Div(bV))
}

func Test_Matrix_Sum(t *testing.T) {
	assert.Equal(t, num.Scalar(21), aM.Sum())
}

func Test_Matrix_Mean(t *testing.T) {
	assert.Equal(t, num.Scalar(3.5), aM.Mean())
}

func Test_Matrix_Transpose(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{1, 4},
		{2, 5},
		{3, 6},
	}, aM.T())
}

func Test_Matrix_Norm(t *testing.T) {
	assert.Equal(t, num.Scalar(9.539392014169456), aM.Norm())
}

func Test_Matrix_DotProduct(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{34, 272},
		{76, 608},
	}, aM.Dot(bM.T()))
}

func Test_Matrix_DotProductScalar(t *testing.T) {
	assert.Equal(t, num.Matrix{
		{2, 4, 6},
		{8, 10, 12},
	}, aM.Dot(b))
}

func Test_Matrix_DotProductVector(t *testing.T) {
	assert.Equal(t, num.Vector{34, 76}, aM.Dot(bV))
}

func Test_Determinant(t *testing.T) {
	m1 := num.Matrix{
		{3, 8},
		{4, 6},
	}
	assert.Equal(t, num.Scalar(-14), m1.Determinant())
	m2 := num.Matrix{
		{6, 1, 1},
		{4, -2, 5},
		{2, 8, 7},
	}
	assert.Equal(t, num.Scalar(-306), m2.Determinant())
}

func Test_Inverse(t *testing.T) {
	m1 := num.Matrix{
		{3, 0, 2},
		{2, 0, -2},
		{0, 1, 1},
	}
	expected := num.Matrix{
		{0.2, 0.2, 0},
		{-0.2, 0.3, 1},
		{0.2, -0.3, 0},
	}
	assert.Equal(t, expected, m1.Inv())
}

func Test_BusAndTrain(t *testing.T) {
	price := num.Matrix{
		{3, 3.5},
		{3.2, 3.6},
	}
	total := num.Matrix{
		{118.4, 135.2},
	}
	expected := num.Matrix{
		{16, 22},
	}
	assert.Equal(t, expected, total.Dot(price.Inv()))
}
