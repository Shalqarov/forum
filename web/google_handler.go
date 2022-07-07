package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	uuid "github.com/satori/go.uuid"
)

const (
	googleClientID     = "820533650499-dj70ovtt4uspgoh9sbdb0m3bdlsf470g.apps.googleusercontent.com"
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
)

func (app *Handler) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s",
		googleClientID,
		"http://localhost:5000/signin/google/callback",
		"profile email",
		"code",
	)

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (app *Handler) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	googleAccessToken, err := getGoogleAccessToken(code)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	googleData, err := getGoogleData(googleAccessToken)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if googleData.Username == "" || googleData.Email == "" {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	u := domain.User{
		Username: googleData.Username,
		Email:    googleData.Email,
		Password: uuid.NewV4().String(),
	}
	user, err := app.UserUsecase.GetUserByEmail(strings.ToLower(u.Email))
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

func getGoogleAccessToken(code string) (string, error) {
	u := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret},
		"redirect_uri":  {"http://localhost:5000/signin/google/callback"},
	}

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(u.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type googleAccesTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var ghresp googleAccesTokenResponse
	err = json.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}
	if ghresp.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	return ghresp.AccessToken, nil
}

func getGoogleData(accessToken string) (*domain.User, error) {
	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/oauth2/v3/userinfo?access_token="+accessToken,
		nil,
	)
	if err != nil {
		return nil, err
	}

	auth := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Username string `json:"given_name"`
		Email    string `json:"email"`
	}
	var data Data
	json.Unmarshal(body, &data)

	u := &domain.User{
		Username: data.Username,
		Email:    data.Email,
	}
	return u, nil
}

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
