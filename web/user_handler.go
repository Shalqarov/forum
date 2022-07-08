package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
)

func (app *Handler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w, r)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id < 1 {
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	user, err := app.UserUsecase.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.ErrorLog.Printf("HANDLERS: profile(): %s", err.Error())
			app.clientError(w, r, http.StatusNotFound, "User doesn't exist :(")
			return
		}
		app.ErrorLog.Printf("HANDLERS: profile(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	posts, err := app.PostUsecase.GetPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: profile(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	likedPosts, err := app.PostUsecase.GetVotedPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: likedPosts(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		app.render(w, r, "login.page.html", &templateData{})
		return
	}
	app.render(w, r, "profile.page.html", &templateData{
		IsSession: session.IsSession(r),
		User: &domain.User{
			ID: userID,
		},
		Profile:    user,
		Posts:      posts,
		LikedPosts: likedPosts,
	})
}

func (app *Handler) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w, r)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	image, err := createAvatar(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
			return
		}
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	err = app.UserUsecase.ChangeAvatarByUserID(userID, image)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile?id=%d", userID), http.StatusSeeOther)
}
