package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultAvatarPath = "/static/images/Blank-profile.jpg"
)

func (app *Handler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || userID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	user, err := app.UserUsecase.GetUserByID(userID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: profile()1: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	posts, err := app.PostUsecase.GetPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: profile()2: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	likedPosts, err := app.PostUsecase.GetVotedPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: likedPosts(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	app.render(w, r, "profile.page.html", &templateData{
		IsSession:  isSession(r),
		User:       user,
		Posts:      posts,
		LikedPosts: likedPosts,
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
		addCookie(w, r, userID)
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
		addCookie(w, r, user.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		app.clientError(w, http.StatusMethodNotAllowed)
		return
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
