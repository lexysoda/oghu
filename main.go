package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/lexysoda/oghu/parser"
)

type Config struct {
	ContentDir   string            `yaml:"ContentDir"`
	PublicDir    string            `yaml:"PublicDir"`
	TemplateGlob string            `yaml:"TemplateGlob"`
	Meta         map[string]string `yaml:"Meta"`
	Tpl          *template.Template
	Parser       parser.Parser
}

func main() {
	b, err := ioutil.ReadFile("oghu.yaml")
	if err != nil {
		panic(err)
	}
	c := &Config{}
	if yaml.Unmarshal(b, c); err != nil {
		panic(err)
	}
	log.Println(c)

	if err := os.RemoveAll(c.PublicDir); err != nil {
		panic(err)
	}

	c.Tpl = template.New("hi").Funcs(template.FuncMap{"urlify": Urlify})
	if _, err := c.Tpl.ParseGlob(c.TemplateGlob); err != nil {
		panic(err)
	}
	c.Parser = parser.New()

	s := NewSite(c)

	if err := s.ParseSite(); err != nil {
		panic(err)
	}
	if err := s.RenderSite(); err != nil {
		panic(err)
	}
}
