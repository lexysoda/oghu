package oghu

import (
	"os"
	"path/filepath"
	"sort"
)

type Site struct {
	Entries []*Entry
	Tags    map[string][]*Entry
	C       *Config
}

func NewSite(c *Config) *Site {
	return &Site{
		Tags: map[string][]*Entry{},
		C:    c,
	}
}

func (s *Site) ParseSite() error {
	if err := filepath.Walk(s.C.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.C.ContentDir, path)
		if err != nil {
			return err
		}
		pubPath := filepath.Join(s.C.PublicDir, relPath)
		if info.IsDir() {
			return os.MkdirAll(pubPath, info.Mode())
		} else if filepath.Ext(path) != ".md" {
			return os.Link(path, pubPath)
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entry, err := ParseEntry(bytes, relPath, s.C.Parser)
		if err != nil {
			return err
		}
		s.Entries = append(s.Entries, entry)
		for _, t := range entry.Tags {
			s.Tags[t] = append(s.Tags[t], entry)
		}

		return nil
	}); err != nil {
		return err
	}
	sort.Slice(s.Entries, func(i, j int) bool {
		return s.Entries[j].Date.After(s.Entries[i].Date)
	})
	return nil
}

func (s *Site) RenderSite() error {
	data := map[string]interface{}{
		"Meta": s.C.Meta,
	}
	for _, e := range s.Entries {
		data["Entry"] = e
		Render(s.C, e.Path, "entry", data)
	}
	delete(data, "Entry")

	categories := []string{"", "blog", "bits"}
	for _, cat := range categories {
		data["Entries"] = Filter(s.Entries, cat)
		data["Title"] = cat
		if cat == "" {
			data["Title"] = "posts"
		}
		if err := Render(s.C, "/"+cat, "list", data); err != nil {
			return err
		}
	}
	delete(data, "Entries")

	data["Tags"] = s.Tags
	data["Title"] = "tags"
	if err := Render(s.C, "/tags", "list", data); err != nil {
		return err
	}
	delete(data, "Tags")

	for tag, entries := range s.Tags {
		data["Entries"] = entries
		data["Title"] = tag
		if err := Render(s.C, "/tags/"+tag, "list", data); err != nil {
			return err
		}
	}

	return nil
}

func Render(c *Config, path, tpl string, data interface{}) error {
	dirPath := filepath.Join(c.PublicDir, path)
	filePath := filepath.Join(dirPath, "index.html")

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := c.Tpl.Lookup(tpl).Execute(f, data); err != nil {
		return err
	}
	return nil
}
