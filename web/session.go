package web

import (
	"net/http"
	"time"
)

var sessions = map[string]session{}

// var cookie sync.Map

const cookieName string = "forum_session"

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	var ok bool
	if err == nil {
		_, ok = sessions[c.Value]
	}
	return ok
}
