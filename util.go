package main

import (
	"regexp"
	"strings"
)

var (
	spaceRegex = regexp.MustCompile(`\s+`)
	wordRegex  = regexp.MustCompile(`[^\w\s]`)
)

func Urlify(name string) string {
	name = wordRegex.ReplaceAllLiteralString(name, "")
	name = spaceRegex.ReplaceAllLiteralString(name, "-")
	return strings.Trim(name, "-")
}

func Filter(entries []*Entry, cat string) []*Entry {
	filter := regexp.MustCompile(`^/*` + cat)
	filtered := []*Entry{}
	for _, e := range entries {
		if filter.MatchString(e.Category) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
