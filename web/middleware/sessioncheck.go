package middleware

import (
	"net/http"

	"github.com/Shalqarov/forum/internal/session"
)

func NeedToBeAuthorized(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !session.IsSession(r) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		h(w, r)
	}
}

func Unauthorized(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if session.IsSession(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		h(w, r)
	}
}
