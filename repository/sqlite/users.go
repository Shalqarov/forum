package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	models "github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser - new user
func (m *Forum) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO "main"."users"(
		"username",
		"email",
		"password"
	) VALUES (?, ?, ?)`

	_, err = m.DB.Exec(stmt, user.Username, user.Email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *Forum) PasswordCompare(login, password string) error {
	s := `SELECT "password" FROM "users" 
	WHERE "username"=? OR "email"=?`
	row := m.DB.QueryRow(s, login, login)
	u := &models.User{}
	err := row.Scan(&u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// GetUserInfo...
func (m *Forum) GetUserInfo(login string) (*models.User, error) {
	statement := "SELECT * FROM users WHERE \"username\" = ? OR \"email\" = ?"
	row := m.DB.QueryRow(statement, login, login)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return u, nil
}
