package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/domain"
)

func (app *Handler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || userID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	user, err := app.UserUsecase.GetUserByID(userID)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	posts, _ := app.PostUsecase.GetAllPostsByUserID(user.ID)
	app.render(w, r, "profile.page.html", &templateData{
		IsSession: isSession(r),
		User:      *user,
		Posts:     posts,
	})
}

func (app *Handler) signup(w http.ResponseWriter, r *http.Request) {
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "register.page.html", &templateData{})
		return
	case http.MethodPost:
		user := domain.User{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Username) == "" {
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "register.page.html", &templateData{
				Error: "incorrect input",
			})
			return
		}

		err := app.UserUsecase.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				app.render(w, r, "register.page.html", &templateData{
					Error: "User already exists",
				})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "register.page.html", &templateData{
				Error: "bad request",
			})
			return
		}
		addCookie(w, r, user.Username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Handler) signin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		app.notFound(w)
		return
	}
	switch r.Method {

	case http.MethodGet:
		app.render(w, r, "login.page.html", &templateData{})
	case http.MethodPost:
		info := &domain.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		if strings.TrimSpace(info.Email) == "" || strings.TrimSpace(info.Password) == "" {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		user, err := app.UserUsecase.GetUserByEmail(info)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{
				Error: "email or password are wrong",
			})
			return
		}

		addCookie(w, r, user.Username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Handler) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.clientError(w, http.StatusBadRequest)
		return
	}
	deleteCookie(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
