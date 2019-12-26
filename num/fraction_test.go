package num_test

import (
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/num"
)

func Test_Fraction(t *testing.T) {
	fmt.Println(num.Fraction{3, 6})
	fmt.Println(num.Fraction{-3, 6})
	fmt.Println(num.Fraction{-3, -6})
	fmt.Println(num.Fraction{3, -6})
	fmt.Println(num.Fraction{24, 170})
	fmt.Println(num.Fraction{170, 24})
	fmt.Println(num.Fraction{-170, 170})
	fmt.Println(num.Fraction{0, -170})
	f := num.Fraction{1, 2}
	f1 := num.Fraction{1, 4}
	f2 := f.Add(f1)
	fmt.Println(f2)
	fmt.Println(f2.Reciprocal())
}
