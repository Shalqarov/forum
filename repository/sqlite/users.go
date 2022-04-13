package repository

import (
	"database/sql"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

type sqliteUserRepo struct {
	db *sql.DB
}

func NewSqliteUserRepo(db *sql.DB) domain.UserRepo {
	return &sqliteUserRepo{db: db}
}

func (u *sqliteUserRepo) CreateUser(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO "user"(
		"username",
		"email",
		"password"
	) VALUES (?, ?, ?)`
	_, err = u.db.Exec(stmt, user.Username, user.Email, hashedPassword)
	if err != nil {
		return domain.ErrConflict
	}

	return nil
}

func (u *sqliteUserRepo) GetByID(id int) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "id"=?`
	user := domain.User{}
	err := u.db.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &user, nil
}

func (u *sqliteUserRepo) GetByEmail(user *domain.User) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "email"=?`
	searchedUser := domain.User{}
	err := u.db.QueryRow(stmt, user.Email).Scan(&searchedUser.ID, &searchedUser.Username, &searchedUser.Email, &searchedUser.Password)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &searchedUser, nil
}
