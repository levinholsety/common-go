package utils

import (
	"fmt"
	"io"

	"github.com/levinholsety/common-go/comm"
)

// TextWriter provides functions to write text.
type TextWriter struct {
	writer        io.Writer
	LineIndent    string
	LineSeparator string
	OnError       func(err error)
	indentLevel   int
}

// NewTextWriter creates a new TextWriter.
func NewTextWriter(w io.Writer) *TextWriter {
	return &TextWriter{
		writer:        w,
		LineIndent:    "    ",
		LineSeparator: comm.LineSeparator(),
		OnError: func(err error) {
			panic(err)
		},
		indentLevel: 0,
	}
}

// WriteString writes string value.
func (p *TextWriter) WriteString(format string, args ...interface{}) (n int, err error) {
	text := fmt.Sprintf(format, args...)
	n, err = p.writer.Write([]byte(text))
	if err != nil {
		p.OnError(err)
	}
	return
}

// WriteLine writes text followed with a line separator.
func (p *TextWriter) WriteLine(format string, args ...interface{}) (n int, err error) {
	var count int
	for i := 0; i < p.indentLevel; i++ {
		count, err = p.WriteString(p.LineIndent)
		if err != nil {
			return
		}
		n += count
	}
	count, err = p.WriteString(format, args...)
	if err != nil {
		return
	}
	n += count
	count, err = p.WriteString(p.LineSeparator)
	if err != nil {
		return
	}
	n += count
	return
}

// Indent indents line when invokes WriteLine function.
func (p *TextWriter) Indent(f func()) {
	p.indentLevel++
	f()
	p.indentLevel--
}
