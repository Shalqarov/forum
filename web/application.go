package web

import (
	"html/template"
	"log"

	"github.com/Shalqarov/forum/pkg/models/sqlite"
)

// Application - web application dependencies
type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Forum         *sqlite.Forum
	TemplateCache map[string]*template.Template
}
