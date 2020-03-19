package utils

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
func (p *TextWriter) Write(value []byte) (int, error) {
	return p.Writer.Write(value)
}

// WriteString writes string value.
func (p *TextWriter) WriteString(value string) (int, error) {
	return p.Write([]byte(value))
}

// WriteLine writes text followed with a line separator.
func (p *TextWriter) WriteLine(text string) (n int, err error) {
	var num int
	for i := 0; i < p.indent; i++ {
		if num, err = p.WriteString(p.IndentationString); err != nil {
			return
		}
		n += num
	}
	if num, err = p.WriteString(text); err != nil {
		return
	}
	n += num
	if num, err = p.WriteString(p.LineSeparator); err != nil {
		return
	}
	n += num
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
