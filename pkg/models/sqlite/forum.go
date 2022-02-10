package sqlite

import (
	"database/sql"
	"errors"

	"github.com/Shalqarov/forum/pkg/models"
)

type Forum struct {
	DB *sql.DB
}

// AddUser - new user
func (m *Forum) AddUser(login, email string) (int64, error) {
	statement := "INSERT INTO users (Login, Email) values(?,?)"

	result, err := m.DB.Exec(statement, login, email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserInfo ...
func (m *Forum) GetUserInfo(id int) (*models.User, error) {
	statement := "SELECT ID,Login,Email FROM Users WHERE ID = ?"
	row := m.DB.QueryRow(statement, id)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Login, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return u, nil
}
