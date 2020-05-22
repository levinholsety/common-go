package num_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/num"
)

func Test_Fraction(t *testing.T) {
	assert.StringEqual(t, "1/2", num.NewFraction(3, 6).String())
	assert.StringEqual(t, "-1/2", num.NewFraction(-3, 6).String())
	assert.StringEqual(t, "1/2", num.NewFraction(-3, -6).String())
	assert.StringEqual(t, "-1/2", num.NewFraction(3, -6).String())
	assert.StringEqual(t, "12/85", num.NewFraction(24, 170).String())
	assert.StringEqual(t, "85/12", num.NewFraction(170, 24).String())
	assert.StringEqual(t, "-1", num.NewFraction(-170, 170).String())
	assert.StringEqual(t, "0", num.NewFraction(0, -170).String())
	f := num.NewFraction(1, 2)
	f1 := num.NewFraction(1, 4)
	f2 := f.Add(f1)
	assert.StringEqual(t, "3/4", f2.String())
	assert.StringEqual(t, "4/3", f2.Reciprocal().String())
}
