package web

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var cookie sync.Map

const cookieName string = "forum_session"

type duration struct {
	expiry map[interface{}]time.Time
	mu     sync.Mutex
}

var sessionDuration = duration{expiry: make(map[interface{}]time.Time)}

func addCookie(w http.ResponseWriter, r *http.Request, id int64) {
	sessionDuration.mu.Lock()
	defer sessionDuration.mu.Unlock()

	u := uuid.NewV4().String()
	deleteExistingCookie(id, u)

	cookie.Store(u, id)
	expire := time.Now().AddDate(0, 0, 1)
	sessionDuration.expiry[u] = expire

	session := &http.Cookie{
		Name:     cookieName,
		Value:    u,
		Path:     "/",
		HttpOnly: true,
		Expires:  expire,
	}
	http.SetCookie(w, session)
}

func deleteCookie(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(cookieName)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	cookie.Delete(c.Value)
}

func deleteExistingCookie(id int64, uuid string) {
	cookie.Range(func(key, value interface{}) bool {
		if id == value.(int64) {
			cookie.Delete(key)
		}
		return true
	})
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	var ok bool
	if err == nil {
		_, ok = cookie.Load(c.Value)
	}
	return ok
}

func getUserIDByCookie(r *http.Request) (int64, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return 0, err
	}
	value, ok := cookie.Load(c.Value)
	if !ok {
		return 0, fmt.Errorf("getUserIDByCookie: cannot load value from cookie store")
	}
	userID := value.(int64)
	return userID, nil
}

func ExpiredSessionsDeletion() {
	for {
		cookie.Range(func(key, value interface{}) bool {
			sessionDuration.mu.Lock()
			if time.Now().Unix() > sessionDuration.expiry[key].Unix() {
				cookie.Delete(key)
			}
			sessionDuration.mu.Unlock()
			return true
		})
		time.Sleep(time.Second * 5)
	}
}
