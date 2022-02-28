package sqlite

import (
	"database/sql"
	"errors"

	"github.com/Shalqarov/forum/pkg/models"
)

type Forum struct {
	DB *sql.DB
}

// CreateUser - new user
func (m *Forum) CreateUser(user *models.User) error {
	stmt := `INSERT INTO users"(
		"login",
		"email"
	) VALUES (?, ?)`

	_, err := m.DB.Exec(stmt, user.Login, user.Email)
	if err != nil {
		return err
	}

	return nil
}

// GetUserInfo...
func (m *Forum) GetUserInfo(login string) (*models.User, error) {
	statement := "SELECT * FROM users WHERE login = ?"
	row := m.DB.QueryRow(statement, login)
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

func (m *Forum) GetAllUsers() ([]*models.User, error) {
	statement := "SELECT * FROM users"
	rows, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*models.User{}
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
