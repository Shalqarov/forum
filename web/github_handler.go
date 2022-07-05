package web

import (
	"fmt"
	"net/http"
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
