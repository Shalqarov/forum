package web

import (
	"html/template"
	"path/filepath"

	models "github.com/Shalqarov/forum/tree/master/domain"
)

type templateData struct {
	User      models.User
	Error     string
	IsSession bool
	Post      *models.Post
	Posts     []*models.PostDTO
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
