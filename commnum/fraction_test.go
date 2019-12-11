package commnum_test

import (
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/commnum"
)

func Test_Fraction(t *testing.T) {
	fmt.Println(commnum.NewFraction(3, 6))
	fmt.Println(commnum.NewFraction(-3, 6))
	fmt.Println(commnum.NewFraction(-3, -6))
	fmt.Println(commnum.NewFraction(3, -6))
	fmt.Println(commnum.NewFraction(24, 170))
	fmt.Println(commnum.NewFraction(170, 24))
	fmt.Println(commnum.NewFraction(-170, 170))
	fmt.Println(commnum.NewFraction(0, -170))
	f := commnum.NewFraction(1, 2)
	f1 := commnum.NewFraction(1, 4)
	f2 := f.Add(f1)
	fmt.Println(f2)
	fmt.Println(f2.Reciprocal())
}
