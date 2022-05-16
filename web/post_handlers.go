package web

import (
	"errors"
	"fmt"
	"log"
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
		fmt.Println("AAAAAAAAAAAAA")
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.PostUsecase.CreatePost(postInfo)
	if err != nil {
		log.Printf("CreatePost: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Handler) postPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || postID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	votes, err := app.PostUsecase.GetVotesCountByPostID(postID)
	if err != nil {
		log.Println("getVotes error:")
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
		Votes:     votes,
	})
}

func (app *Handler) votePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || postID < 1 {
		log.Println("VotePost:", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	vote, err := strconv.ParseInt(r.URL.Query().Get("vote"), 10, 64)
	if err != nil || vote != 1 && vote != -1 {
		log.Println("VotePost:", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := app.UserUsecase.GetUserIDByUsername(getUserNameByCookie(r))
	if err != nil {
		log.Println("VotePost:", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.PostUsecase.VotePost(postID, userID, int(vote))
	if err != nil {
		log.Println("VotePost:", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}

func (app *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
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
	comment := r.FormValue("comment")
	if len(comment) > 255 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	comm := &domain.Comment{
		UserID:  userID,
		PostID:  postID,
		Author:  userName,
		Content: comment,
	}
	app.CommentUsecase.CreateComment(comm)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}
