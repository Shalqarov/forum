package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Shalqarov/forum/domain"
)

type Handler struct {
	usecase       domain.Usecase
	TemplateCache map[string]*template.Template
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
}

func NewHandler(usecase domain.Usecase, template map[string]*template.Template) *http.ServeMux {
	handler := &Handler{
		usecase:       usecase,
		TemplateCache: template,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.home)
	mux.HandleFunc("/signup", handler.signup)
	mux.HandleFunc("/signin", handler.signin)
	mux.HandleFunc("/logout", handler.logout)
	mux.HandleFunc("/welcome", handler.welcome)
	mux.HandleFunc("/createpost", handler.createPost)
	return mux
}

func (app *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "home.page.html", &templateData{
			IsSession: isSession(r),
		})
	default:
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
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
		err := app.usecase.CreateUser(&user)
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

		user, err := app.usecase.GetUserByEmail(info)
		if err != nil {
			fmt.Println("wrong login or password")
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

func (app *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	switch r.Method {

	case http.MethodGet:
		app.render(w, r, "createpost.page.html", &templateData{})
	case http.MethodPost:
		userID, err := app.usecase.GetUserIDByUsername(getUserNameByCookie(r))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		postInfo := &domain.Post{
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
			UserID:   userID,
			Category: r.FormValue("category"),
		}
		fmt.Println(postInfo)
		err = app.usecase.CreatePost(postInfo)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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

	app.render(w, r, "welcome.page.html", &templateData{})
	w.Write([]byte(fmt.Sprintf("Welcome %s!", userName)))
}
