package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	postgres "github.com/Shalqarov/forum/internal/repository/postgres"
	repository "github.com/Shalqarov/forum/internal/repository/sqlite"
	"github.com/Shalqarov/forum/internal/session"
	"github.com/Shalqarov/forum/internal/usecase"
	"github.com/Shalqarov/forum/web"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":5000", "Network address HTTP")
	dsn := flag.String("dsn", "forum.db", "Database name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbConn, err := repository.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	postgre, err := postgres.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	router := http.NewServeMux()
	userRepo := postgres.NewPostgresUserRepo(postgre)
	postRepo := repository.NewSqlitePostRepo(dbConn)
	commRepo := repository.NewSqliteCommentRepo(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo)
	postUsecase := usecase.NewPostUsecase(postRepo)
	commUsecase := usecase.NewCommentUsecase(commRepo)

	templateCache, err := web.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	web.NewHandler(router, &web.Handler{
		UserUsecase:    userUsecase,
		PostUsecase:    postUsecase,
		CommentUsecase: commUsecase,
		TemplateCache:  templateCache,
		ErrorLog:       errorLog,
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	go session.ExpiredSessionsDeletion()
	err = srv.ListenAndServe()
	errorLog.Fatalln(err)
}
