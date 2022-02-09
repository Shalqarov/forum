package web

import (
	"log"
)

// Application - web application dependencies
type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}
