package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultAvatarPath = "/static/images/default-avatar.jpg"
)

func (app *Handler) signup(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "register.page.html", &templateData{
				Error: "Credentials are empty",
			})
			return
		}
		app.setUser(w, r, &user)
	default:
		app.methodNotAllowed(w, r)
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
			app.clientError(w, r, http.StatusBadRequest, "Some of the credentials are empty")
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
		app.methodNotAllowed(w, r)
		return
	}
}

func (app *Handler) setUser(w http.ResponseWriter, r *http.Request, u *domain.User) {
	_, err := app.UserUsecase.GetUserByEmail(u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err := app.UserUsecase.CreateUser(u)
			if err != nil {
				app.ErrorLog.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				app.render(w, r, "register.page.html", &templateData{
					Error: "User already exists",
				})
				return
			}
			session.AddCookie(w, r, userID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	app.render(w, r, "register.page.html", &templateData{
		Error: "User already exists",
	})
}

func (app *Handler) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(session.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.clientError(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	session.DeleteCookie(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Handler) changePassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "changepass.page.html", &templateData{})
		return
	case http.MethodPost:
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		if strings.TrimSpace(password) == "" || strings.TrimSpace(confirmPassword) == "" {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "changepass.page.html", &templateData{Error: "Some of the credentials are empty"})
			return
		}
		if password != confirmPassword {
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "changepass.page.html", &templateData{
				Error: "passwords not similar",
			})
			return
		}
		userID, err := session.GetUserIDByCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{
				Error: "You must be logged in to change password",
			})
			return
		}
		err = app.UserUsecase.ChangePassword(password, userID)
		if err != nil {
			app.ErrorLog.Printf("HANDLERS: changePass(): %s", err.Error())
			app.clientError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", userID), http.StatusSeeOther)
	default:
		app.methodNotAllowed(w, r)
		return
	}
}
