package highlight

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

var (
	reNumber = regexp.MustCompile(`^\d+(\.\d+)?$`)
)

// TextType represents type of text.
type TextType int

// TextTypes
const (
	Normal TextType = iota
	DataType
	Keyword
	Comment
	String
	Number
)

// Text represents text with type.
type Text struct {
	Text string
	Type TextType
}

// Config represents highlight config.
type Config struct {
	DataTypes  []string
	Keywords   []string
	KeySymbols []rune
	Blocks     TextBlocks
}

// IsDataType returns true if string v is data type.
func (p *Config) IsDataType(v string) bool {
	v = strings.ToLower(v)
	for _, e := range p.DataTypes {
		if e == v {
			return true
		}
	}
	return false
}

// IsKeyword returns true if string v is keyword.
func (p *Config) IsKeyword(v string) bool {
	v = strings.ToLower(v)
	for _, e := range p.Keywords {
		if e == v {
			return true
		}
	}
	return false
}

// IsKeySymbol returns true if string v is key symbol.
func (p *Config) IsKeySymbol(r rune) bool {
	for _, e := range p.KeySymbols {
		if e == r {
			return true
		}
	}
	return false
}

// TextBlock represents text block.
type TextBlock struct {
	BeginIdentifier string
	EndIdentifier   string
	EscapeChar      rune
	TextType        TextType
	escape          bool
}

// CheckBegin returns true if founds the begin of the TextBlock.
func (p *TextBlock) CheckBegin(buf *bytes.Buffer) bool {
	if buf.Len() < len(p.BeginIdentifier) {
		return false
	}
	return p.BeginIdentifier == string(buf.Bytes()[buf.Len()-len(p.BeginIdentifier):])
}

// CheckEnd returns true if founds the end of the TextBlock.
func (p *TextBlock) CheckEnd(buf *bytes.Buffer, r rune) bool {
	if p.escape {
		p.escape = false
		return false
	}
	if r == p.EscapeChar {
		p.escape = true
		return false
	}
	if buf.Len() < len(p.EndIdentifier) {
		return false
	}
	return p.EndIdentifier == string(buf.Bytes()[buf.Len()-len(p.EndIdentifier):])
}

// TextBlocks represents a group of TextBlock.
type TextBlocks []*TextBlock

// CheckBegin returns the matched TextBlock if founds the begin of any TextBlock.
func (p TextBlocks) CheckBegin(buf *bytes.Buffer) *TextBlock {
	for _, b := range p {
		if b.CheckBegin(buf) {
			return b
		}
	}
	return nil
}

// Parse parses the input text with specified config and returns the result.
func Parse(text string, cfg *Config) (result []*Text) {
	buf := &bytes.Buffer{}
	var block *TextBlock
	appendText := func(text string) {
		if cfg.IsDataType(text) {
			result = append(result, &Text{Text: text, Type: DataType})
		} else if cfg.IsKeyword(text) {
			result = append(result, &Text{Text: text, Type: Keyword})
		} else if reNumber.MatchString(text) {
			result = append(result, &Text{Text: text, Type: Number})
		} else {
			result = append(result, &Text{Text: text})
		}
	}
	for _, r := range text {
		buf.WriteRune(r)
		if block != nil {
			if block.CheckEnd(buf, r) {
				text := block.BeginIdentifier + buf.String()
				buf.Reset()
				result = append(result, &Text{Text: text, Type: block.TextType})
				block = nil
			}
		} else {
			block = cfg.Blocks.CheckBegin(buf)
			if block != nil {
				text := buf.String()
				buf.Reset()
				text = text[:len(text)-len(block.BeginIdentifier)]
				appendText(text)
			} else if cfg.IsKeySymbol(r) || unicode.IsSpace(r) {
				text := buf.String()
				buf.Reset()
				ch := string([]rune{r})
				text = text[:len(text)-len(ch)]
				appendText(text)
				result = append(result, &Text{Text: ch, Type: Keyword})
			}
		}
	}
	return
}
