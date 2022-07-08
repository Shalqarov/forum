package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/Shalqarov/forum/internal/domain"
	"github.com/Shalqarov/forum/internal/session"
	"github.com/Shalqarov/forum/web/middleware"
)

type Handler struct {
	UserUsecase    domain.UserUsecase
	PostUsecase    domain.PostUsecase
	CommentUsecase domain.CommentUsecase
	TemplateCache  map[string]*template.Template
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
}

func NewHandler(r *http.ServeMux, h *Handler) {
	r.HandleFunc("/", h.home)
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/profile", h.profile)
	r.HandleFunc("/post", h.postPage)
	r.HandleFunc("/filter/category", h.postCategory)

	r.Handle("/signup", middleware.Unauthorized(h.signup))
	r.Handle("/signin", middleware.Unauthorized(h.signin))

	r.Handle("/signin/google/auth", middleware.Unauthorized(h.googleLoginHandler))
	r.Handle("/signin/google/callback", middleware.Unauthorized(h.googleLoginCallbackHandler))
	r.Handle("/signup/google/auth", middleware.Unauthorized(h.googleRegisterHandler))
	r.Handle("/signup/google/callback", middleware.Unauthorized(h.googleRegisterCallbackHandler))

	r.Handle("/signin/github/auth", middleware.Unauthorized(h.githubLoginHandler))
	r.Handle("/signin/github/callback", middleware.Unauthorized(h.githubLoginCallBackHandler))
	r.Handle("/signup/github/auth", middleware.Unauthorized(h.githubRegisterHandler))
	r.Handle("/signup/github/callback", middleware.Unauthorized(h.githubRegisterCallbackHandler))

	r.Handle("/post/vote", middleware.NeedToBeAuthorized(h.votePost))
	r.Handle("/createpost", middleware.NeedToBeAuthorized(h.createPost))
	r.Handle("/post/votecomment", middleware.NeedToBeAuthorized(h.voteComment))
	r.Handle("/post/createcomment", middleware.NeedToBeAuthorized(h.createComment))
	r.Handle("/profile/upload-avatar", middleware.NeedToBeAuthorized(h.uploadAvatar))
	r.Handle("/profile/changepassword", middleware.NeedToBeAuthorized(h.changePassword))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static", http.NotFoundHandler())
	r.Handle("/static/", http.StripPrefix("/static", fileServer))
}

func (app *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	user := &domain.User{}
	if session.IsSession(r) {
		userID, err := session.GetUserIDByCookie(r)
		if err != nil {
			if err == http.ErrNoCookie {
				app.clientError(w, http.StatusUnauthorized)
				return
			}
			app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		user.ID = userID
	}
	posts, err := app.PostUsecase.GetAllPosts()
	if err != nil {
		if err == sql.ErrNoRows {
			app.InfoLog.Println(err)
		}
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	app.render(w, r, "home.page.html", &templateData{
		IsSession: session.IsSession(r),
		User:      user,
		Posts:     posts,
	})
}
