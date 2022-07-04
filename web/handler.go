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
	r.HandleFunc("/signup", h.signup)
	r.HandleFunc("/signin", h.signin)
	r.HandleFunc("/signin/google/auth", h.googleAuthSignIn)
	r.HandleFunc("/signin/google/callback", h.googleSignIn)
	r.HandleFunc("/signup/google/auth", h.googleAuthSignUp)
	r.HandleFunc("/signup/google/callback", h.googleSignUp)
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/profile", h.profile)
	r.HandleFunc("/post", h.postPage)
	r.HandleFunc("/filter/category", h.postCategory)

	r.Handle("/profile/changepassword", middleware.SessionChecker(h.changePassword))
	r.Handle("/createpost", middleware.SessionChecker(h.createPost))
	r.Handle("/post/createcomment", middleware.SessionChecker(h.createComment))
	r.Handle("/post/votecomment", middleware.SessionChecker(h.voteComment))
	r.Handle("/post/vote", middleware.SessionChecker(h.votePost))
	r.Handle("/profile/upload-avatar", middleware.SessionChecker(h.uploadAvatar))

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
