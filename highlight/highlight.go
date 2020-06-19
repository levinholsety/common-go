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

func (v TextType) String() string {
	switch v {
	case Normal:
		return "Normal"
	case DataType:
		return "DataType"
	case Keyword:
		return "Keyword"
	case Comment:
		return "Comment"
	case String:
		return "String"
	case Number:
		return "Number"
	default:
		return "Unknown"
	}
}

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
	DataTypes    []string   `json:"dataTypes"`
	Keywords     []string   `json:"keywords"`
	KeySymbols   []rune     `json:"keySymbols"`
	Blocks       TextBlocks `json:"blocks"`
	dataTypesMap map[string]int
	keywordsMap  map[string]int
}

// IsDataType returns true if string v is data type.
func (p *Config) IsDataType(v string) bool {
	if p.dataTypesMap == nil {
		p.dataTypesMap = make(map[string]int)
		for i, e := range p.DataTypes {
			p.dataTypesMap[strings.ToLower(e)] = i
		}
	}
	_, ok := p.dataTypesMap[strings.ToLower(v)]
	return ok
}

// IsKeyword returns true if string v is keyword.
func (p *Config) IsKeyword(v string) bool {
	if p.keywordsMap == nil {
		p.keywordsMap = make(map[string]int)
		for i, e := range p.Keywords {
			p.keywordsMap[strings.ToLower(e)] = i
		}
	}
	_, ok := p.keywordsMap[strings.ToLower(v)]
	return ok
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
	BeginIdentifier string   `json:"beginIdentifier"`
	EndIdentifier   string   `json:"endIdentifier"`
	EscapeChar      rune     `json:"escapeChar"`
	TextType        TextType `json:"textType"`
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
