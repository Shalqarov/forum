package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	models "github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	uuid "github.com/satori/go.uuid"
)

const (
	githubClientID     = "e6775f24b7ae5d23a9bf"
	githubClientSecret = "9631fe307239dd7150112696a6aaf3557aac79c9"
)

func (app *Handler) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		githubClientID,
		"http://localhost:5000/signin/github/callback",
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (app *Handler) githubCallBackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken, err := getGithubAccesToken(code)
	if err != nil {
		app.clientError(w, 400)
		return
	}

	githubLogin, err := getGithubData(githubAccessToken)
	if err != nil {
		app.clientError(w, 400)
		return
	}

	githubEmail, err := getGithubEmail(githubAccessToken)
	if err != nil {
		app.clientError(w, 400)
		return
	}

	if githubLogin == "" || githubEmail == "" {
		app.ErrorLog.Fatal("getting empty login or email ")
	}

	u := &models.User{
		Username: githubLogin,
		Email:    githubEmail,
		Password: uuid.NewV4().String(),
	}
	fmt.Println(u.Email)
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

func getGithubAccesToken(code string) (string, error) {
	requestBodyMap := map[string]string{
		"client_id":     githubClientID,
		"client_secret": githubClientSecret,
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
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var ghresp githubAccessTokenResponse
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
		Username string `json:"given_name"`
		Email    string `json:"email"`
	}
	var data []Data
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", err
	}
	return data[0].Email, nil
}
