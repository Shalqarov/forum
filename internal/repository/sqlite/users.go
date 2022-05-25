package sqlite

import (
	"database/sql"

	"github.com/Shalqarov/forum/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	queryCreateUser = `
	INSERT INTO "user"(
		"username",
		"email",
		"password"
	) VALUES (?, ?, ?)`

	queryGetUserIDByUsername = `
	SELECT "id" 
	FROM "user" 
	WHERE "username"=?`

	queryGetUserByID = `
	SELECT * 
	FROM "user" 
	WHERE "id"=?`

	queryGetUserByEmail = `
	SELECT * 
	FROM "user"
	WHERE "email"=?`
)

type sqliteRepo struct {
	db *sql.DB
}

func (s *sqliteRepo) Close() {
	s.db.Close()
}

func NewSqliteUserRepo(db *sql.DB) domain.UserRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreateUser(user *domain.User) (int64, error) {
	id, err := u.db.Exec(queryCreateUser, user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	return id.LastInsertId()
}

func (u *sqliteRepo) GetUserIDByUsername(username string) (int64, error) {
	user := domain.User{}
	err := u.db.QueryRow(queryGetUserIDByUsername, username).Scan(&user.ID)
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (u *sqliteRepo) GetUserByID(id int64) (*domain.User, error) {
	user := domain.User{}
	err := u.db.QueryRow(queryGetUserByID, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return &user, err
}

func (u *sqliteRepo) GetUserByEmail(user *domain.User) (*domain.User, error) {
	searchedUser := domain.User{}
	err := u.db.QueryRow(queryGetUserByEmail, user.Email).Scan(&searchedUser.ID, &searchedUser.Username, &searchedUser.Email, &searchedUser.Password)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(searchedUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}
	return &searchedUser, nil
}
