package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Shalqarov/forum/internal/domain"
	uuid "github.com/satori/go.uuid"
)

const (
	// ClientID     = "820533650499-t1lg2j1tl2162t2sldeo9tp3sj4itj3k.apps.googleusercontent.com"
	ClientID = "820533650499-dj70ovtt4uspgoh9sbdb0m3bdlsf470g.apps.googleusercontent.com"
	// ClientSecret = "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"
	ClientSecret = "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"
)

var (
	googleOauthConfig = &Config{
		RedirectURL:  "http://localhost:5000/signin/google/callback",
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
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

func (app *Handler) googleAuth(w http.ResponseWriter, r *http.Request) {
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", errors.New("state is not valid"))
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	token, err := googleOauthConfig.GetTokenByCode(r.FormValue("code"))
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
	fmt.Println(string(content))
	userInfo := &UserInfo{}
	json.Unmarshal(content, &userInfo)
	fmt.Println(userInfo.Email)
	fmt.Println(userInfo.Username)

	user, err := app.UserUsecase.GetUserByEmail(userInfo.Email)
	fmt.Println(user, err)
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
			addCookie(w, r, userID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	addCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
