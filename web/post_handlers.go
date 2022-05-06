package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/domain"
)

func (app *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		app.render(w, r, "createpost.page.html", &templateData{})
		return
	}

	userID, err := app.UserUsecase.GetUserIDByUsername(getUserNameByCookie(r))
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	postInfo := &domain.Post{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		UserID:   userID,
		Category: r.FormValue("category"),
	}
	if strings.TrimSpace(postInfo.Title) == "" || strings.TrimSpace(postInfo.Content) == "" || strings.TrimSpace(postInfo.Category) == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.PostUsecase.CreatePost(postInfo)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Handler) PostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	post, err := app.PostUsecase.GetPostByID(postID)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	comments, err := app.CommentUsecase.GetCommentsByPostID(postID)
	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	app.render(w, r, "post.page.html", &templateData{
		IsSession: isSession(r),
		Post:      post,
		Comments:  comments,
	})
}

func (app *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || postID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userName := getUserNameByCookie(r)
	userID, err := app.UserUsecase.GetUserIDByUsername(userName)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	comm := &domain.Comment{
		UserID:  userID,
		PostID:  postID,
		Author:  userName,
		Content: r.FormValue("comment"),
	}
	app.CommentUsecase.CreateComment(comm)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}
