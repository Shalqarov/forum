package web

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
)

func (app *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		app.render(w, r, "createpost.page.html", &templateData{
			User:      &domain.User{ID: userID},
			IsSession: session.IsSession(r),
		})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)

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
	postInfo.Image, err = imageUpload(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	_, err = app.PostUsecase.CreatePost(postInfo)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
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
		app.ErrorLog.Printf("HANDLERS: postPage(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := app.PostUsecase.GetPostByID(postID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: postPage(): %s", err.Error())
		app.clientError(w, http.StatusNotFound)
		return
	}

	comments, err := app.CommentUsecase.GetCommentsByPostID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.ErrorLog.Printf("HANDLERS: postPage(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	tempUser := &domain.User{}
	if session.IsSession(r) {
		tempUser.ID, err = session.GetUserIDByCookie(r)
		if err != nil {
			log.Println("VotePost: GetUserIDByUsername: ", err)
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}

	app.render(w, r, "post.page.html", &templateData{
		IsSession: session.IsSession(r),
		User:      tempUser,
		Post:      post,
		Comments:  comments,
	})
}

func (app *Handler) postCategory(w http.ResponseWriter, r *http.Request) {
	user := &domain.User{}
	if session.IsSession(r) {
		userID, err := session.GetUserIDByCookie(r)
		if err != nil {
			app.ErrorLog.Printf("postCategory: getUserIDByCookie: %s", err.Error())
			app.clientError(w, http.StatusBadRequest)
			return
		}
		user.ID = userID
	}
	category := r.URL.Query().Get("category")
	posts, err := app.PostUsecase.GetPostsByCategory(category)
	if err != nil {
		app.ErrorLog.Printf("postCategory: GetPostsByCategory: %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	app.render(w, r, "home.page.html", &templateData{
		IsSession: session.IsSession(r),
		User:      user,
		Posts:     posts,
	})
}

func (app *Handler) votePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || postID < 1 {
		app.ErrorLog.Printf("HANDLERS: votePost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	vote, err := strconv.ParseInt(r.URL.Query().Get("vote"), 10, 64)
	if err != nil || vote != 1 && vote != -1 {
		app.ErrorLog.Printf("HANDLERS: votePost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: votePost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.PostUsecase.VotePost(postID, userID, int(vote))
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: votePost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}

func (app *Handler) voteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: voteComment(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	vote, err := strconv.ParseInt(r.URL.Query().Get("vote"), 10, 64)
	if err != nil || vote != 1 && vote != -1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	commentID, err := strconv.ParseInt(r.URL.Query().Get("comm"), 10, 64)
	if err != nil || vote != 1 && vote != -1 {
		app.ErrorLog.Printf("HANDLERS: voteComment(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserIDByCookie(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: voteComment(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.CommentUsecase.VoteComment(commentID, userID, int(vote))
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: voteComment(): %s", err.Error())
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
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || postID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := session.GetUserIDByCookie(r)
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
		Content: comment,
	}
	app.CommentUsecase.CreateComment(comm)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}
