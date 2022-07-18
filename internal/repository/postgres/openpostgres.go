package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

const (
	host     = "forumdb"
	port     = ":5432"
	user     = "forum"
	password = "mangothebest"
	dbname   = "forum"
)

func OpenDB(dsn string) (*sql.DB, error) {
	fmt.Println(host)
	// dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = setup(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setup(db *sql.DB) error {
	query, err := ioutil.ReadFile("./internal/repository/postgres/setup.sql")
	if err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("setup: %s", err)
	}
	return nil
}
