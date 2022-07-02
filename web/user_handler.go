package web

import (
	"net/http"
	"strconv"

	"github.com/Shalqarov/forum/internal/session"
)

const (
	defaultAvatarPath = "/static/images/Blank-profile.jpg"
)

func (app *Handler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || userID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	user, err := app.UserUsecase.GetUserByID(userID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: profile()1: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	posts, err := app.PostUsecase.GetPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: profile()2: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	likedPosts, err := app.PostUsecase.GetVotedPostsByUserID(user.ID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: likedPosts(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	app.render(w, r, "profile.page.html", &templateData{
		IsSession:  session.IsSession(r),
		User:       user,
		Posts:      posts,
		LikedPosts: likedPosts,
	})
}

func (app *Handler) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	image, err := createAvatar(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
			return
		}
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	err = app.UserUsecase.ChangeAvatarByUserID(userID, image)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	redirectID := strconv.Itoa(int(userID))
	http.Redirect(w, r, "/profile?id="+redirectID, http.StatusSeeOther)
}
