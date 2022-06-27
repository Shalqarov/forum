package web

import (
	"database/sql"
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
)

func imageUpload(r *http.Request) (string, error) {
	file, _, err := r.FormFile("myFile")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return "", nil
		}
		return "", err
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	fileType, err := contentType(fileBytes)
	if err != nil {
		return "", err
	}
	tempFile, err := ioutil.TempFile("./ui/static/images", "upload-*."+fileType)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	tempFile.Write(fileBytes)
	return strings.ReplaceAll(tempFile.Name(), "./ui", ""), nil
}

func (app *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	userID, err := getUserIDByCookie(r)
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
			IsSession: isSession(r),
		})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)

	user, err := app.UserUsecase.GetUserByID(userID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	postInfo := &domain.Post{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		UserID:   userID,
		Category: r.FormValue("category"),
		Author:   user.Username,
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
	if isSession(r) {
		_, err := getUserIDByCookie(r)
		if err != nil {
			log.Println("VotePost: GetUserIDByUsername: ", err)
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}
	user, err := app.UserUsecase.GetUserByID(post.UserID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: postPage(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.render(w, r, "post.page.html", &templateData{
		IsSession: isSession(r),
		User:      user,
		Post:      post,
		Comments:  comments,
	})
}

func (app *Handler) postCategory(w http.ResponseWriter, r *http.Request) {
	user := &domain.User{}
	if isSession(r) {
		userID, err := getUserIDByCookie(r)
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
		IsSession: isSession(r),
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
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
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
	userID, err := getUserIDByCookie(r)
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
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
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
	userID, err := getUserIDByCookie(r)
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
	if !isSession(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	postID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || postID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := getUserIDByCookie(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	user, err := app.UserUsecase.GetUserByID(userID)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createComment(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	comment := r.FormValue("comment")
	if len(comment) > 255 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	comm := &domain.Comment{
		UserID:     userID,
		PostID:     postID,
		Author:     user.Username,
		Content:    comment,
		UserAvatar: user.Avatar,
	}
	app.CommentUsecase.CreateComment(comm)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
}
