package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Shalqarov/forum/pkg/models"
	uuid "github.com/satori/go.uuid"
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

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		app.notFound(w)
		return
	}
	switch r.Method {

	case http.MethodGet:
		app.render(w, r, "register.page.html", &templateData{})

	case http.MethodPost:

		user := models.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		err := app.Forum.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				app.render(w, r, "register.page.html", &templateData{
					Error: true,
				})
				return
			}
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

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
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
			fmt.Println("NO user")
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{})
			return
		}

		err = app.Forum.PasswordCompare(info, password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{})
			return
		}

		sessionToken := uuid.NewV4().String()
		expiresAt := time.Now().Add(120 * time.Second)

		sessions[sessionToken] = session{
			username: user.Username,
			expiry:   expiresAt,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresAt,
		})

		fmt.Println("success")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Application) welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	app.render(w, r, "welcome.page.html", &templateData{})
	// If the session is valid, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
}
