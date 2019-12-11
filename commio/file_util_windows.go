package commio

import "bytes"

// EscapeFileName escapes file name.
func EscapeFileName(name string) string {
	var buf bytes.Buffer
	for _, ch := range name {
		switch ch {
		case '\\':
			buf.WriteString("%5c")
		case '/':
			buf.WriteString("%2f")
		case ':':
			buf.WriteString("%3a")
		case '*':
			buf.WriteString("%2a")
		case '?':
			buf.WriteString("%3f")
		case '"':
			buf.WriteString("%22")
		case '<':
			buf.WriteString("%3c")
		case '>':
			buf.WriteString("%3e")
		case '|':
			buf.WriteString("%7c")
		default:
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
