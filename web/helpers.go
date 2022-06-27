package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)

func (app *Handler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Handler) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s doesn't exists", name))
		return
	}
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

func contentType(filebytes []byte) (string, error) {
	t := http.DetectContentType(filebytes)
	if strings.Contains(t, "image/jpeg") {
		return "jpeg", nil
	}
	if strings.Contains(t, "image/jpg") {
		return "jpg", nil
	}
	if strings.Contains(t, "image/png") {
		return "png", nil
	}
	if strings.Contains(t, "image/gif") {
		return "gif", nil
	}
	return "", errors.New("content is not an image")
}

func avatarType(filebytes []byte) (string, error) {
	t := http.DetectContentType(filebytes)
	if strings.Contains(t, "image/jpeg") {
		return "jpeg", nil
	}
	if strings.Contains(t, "image/jpg") {
		return "jpg", nil
	}
	if strings.Contains(t, "image/png") {
		return "png", nil
	}
	return "", errors.New("content is not an image")
}
