package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	repository "github.com/Shalqarov/forum/repository/sqlite"
	"github.com/Shalqarov/forum/usecase"
	"github.com/Shalqarov/forum/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":8080", "Network address HTTP")
	dsn := flag.String("dsn", "forum.db", "Database name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbConn, err := repository.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	templateCache, err := web.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	userRepository := repository.NewSqliteRepo(dbConn)
	userUseCase := usecase.NewUserUsecase(userRepository)

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      web.NewUserHandler(userUseCase, templateCache),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	go web.ExpiredSessionsDeletion()
	err = srv.ListenAndServe()
	errorLog.Fatalln(err)
}
