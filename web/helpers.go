package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Handler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Handler) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
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
