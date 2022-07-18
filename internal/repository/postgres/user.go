package postgres

import (
	"database/sql"
	"os"
	"strings"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryCreateUser = `
	INSERT INTO "user" (
		"username",
		"email",
		"password",
		"avatar"
	) VALUES ($1, $2, $3, $4, $5) RETURNING "user_id"`

	queryGetUserIDByUsername = `
	SELECT "user_id" 
	FROM "user" 
	WHERE "username"=$1`

	queryGetUserByID = `
	SELECT * 
	FROM "user" 
	WHERE "user_id"=$1`

	queryGetUserByEmail = `
	SELECT * 
	FROM "user"
	WHERE "email"=$1`

	queryChangeAvatar = `
	UPDATE "user" 
	SET "avatar"=$1 
	WHERE "user_id" = $2`

	queryChangePassword = `
	UPDATE "user" 
	SET "password"=$1 
	WHERE "user_id"=$2`
)

type repo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) domain.UserRepo {
	return &repo{db: db}
}

func (u *repo) CreateUser(user *domain.User) (int64, error) {
	var lastInsertId int64
	err := u.db.QueryRow(
		queryCreateUser,
		user.Username,
		user.Password,
		user.Avatar,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, nil
	}
	return lastInsertId, nil
}

func (u *repo) ChangeAvatarByUserID(userID int64, image string) error {
	imagePath := ""
	err := u.db.QueryRow(
		`SELECT "avatar" 
		FROM user 
		WHERE "user_id" = $1`,
		userID).Scan(&imagePath)
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

func (u *repo) GetUserIDByUsername(username string) (int64, error) {
	user := domain.User{}
	err := u.db.QueryRow(queryGetUserIDByUsername, username).Scan(&user.ID)
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (u *repo) GetUserByID(id int64) (*domain.User, error) {
	user := domain.User{}
	err := u.db.QueryRow(queryGetUserByID, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar)
	return &user, err
}

func (u *repo) GetUserByEmail(email string) (*domain.User, error) {
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

func (u *repo) ChangePassword(newPassword string, userID int64) error {
	_, err := u.db.Exec(queryChangePassword, newPassword, userID)
	return err
}
