package web

import "net/http"

// Routes - initialize routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/signup", app.signup)
	mux.HandleFunc("/signin", app.signin)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/welcome", app.welcome)
	mux.HandleFunc("/createpost", app.createPost)
	return mux
}
