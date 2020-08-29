package app

import (
	"path/filepath"
	"text/template"
)

type templateCache map[string]*template.Template

func newTemplateCache(files []string) (templateCache, error) {
	cache := templateCache{}

	for _, f := range files {
		name := filepath.Base(f)
		t, err := template.New(name).ParseFiles(f)
		if err != nil {
			return nil, err
		}

		cache[name] = t
	}

	return cache, nil
}
