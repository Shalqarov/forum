package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	models "github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	uuid "github.com/satori/go.uuid"
)

const (
	githubLoginClientID     = "e6775f24b7ae5d23a9bf"
	githubLoginClientSecret = "9631fe307239dd7150112696a6aaf3557aac79c9"

	githubRegisterClientID     = "ad57429fdd16c9830c94"
	githubRegisterClientSecret = "b31feac706a76b9e5a79de7515ce3d7a82f33aae"
)

func (app *Handler) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		githubLoginClientID,
		"http://localhost:5000/signin/github/callback",
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *Handler) githubRegisterHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		githubRegisterClientID,
		"http://localhost:5000/signup/github/callback",
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *Handler) githubRegisterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken, err := getGithubAccesToken(code, githubRegisterClientID, githubRegisterClientSecret)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	u, err := getGithubUserInfo(r, accessToken)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	u.Password = uuid.NewV4().String()
	u.Avatar = defaultAvatarPath
	_, err = app.UserUsecase.GetUserByEmail(u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err := app.UserUsecase.CreateUser(u)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				app.ErrorLog.Println(err)
				app.render(w, r, "register.page.html", &templateData{
					Error: "User already exists",
				})
				return
			}
			session.AddCookie(w, r, userID)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.ErrorLog.Printf("HANDLERS: github: %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	app.render(w, r, "register.page.html", &templateData{
		Error: "User already exists",
	})
}

func (app *Handler) githubLoginCallBackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	accessToken, err := getGithubAccesToken(code, githubLoginClientID, githubLoginClientSecret)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	info, err := getGithubUserInfo(r, accessToken)
	if err != nil {
		app.ErrorLog.Println(err)
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	u := &models.User{
		Username: info.Username,
		Email:    info.Email,
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
		app.ErrorLog.Printf("HANDLERS: githubCallback(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	session.AddCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGithubUserInfo(r *http.Request, accessToken string) (*models.User, error) {
	username, err := getGithubData(accessToken)
	if err != nil {
		return nil, err
	}
	email, err := getGithubEmail(accessToken)
	if err != nil {
		return nil, err
	}
	userInfo := &models.User{
		Username: username,
		Email:    email,
	}
	if userInfo.Email == "" || userInfo.Username == "" {
		return nil, errors.New("getting empty login or email")
	}
	return userInfo, nil
}

func getGithubAccesToken(code, clientID, clientSecret string) (string, error) {
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}

	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ghresp Token
	err = json.Unmarshal(respBody, &ghresp)
	if err != nil {
		return "", err
	}
	return ghresp.AccessToken, nil
}

func getGithubData(accessToken string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return "", err
	}

	authorizationHeader := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeader)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Data struct {
		Username string `json:"login"`
		Email    string `json:"email"`
	}
	var data Data
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", err
	}
	return data.Username, nil
}

func getGithubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}

	authorizationHeader := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeader)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Data struct {
		Username string `json:"login"`
		Email    string `json:"email"`
	}
	var data []Data
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", err
	}
	return data[0].Email, nil
}
