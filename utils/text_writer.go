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

// NewTextWriter creates and returns a new TextWriter.
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

// WriteString writes string to TextWriter.
func (p *TextWriter) WriteString(s string) (n int, err error) {
	n, err = p.writer.Write([]byte(s))
	if err != nil {
		p.OnError(err)
	}
	return
}

// WriteFormat writes formatted string to TextWriter.
func (p *TextWriter) WriteFormat(format string, args ...interface{}) (n int, err error) {
	return p.WriteString(fmt.Sprintf(format, args...))
}

// WriteLine writes string followed with a line separator to TextWriter.
func (p *TextWriter) WriteLine(s string) (n int, err error) {
	counter := &comm.Counter{}
	for i := 0; i < p.indentLevel; i++ {
		err = counter.Add(p.WriteString(p.LineIndent))
		if err != nil {
			return
		}
	}
	err = counter.Add(p.WriteString(s))
	if err != nil {
		return
	}
	err = counter.Add(p.WriteString(p.LineSeparator))
	if err != nil {
		return
	}
	n = counter.Count()
	return
}

// WriteLineFormat writes formatted string followed with a line separator to TextWriter.
func (p *TextWriter) WriteLineFormat(format string, args ...interface{}) (n int, err error) {
	return p.WriteLine(fmt.Sprintf(format, args...))
}

// Indent indents line when invokes WriteLine or WriteLineFormat function.
func (p *TextWriter) Indent(f func()) {
	p.indentLevel++
	f()
	p.indentLevel--
}
