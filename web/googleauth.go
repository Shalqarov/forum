package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	uuid "github.com/satori/go.uuid"
)

const (
	// ClientID     = "820533650499-t1lg2j1tl2162t2sldeo9tp3sj4itj3k.apps.googleusercontent.com"
	googleClientID = "820533650499-dj70ovtt4uspgoh9sbdb0m3bdlsf470g.apps.googleusercontent.com"
	// ClientSecret = "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"
	googleClientSecret = "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"
)

var (
	googleConfigSignIn = &Config{
		RedirectURL:  "http://localhost:5000/signin/google/callback",
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: googleEndPoint,
	}
	googleConfigSignUp = &Config{
		RedirectURL:  "http://localhost:5000/signup/google/callback",
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: googleEndPoint,
	}
	googleEndPoint = Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	}
	state = uuid.NewV4().String()
)

func (app *Handler) googleAuthSignIn(w http.ResponseWriter, r *http.Request) {
	if session.IsSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	url := googleConfigSignIn.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Handler) googleAuthSignUp(w http.ResponseWriter, r *http.Request) {
	if session.IsSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	url := googleConfigSignUp.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Handler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", errors.New("state is not valid"))
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	token, err := googleConfigSignIn.GetTokenByCode(r.FormValue("code"))
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userInfo := &UserInfo{}
	json.Unmarshal(content, &userInfo)
	user, err := app.UserUsecase.GetUserByEmail(userInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			app.render(w, r, "login.page.html", &templateData{
				Error: "User doesn't exists",
			})
			return
		}
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	session.AddCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Handler) googleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", errors.New("state is not valid"))
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	token, err := googleConfigSignUp.GetTokenByCode(r.FormValue("code"))
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response %s", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	userInfo := &UserInfo{}
	json.Unmarshal(content, &userInfo)

	_, err = app.UserUsecase.GetUserByEmail(userInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			u := domain.User{
				Email:    userInfo.Email,
				Username: userInfo.Username,
				Password: uuid.NewV4().String(),
				Avatar:   defaultAvatarPath,
			}
			userID, err := app.UserUsecase.CreateUser(&u)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			session.AddCookie(w, r, userID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	app.render(w, r, "register.page.html", &templateData{
		Error: "User already exists",
	})
}
