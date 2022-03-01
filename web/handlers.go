package web

import (
	"net/http"

	"github.com/Shalqarov/forum/pkg/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "home.page.html", &templateData{})
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		app.notFound(w)
		return
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "register.page.html", &templateData{})
	case http.MethodPost:
		user := models.User{
			Email: r.FormValue("email"),
			Login: r.FormValue("login"),
		}
		err := app.Forum.CreateUser(&user)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Application) registered(w http.ResponseWriter, r *http.Request) {
	users, _ := app.Forum.GetAllUsers()
	data := &templateData{
		Users: users,
	}
	app.render(w, r, "test.page.html", data)
}
