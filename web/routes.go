package web

import "net/http"

// Routes - initialize routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/register", app.register)
	mux.HandleFunc("/registered", app.registered)
	mux.HandleFunc("/login", app.login)
	return mux
}
