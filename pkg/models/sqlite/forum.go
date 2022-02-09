package sqlite

import "database/sql"

type Forum struct {
	DB *sql.DB
}

func (m *Forum) AddUser(login, email string) (int64, error) {

}
