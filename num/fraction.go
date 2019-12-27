package num

import (
	"fmt"
)

// Fraction represents a fraction.
type Fraction struct {
	Numerator   int
	Denominator int
}

// Reduce returns reduction of fraction.
func (f Fraction) Reduce() Fraction {
	if f.Numerator == 0 {
		return Fraction{0, f.Denominator / f.Denominator}
	}
	min, max := Sort(Abs(f.Numerator), Abs(f.Denominator))
	for remainder := max % min; remainder > 0; remainder = max % min {
		max, min = min, remainder
	}
	if f.Denominator < 0 {
		min = -min
	}
	return Fraction{f.Numerator / min, f.Denominator / min}
}

// Reciprocal returns reciprocal of fraction.
func (f Fraction) Reciprocal() Fraction {
	return Fraction{f.Denominator, f.Numerator}
}

func (f Fraction) String() string {
	if f.Numerator == 0 {
		return "0"
	}
	f = f.Reduce()
	if f.Denominator == 1 {
		return fmt.Sprintf("%d", f.Numerator/f.Denominator)
	}
	return fmt.Sprintf("%d/%d", f.Numerator, f.Denominator)
}

// Float64 returns float64 representation of fraction.
func (f Fraction) Float64() float64 {
	return float64(f.Numerator) / float64(f.Denominator)
}

// Add adds specified fraction to current fraction.
func (f Fraction) Add(f1 Fraction) Fraction {
	return Fraction{f.Numerator*f1.Denominator + f1.Numerator*f.Denominator, f.Denominator * f1.Denominator}.Reduce()
}

// Subtract subtracts specified fraction from current fraction.
func (f Fraction) Subtract(f1 Fraction) Fraction {
	return Fraction{f.Numerator*f1.Denominator - f1.Numerator*f.Denominator, f.Denominator * f1.Denominator}.Reduce()
}

// Multiply multiplies two fractions.
func (f Fraction) Multiply(f1 Fraction) Fraction {
	return Fraction{f.Numerator * f1.Numerator, f.Denominator * f1.Denominator}.Reduce()
}

// Divide divides current fraction by specified fraction.
func (f Fraction) Divide(f1 Fraction) Fraction {
	return Fraction{f.Numerator * f1.Denominator, f.Denominator * f1.Numerator}.Reduce()
}
