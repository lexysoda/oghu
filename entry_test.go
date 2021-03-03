package oghu

import (
	"html/template"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	inputPath  = "fixtures/content/catA/subCatB/1.md"
	relPath    = "catA/subCatB/1.md"
	outputPath = "fixtures/public/catA/subCatB/this-is-title/index.html"
)

func TestParseEntry(t *testing.T) {
	d, err := time.Parse("02.01.2006", "30.12.1337")
	if err != nil {
		t.Errorf("Failed to parse date: %w", err)
	}
	contentBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Errorf("Failed to read file: %w", err)
	}
	expected := Entry{
		Title:    "this is title",
		Date:     d,
		Tags:     []string{"a", "b", "c"},
		Category: "catA/subCatB",
		Path:     "catA/subCatB/this-is-title",
		Content:  template.HTML(string(contentBytes)),
	}

	in, err := os.ReadFile(inputPath)
	if err != nil {
		t.Errorf("Failed to read file: %w", err)
	}
	p := GetParser()
	generated, err := ParseEntry(in, relPath, p)
	if err != nil {
		t.Errorf("Failed to parse: %w", err)
	}
	if !reflect.DeepEqual(*generated, expected) {
		t.Errorf("Got: %+v\nExpected: %+v", *generated, expected)
	}
}
