package oghu

import (
	"bytes"
	"html/template"
	"path/filepath"
	"time"
)

type Entry struct {
	Title    string
	Date     time.Time
	Tags     []string
	Content  template.HTML
	Category string
	Path     string
}

func ParseEntry(b []byte, path string, p *Parser) (*Entry, error) {
	var buf bytes.Buffer
	meta, err := p.Parse(b, &buf)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	category := filepath.Dir(path)
	return &Entry{
		Title:    meta.Title,
		Date:     meta.Date,
		Tags:     meta.Tags,
		Category: category,
		Path:     filepath.Join(category, Urlify(meta.Title)),
		Content:  template.HTML(buf.String()),
	}, nil
}
