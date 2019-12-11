package commio

import (
	"io"

	"github.com/levinholsety/common-go/comm"
)

// TextWriter provides functions to write text.
type TextWriter struct {
	Writer            io.Writer
	LineSeparator     string
	IndentationString string
	indent            int
}

// NewTextWriter creates a new TextWriter.
func NewTextWriter(w io.Writer) *TextWriter {
	return &TextWriter{
		Writer:            w,
		LineSeparator:     comm.LineSeparator,
		IndentationString: "    ",
	}
}

// Write writes byte array.
func (p *TextWriter) Write(value []byte) int {
	n, err := p.Writer.Write(value)
	if err != nil {
		panic(err)
	}
	return n
}

// WriteString writes string value.
func (p *TextWriter) WriteString(value string) int {
	return p.Write([]byte(value))
}

// WriteLine writes text followed with a line separator.
func (p *TextWriter) WriteLine(text string) (n int) {
	for i := 0; i < p.indent; i++ {
		n += p.WriteString(p.IndentationString)
	}
	n += p.WriteString(text)
	n += p.WriteString(p.LineSeparator)
	return
}

// IncreaseIndent increases indent.
func (p *TextWriter) IncreaseIndent() {
	p.indent++
}

// DecreaseIndent decreases indent.
func (p *TextWriter) DecreaseIndent() {
	p.indent--
}
