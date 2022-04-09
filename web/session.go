package web

import (
	"net/http"
	"sync"
	"time"
)

var sessions = map[string]session{}

var cookie sync.Map

const cookieName string = "forum_session"

func interfaceToStruct(object interface{}) session {
	session, ok := object.(session)
	if ok {
		return session
	}
	return session
}

type session struct {
	username string
	expiry   time.Time
	mu       sync.Mutex
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	userSession, ok := cookie.Load(c.Value)
	if !ok {
		return false
	}
	userSessionStruct := interfaceToStruct(userSession)
	userSessionStruct.mu.Lock()
	if userSessionStruct.isExpired() {
		cookie.Delete(c.Value)
		userSessionStruct.mu.Unlock()
		return false
	}
	userSessionStruct.mu.Unlock()
	return ok
}
