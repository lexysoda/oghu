package oghu

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ContentDir   string            `yaml:"ContentDir"`
	PublicDir    string            `yaml:"PublicDir"`
	TemplateGlob string            `yaml:"TemplateGlob"`
	Meta         map[string]string `yaml:"Meta"`
	Tpl          *template.Template
	Parser       *Parser
}

func parseConfig(confPath string) (*Config, error) {
	b, err := ioutil.ReadFile("oghu.yaml")
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func Oghu() (string, error) {
	c, err := parseConfig("config.yaml")
	if err != nil {
		return "", err
	}

	if err := os.RemoveAll(c.PublicDir); err != nil {
		return "", err
	}

	c.Tpl = template.New("hi").Funcs(template.FuncMap{"urlify": Urlify})
	if _, err := c.Tpl.ParseGlob(c.TemplateGlob); err != nil {
		return "", err
	}
	c.Parser = GetParser()

	s := NewSite(c)

	if err := s.ParseSite(); err != nil {
		return "", err
	}
	if err := s.RenderSite(); err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully rendered %s to %s.", c.ContentDir, c.PublicDir), err
}

func New(path string, data map[string]interface{}) (string, error) {
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("Failed to open %s: %w", path, err)
	}
	if err := template.Must(template.New("kek").Parse(`---
Title: {{ .Title }}
Date: {{ .Date }}
Tags:
---

`)).Execute(file, data); err != nil {
		return "", fmt.Errorf("Failed to execute template: %w", err)
	}
	return fmt.Sprintf("Successfully created %s", path), nil
}
