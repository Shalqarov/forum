package web

import (
	"fmt"
	"net/http"

	"github.com/Shalqarov/forum/pkg/models"
	"golang.org/x/crypto/bcrypt"
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
			Email:    r.FormValue("email"),
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
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

func (app *Application) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		app.notFound(w)
		return
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "login.page.html", &templateData{})
	case http.MethodPost:
		info := r.FormValue("email")
		password := r.FormValue("password")
		user, err := app.Forum.GetUserInfo(info)
		if err != nil {
			// ДОДЕЛАТЬ
			app.render(w, r, "login.page.html", &templateData{})
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			// ДОДЕЛАТЬ
			app.render(w, r, "login.page.html", &templateData{})
		}

		fmt.Println("success")
	}
}

func (app *Application) registered(w http.ResponseWriter, r *http.Request) {
	users, _ := app.Forum.GetAllUsers()
	data := &templateData{
		Users: users,
	}
	app.render(w, r, "test.page.html", data)
}
