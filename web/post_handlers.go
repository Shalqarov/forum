package web

import (
	"errors"
	"net/http"
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