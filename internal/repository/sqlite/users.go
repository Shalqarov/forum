package sqlite

import (
	"database/sql"
	"os"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryCreateUser = `
	INSERT INTO "user"(
		"username",
		"email",
		"password",
		"avatar"
	) VALUES (?, ?, ?, ?)`

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

	queryChangeAvatar = `
	UPDATE "user" SET "avatar"=? WHERE "id" = ?
	`

	queryChangePassword = `
	UPDATE "user" SET "password"=? WHERE "id"=?
	`
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
	id, err := u.db.Exec(queryCreateUser, user.Username, user.Email, user.Password, user.Avatar)
	if err != nil {
		return 0, err
	}

	return id.LastInsertId()
}

func (u *sqliteRepo) ChangeAvatarByUserID(userID int64, image string) error {
	imagePath := ""
	err := u.db.QueryRow(`SELECT "avatar" FROM user WHERE "id" = ?`, userID).Scan(&imagePath)
	if err != nil {
		return err
	}
	if !strings.Contains(imagePath, "default-avatar.jpg") {
		err = os.Remove("ui" + imagePath)
		if err != nil {
			return err
		}
	}
	_, err = u.db.Exec(queryChangeAvatar, image, userID)
	if err != nil {
		return err
	}
	return nil
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
	err := u.db.QueryRow(queryGetUserByID, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar)
	return &user, err
}

func (u *sqliteRepo) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	err := u.db.QueryRow(queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Avatar,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *sqliteRepo) ChangePassword(newPassword string, userID int64) error {
	_, err := u.db.Exec(queryChangePassword, newPassword, userID)
	return err
}
