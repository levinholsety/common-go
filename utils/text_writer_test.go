package utils_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/utils"
)

func Test_TextWriter(t *testing.T) {
	w := utils.NewTextWriter(&bytes.Buffer{})
	w.LineIndent = "    "
	w.LineSeparator = "\n"
	w.OnError = func(err error) {}
	n, err := w.WriteLine("abc")
	assert.NoError(t, err)
	assert.IntEqual(t, 4, n)
	w.Indent(func() {
		n, err = w.WriteLine("abc")
		assert.NoError(t, err)
		assert.IntEqual(t, 8, n)
	})
}

func Test_WriteJavaClass(t *testing.T) {
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
