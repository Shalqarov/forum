package web

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
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
	r.HandleFunc("/signin/google/auth", h.googleAuth)
	r.HandleFunc("/signin/google/callback", h.googleCallback)
	r.HandleFunc("/logout", h.logout)
	r.HandleFunc("/profile", h.profile)
	r.HandleFunc("/post", h.postPage)
	r.HandleFunc("/filter/category", h.postCategory)
	r.Handle("/createpost", sessionChecker(h.createPost))
	r.Handle("/post/createcomment", sessionChecker(h.createComment))
	r.Handle("/post/votecomment", sessionChecker(h.voteComment))
	r.Handle("/post/vote", sessionChecker(h.votePost))
	r.Handle("/profile/upload-avatar", sessionChecker(h.uploadAvatar))
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.Handle("/static", http.NotFoundHandler())
	r.Handle("/static/", http.StripPrefix("/static", fileServer))
}

func sessionChecker(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isSession(r) {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		h(w, r)
	}
}

func (app *Handler) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	image, err := createAvatar(r)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: createPost(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	userID, err := getUserIDByCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			app.clientError(w, http.StatusUnauthorized)
			return
		}
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	err = app.UserUsecase.ChangeAvatarByUserID(userID, image)
	if err != nil {
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusBadRequest)
		return
	}
	redirectID := strconv.Itoa(int(userID))
	http.Redirect(w, r, "/profile?id="+redirectID, http.StatusSeeOther)
}

func createAvatar(r *http.Request) (string, error) {
	file, _, err := r.FormFile("avatar")
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
	fileType, err := avatarType(fileBytes)
	if err != nil {
		return "", err
	}
	tempFile, err := ioutil.TempFile("./ui/static/images", "avatar-*."+fileType)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	tempFile.Write(fileBytes)
	return strings.ReplaceAll(tempFile.Name(), "./ui", ""), nil
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
	if isSession(r) {
		userID, err := getUserIDByCookie(r)
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
		if err != sql.ErrNoRows {
			app.ErrorLog.Println(err)
		}
		app.ErrorLog.Printf("HANDLERS: home(): %s", err.Error())
		app.clientError(w, http.StatusInternalServerError)
		return
	}
	app.render(w, r, "home.page.html", &templateData{
		IsSession: isSession(r),
		User:      user,
		Posts:     posts,
	})
}
