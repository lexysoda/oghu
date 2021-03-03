package oghu

import (
	"crypto/md5"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"testing"
)

const (
	content   = "fixtures/content"
	public    = "fixtures/public"
	templates = "fixtures/tpl/*"
	tempGlob  = "oghu-public-test-*"
)

func TestSite(t *testing.T) {
	outputDir, err := ioutil.TempDir(".", tempGlob)
	defer os.RemoveAll(outputDir)
	if err != nil {
		t.Fatalf("Couldn't create tempdir: %s", err)
	}
	tpl, err := template.New("hi").Funcs(template.FuncMap{"urlify": Urlify}).ParseGlob(templates)
	if err != nil {
		t.Fatalf("Failed creating templates: %s", err)
	}

	c := &Config{
		ContentDir: content,
		PublicDir:  outputDir,
		Tpl:        tpl,
		Parser:     GetParser(),
	}

	s := NewSite(c)
	if err := s.ParseSite(); err != nil {
		t.Fatalf("Failed parsing site: %s", err)
	}
	if err := s.RenderSite(); err != nil {
		t.Fatalf("Failed rendering site: %s", err)
	}

	generated, err := hashDir(outputDir)
	if err != nil {
		t.Fatalf("Failed hashing %s: %s", outputDir, err)
	}
	wanted, err := hashDir(public)
	if err != nil {
		t.Fatalf("Failed hashing %s: %s", public, err)
	}

	if generated != wanted {
		t.Errorf("Hash mismatch: got: %s wanted: %s", generated, wanted)
	}
}

func hashDir(dir string) (string, error) {
	h := md5.New()
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			io.WriteString(h, relPath)
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}
