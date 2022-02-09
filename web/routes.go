package web

import "net/http"

// Routes - initialize routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	return mux
}
