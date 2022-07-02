package web

import (
	"net/http"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	"golang.org/x/crypto/bcrypt"
)

func (app *Handler) signup(w http.ResponseWriter, r *http.Request) {
	if session.IsSession(r) {
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
			Avatar:   defaultAvatarPath,
		}

		if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Username) == "" {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		userID, err := app.UserUsecase.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				app.render(w, r, "register.page.html", &templateData{
					Error: "User already exists",
				})
				return
			}
			app.clientError(w, http.StatusBadRequest)
			return
		}
		session.AddCookie(w, r, userID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Handler) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "login.page.html", &templateData{})
		return
	case http.MethodPost:
		info := &domain.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		if strings.TrimSpace(info.Email) == "" || strings.TrimSpace(info.Password) == "" {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		user, err := app.UserUsecase.GetUserByEmail(info.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{
				Error: "user doesn't exists",
			})
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{
				Error: "email or password are wrong",
			})
			return
		}
		session.AddCookie(w, r, user.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *Handler) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(session.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.clientError(w, http.StatusBadRequest)
		return
	}
	session.DeleteCookie(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
