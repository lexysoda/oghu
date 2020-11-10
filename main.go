package main

import (
	"html/template"
	"os"

	"github.com/lexysoda/oghu/parser"
)

type Config struct {
	ContentDir string
	PublicDir  string
	Tpl        *template.Template
	Parser     parser.Parser
}

func main() {
	if err := os.RemoveAll("public"); err != nil {
		panic(err)
	}

	tpl := template.New("hi").Funcs(template.FuncMap{"urlify": Urlify})
	if _, err := tpl.ParseGlob("tpl/*"); err != nil {
		panic(err)
	}
	c := &Config{
		ContentDir: "content",
		PublicDir:  "public",
		Tpl:        tpl,
		Parser:     parser.New(),
	}
	s := NewSite(c)

	if err := s.ParseSite(); err != nil {
		panic(err)
	}
	if err := s.RenderSite(); err != nil {
		panic(err)
	}
}
