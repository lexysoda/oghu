package oghu

import (
	"html/template"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	inputPath  = "fixtures/entry/content/catA/subCatB/1.md"
	relPath    = "catA/subCatB/1.md"
	outputPath = "fixtures/entry/public/catA/subCatB/this-is-title/index.html"
)

func TestParseEntry(t *testing.T) {
	d, err := time.Parse("02.01.2006", "30.12.1337")
	if err != nil {
		t.Fatalf("Failed to parse date: %s", err)
	}
	contentBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
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
		t.Fatalf("Failed to read file: %s", err)
	}
	p := GetParser()
	generated, err := ParseEntry(in, relPath, p)
	if err != nil {
		t.Fatalf("Failed to parse: %s", err)
	}
	if !reflect.DeepEqual(*generated, expected) {
		t.Errorf("Got: %+v\nExpected: %+v", *generated, expected)
	}
}
