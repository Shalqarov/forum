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

type googleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

var (
	loginConfig = googleConfig{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURI:  "http://localhost:5000/signin/google/callback",
	}
	registerConfig = googleConfig{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURI:  "http://localhost:5000/signup/google/callback",
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

func (app *Handler) googleRegisterHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s",
		googleClientID,
		"http://localhost:5000/signup/google/callback",
		"profile email",
		"code",
	)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (app *Handler) googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	u, err := getGoogleUserInfo(r, &loginConfig)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := app.UserUsecase.GetUserByEmail(strings.ToLower(u.Email))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			app.render(w, r, "login.page.html", &templateData{
				Error: "User doesn't exists",
			})
			return
		}
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	session.AddCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Handler) googleRegisterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	u, err := getGoogleUserInfo(r, &registerConfig)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: googleCallback(): %s", err.Error())
		app.clientError(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	app.setUser(w, r, u)
}

func getGoogleUserInfo(r *http.Request, c *googleConfig) (*domain.User, error) {
	code := r.URL.Query().Get("code")
	googleAccessToken, err := getGoogleAccessToken(code, c.RedirectURI)
	if err != nil {
		return nil, err
	}
	googleData, err := getGoogleData(googleAccessToken)
	if err != nil {
		return nil, err
	}
	if googleData.Username == "" || googleData.Email == "" {
		return nil, errors.New("data is nil")
	}
	u := &domain.User{
		Username: googleData.Username,
		Email:    googleData.Email,
		Password: uuid.NewV4().String(),
		Avatar:   defaultAvatarPath,
	}
	return u, nil
}

func getGoogleAccessToken(code, redirect_uri string) (string, error) {
	u := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret},
		"redirect_uri":  {redirect_uri},
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

	var ghresp Token
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
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data UserInfo
	json.Unmarshal(body, &data)

	u := &domain.User{
		Username: data.Username,
		Email:    data.Email,
	}
	return u, nil
}
