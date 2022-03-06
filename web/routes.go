package web

import "net/http"

// Routes - initialize routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signup", app.signup)
	mux.HandleFunc("/signin", app.signin)
	mux.HandleFunc("/welcome", app.welcome)
	return mux
}
