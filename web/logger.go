package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/levinholsety/common-go/comm"
)

// NewLogger creates an instance of Logger and returns it.
func NewLogger() *Logger {
	buf, err := comm.RandomBytes(4)
	if err != nil {
		panic(err)
	}
	return &Logger{
		ID: fmt.Sprintf("%08x", buf),
	}
}

// Logger provides methods for logging.
type Logger struct {
	ID string
}

// Log prints log.
func (p *Logger) Log(format string, args ...interface{}) {
	p.Logi(0, format, args...)
}

const timeFormat = "2006-01-02 15:04:05.000"

// Logi prints log with indent.
func (p *Logger) Logi(indent int, format string, args ...interface{}) {
	fmt.Println(time.Now().Format(timeFormat) + "/" + p.ID + ": " + strings.Repeat("    ", indent) + fmt.Sprintf(format, args...))
}
