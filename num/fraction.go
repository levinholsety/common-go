package num

import (
	"errors"
	"fmt"
)

var errDivideByZero = errors.New("/ by zero")

// NewFraction creates a fraction with specified numerator and denominator.
func NewFraction(numerator, denominator int) (f Fraction) {
	f = Fraction{
		Numerator:   numerator,
		Denominator: denominator,
	}
	f.validate()
	f.Numerator, f.Denominator = reduce(f.Numerator, f.Denominator)
	return
}

// Fraction provides operations of fraction.
type Fraction struct {
	Numerator   int
	Denominator int
}

func (f Fraction) validate() {
	if f.Denominator == 0 {
		panic(errDivideByZero)
	}
}

func (f Fraction) String() string {
	f.validate()
	if f.Numerator == 0 {
		return "0"
	}
	f.Numerator, f.Denominator = reduce(f.Numerator, f.Denominator)
	if f.Denominator == 1 {
		return fmt.Sprintf("%d", f.Numerator)
	}
	return fmt.Sprintf("%d/%d", f.Numerator, f.Denominator)
}

// IsZero returns true if this fraction is zero.
func (f Fraction) IsZero() bool {
	f.validate()
	return f.Numerator == 0
}

// IsPositive returns true if this fraction is positive.
func (f Fraction) IsPositive() bool {
	return (f.Numerator > 0 && f.Denominator > 0) || (f.Numerator < 0 && f.Denominator < 0)
}

// IsNegative returns true if this fraction is negative.
func (f Fraction) IsNegative() bool {
	return (f.Numerator > 0 && f.Denominator < 0) || (f.Numerator < 0 && f.Denominator > 0)
}

// Add adds specified fraction to current fraction.
func (f Fraction) Add(f1 Fraction) Fraction {
	return NewFraction(f.Numerator*f1.Denominator+f1.Numerator*f.Denominator, f.Denominator*f1.Denominator)
}

// Sub subtracts specified fraction from current fraction.
func (f Fraction) Sub(f1 Fraction) Fraction {
	return NewFraction(f.Numerator*f1.Denominator-f1.Numerator*f.Denominator, f.Denominator*f1.Denominator)
}

// Mul multiplies two fractions.
func (f Fraction) Mul(f1 Fraction) Fraction {
	return NewFraction(f.Numerator*f1.Numerator, f.Denominator*f1.Denominator)
}

// Div divides current fraction by specified fraction.
func (f Fraction) Div(f1 Fraction) Fraction {
	return NewFraction(f.Numerator*f1.Denominator, f.Denominator*f1.Numerator)
}

// Reciprocal returns reciprocal of current fraction.
func (f Fraction) Reciprocal() (rcpl Fraction) {
	rcpl = Fraction{
		Numerator:   f.Denominator,
		Denominator: f.Numerator,
	}
	rcpl.validate()
	return
}

func reduce(numerator int, denominator int) (int, int) {
	if denominator == 0 {
		panic(errDivideByZero)
	}
	if numerator == 0 {
		return 0, 1
	}
	sign := 1
	if numerator < 0 {
		sign = -sign
		numerator = -numerator
	}
	if denominator < 0 {
		sign = -sign
		denominator = -denominator
	}
	if numerator == denominator {
		return 1, 1
	}
	var max, min int
	if numerator > denominator {
		max, min = numerator, denominator
	} else {
		max, min = denominator, numerator
	}
	for mod := max % min; mod > 0; mod = max % min {
		max, min = min, mod
	}
	return sign * numerator / min, denominator / min
}
