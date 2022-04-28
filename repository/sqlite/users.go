package repository

import (
	"database/sql"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

type sqliteRepo struct {
	db *sql.DB
}

func NewSqliteUserRepo(db *sql.DB) domain.UserRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreateUser(user *domain.User) error {
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

	return err
}

func (u *sqliteRepo) GetUserIDByUsername(username string) (int, error) {
	stmt := `SELECT "id" FROM "user" WHERE "username"=?`
	user := domain.User{}
	err := u.db.QueryRow(stmt, username).Scan(&user.ID)
	if err != nil {
		return -1, domain.ErrNotFound
	}
	return user.ID, nil
}

func (u *sqliteRepo) GetUserByID(id int) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "id"=?`
	user := domain.User{}
	err := u.db.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return &user, err
}

func (u *sqliteRepo) GetUserByEmail(user *domain.User) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "email"=?`
	searchedUser := domain.User{}
	err := u.db.QueryRow(stmt, user.Email).Scan(&searchedUser.ID, &searchedUser.Username, &searchedUser.Email, &searchedUser.Password)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(searchedUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}
	return &searchedUser, nil
}
