package num

import (
	"fmt"
)

// simplify simplify(reduce) the fraction.
func simplify(a, b uint) (uint, uint) {
	var max, min uint
	if a > b {
		max, min = a, b
	} else {
		max, min = b, a
	}
	for remainder := max % min; remainder > 0; remainder = max % min {
		max, min = min, remainder
	}
	return a / min, b / min
}

// Fraction represents a fraction.
type Fraction struct {
	isNegative  bool
	numerator   uint
	denominator uint
}

// IsNegative returns true if the fraction is negative.
func (v Fraction) IsNegative() bool {
	return v.isNegative
}

// Numerator returns the numerator of the fraction.
func (v Fraction) Numerator() uint {
	return v.numerator
}

// Denominator returns the denominator of the fraction.
func (v Fraction) Denominator() uint {
	return v.denominator
}

func (v Fraction) String() string {
	if v.numerator == 0 {
		return "0"
	}
	if v.denominator == 1 {
		if v.isNegative {
			return fmt.Sprintf("-%d", v.numerator/v.denominator)
		}
		return fmt.Sprintf("%d", v.numerator/v.denominator)
	}
	if v.isNegative {
		return fmt.Sprintf("-%d/%d", v.numerator, v.denominator)
	}
	return fmt.Sprintf("%d/%d", v.numerator, v.denominator)
}

// Reciprocal returns reciprocal of fraction.
func (v Fraction) Reciprocal() Fraction {
	if v.numerator == 0 {
		panic(errDividedByZero)
	}
	return Fraction{
		isNegative:  v.isNegative,
		numerator:   v.denominator,
		denominator: v.numerator,
	}
}

// Float64 returns float64 representation of fraction.
func (v Fraction) Float64() (result float64) {
	result = float64(v.numerator) / float64(v.denominator)
	if v.isNegative {
		result = -result
	}
	return
}

func (v Fraction) negativeFactor() int {
	if v.isNegative {
		return -1
	}
	return 1
}

// Add adds specified fraction to current fraction.
func (v Fraction) Add(a Fraction) Fraction {
	return NewFraction(
		v.negativeFactor()*int(v.numerator*a.denominator)+a.negativeFactor()*int(a.numerator*v.denominator),
		int(v.denominator*a.denominator),
	)
}

// Sub subtracts specified fraction from current fraction.
func (v Fraction) Sub(a Fraction) Fraction {
	return NewFraction(
		v.negativeFactor()*int(v.numerator*a.denominator)-a.negativeFactor()*int(a.numerator*v.denominator),
		int(v.denominator*a.denominator),
	)
}

// Mul multiplies two fractions.
func (v Fraction) Mul(a Fraction) Fraction {
	return NewFraction(
		v.negativeFactor()*a.negativeFactor()*int(v.numerator*a.numerator),
		int(v.denominator*a.denominator),
	)
}

// Div divides current fraction by specified fraction.
func (v Fraction) Div(a Fraction) Fraction {
	return NewFraction(
		v.negativeFactor()*a.negativeFactor()*int(v.numerator*a.denominator),
		int(v.denominator*a.numerator),
	)
}

// NewFraction creates a fraction with specified numerator and denominator and returns it.
func NewFraction(numerator, denominator int) Fraction {
	if denominator == 0 {
		panic(errDividedByZero)
	}
	if numerator == 0 {
		return Fraction{
			denominator: 1,
		}
	}
	var isNegative bool
	if numerator < 0 {
		isNegative = !isNegative
		numerator = -numerator
	}
	if denominator < 0 {
		isNegative = !isNegative
		denominator = -denominator
	}
	a, b := simplify(uint(numerator), uint(denominator))
	return Fraction{
		isNegative:  isNegative,
		numerator:   a,
		denominator: b,
	}
}
