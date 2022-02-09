package pkg

import (
	"database/sql"
	"io/ioutil"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
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
	query, err := ioutil.ReadFile("./pkg/setup.sql")
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(query)); err != nil {
		return err
	}
	return nil
}
