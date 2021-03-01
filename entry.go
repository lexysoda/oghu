package oghu

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
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

func GenerateEntry(path, relPath string, c *Config) (*Entry, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	meta, err := c.Parser.Parse(source, &buf)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	category := filepath.Dir(relPath)
	return &Entry{
		Title:    meta.Title,
		Date:     meta.Date,
		Tags:     meta.Tags,
		Category: category,
		Path:     filepath.Join(category, Urlify(meta.Title)),
		Content:  template.HTML(buf.String()),
	}, nil
}

func (e *Entry) Render(c *Config) error {
	dirPath := filepath.Join(c.PublicDir, e.Category, e.Title)
	filePath := filepath.Join(dirPath, "index.html")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := c.Tpl.Lookup("entry").Execute(f, e); err != nil {
		return err
	}
	return nil
}
