package oghu

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	goldParser "github.com/yuin/goldmark/parser"
)

type Meta struct {
	Title string
	Date  time.Time
	Tags  []string
}

type Parser struct {
	goldmark.Markdown
}

func GetParser() *Parser {
	return &Parser{goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
				),
			),
			meta.Meta,
			extension.DefinitionList,
		),
	)}
}

func (p *Parser) Parse(source []byte, writer io.Writer) (*Meta, error) {
	context := goldParser.NewContext()
	if err := p.Convert(source, writer, goldParser.WithContext(context)); err != nil {
		return nil, err
	}
	m := &Meta{}
	metaData := meta.Get(context)
	if s, ok := metaData["Title"].(string); ok {
		m.Title = s
	} else {
		return nil, fmt.Errorf("No Title")
	}
	if s, ok := metaData["Date"].(string); ok {
		if d, err := time.Parse("02.01.2006", s); err == nil {
			m.Date = d
		} else {
			return nil, err
		}
	}
	if s, ok := metaData["Tags"].(string); ok {
		m.Tags = strings.Split(s, " ")
	}

	return m, nil
}
