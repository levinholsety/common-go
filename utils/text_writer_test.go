package utils_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/utils"
)

func Test_TextWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	w := utils.NewTextWriter(buf)
	w.WriteLine("public class Test {")
	w.Indent(func() {
		w.WriteLine("public static void main(String[] args) {")
		w.Indent(func() {
			w.WriteLine("System.out.println(\"Hello World!\")")
		})
		w.WriteLine("}")
	})
	w.WriteLine("}")
	fmt.Println(buf.String())
}
