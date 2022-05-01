package web

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	fmt.Println(posts)
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
			Username: r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		if strings.TrimSpace(user.Email) != "" || strings.TrimSpace(user.Password) != "" || strings.TrimSpace(user.Username) != "" {
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "register.page.html", &templateData{
				Error: true,
			})
			return
		}

		err := app.UserUsecase.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				app.render(w, r, "register.page.html", &templateData{
					Error: true,
				})
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			app.render(w, r, "register.page.html", &templateData{
				Error: true,
			})
			return
		}
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
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusUnauthorized)
				app.render(w, r, "login.page.html", &templateData{})
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{})
			return
		}

		addCookie(w, r, user.Username)

		log.Println("success signin - ", user.Username)
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

func (app *Handler) welcome(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	_, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userName := getUserNameByCookie(r)

	w.Write([]byte(fmt.Sprintf("Welcome %s!", userName)))
}