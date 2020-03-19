package comm_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/comm"
)

func Test_Fraction(t *testing.T) {
	assert.StringEqual(t, "1/2", comm.Fraction{3, 6}.String())
	assert.StringEqual(t, "-1/2", comm.Fraction{-3, 6}.String())
	assert.StringEqual(t, "1/2", comm.Fraction{-3, -6}.String())
	assert.StringEqual(t, "-1/2", comm.Fraction{3, -6}.String())
	assert.StringEqual(t, "12/85", comm.Fraction{24, 170}.String())
	assert.StringEqual(t, "85/12", comm.Fraction{170, 24}.String())
	assert.StringEqual(t, "-1", comm.Fraction{-170, 170}.String())
	assert.StringEqual(t, "0", comm.Fraction{0, -170}.String())
	f := comm.Fraction{1, 2}
	f1 := comm.Fraction{1, 4}
	f2 := f.Add(f1)
	assert.StringEqual(t, "3/4", f2.String())
	assert.StringEqual(t, "4/3", f2.Reciprocal().String())
}
