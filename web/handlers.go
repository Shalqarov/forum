package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Shalqarov/forum/domain"
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
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/profile", h.profile)
	r.HandleFunc("/createpost", h.createPost)
	r.HandleFunc("/post/createcomment", h.createComment)
	r.HandleFunc("/post/vote", h.votePost)
	r.HandleFunc("/post/votecomment", h.voteComment)
	r.HandleFunc("/post", h.postPage)
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
	user := domain.User{}
	if isSession(r) {
		userID, err := getUserIDByCookie(r)
		if err != nil {
			log.Printf("home: GetUserIDByUsername: %s", err.Error())
			app.clientError(w, http.StatusInternalServerError)
			return
		}
		user.ID = userID
	}
	posts, err := app.PostUsecase.GetAllPosts()
	if err != nil {
		log.Println(err)
	}
	app.render(w, r, "home.page.html", &templateData{
		IsSession: isSession(r),
		User:      user,
		Posts:     posts,
	})
}
