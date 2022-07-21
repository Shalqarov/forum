package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryCreateUser = `
	INSERT INTO "user" (
		"username",
		"email",
		"password",
		"avatar"
	) VALUES ($1, $2, $3, $4) RETURNING "user_id"`

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

	queryGetVotedPostsByUserID = `
	SELECT p."post_id", p."user_id", p."title", p."category", p."date", u."username"
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p."user_id" = u."user_id"
	INNER JOIN "post_vote" AS v
		ON v."user_id" = u.user_id AND v."post_id" = p."post_id"
	WHERE u."user_id"=$1 AND v."vote"=1
	ORDER BY p."date" DESC
	`
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
		user.Email,
		user.Password,
		user.Avatar,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	fmt.Println(lastInsertId)
	return lastInsertId, nil
}

func (u *repo) GetVotedPostsByUserID(userID int64) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetVotedPostsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*domain.PostDTO{}
	for rows.Next() {
		date := time.Time{}
		post := &domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Category, &date, &post.Author)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = date.Format("01-02-2006 15:04:05")
		posts = append(posts, post)
	}
	fmt.Println(posts)
	return posts, nil
}

func (u *repo) ChangeAvatarByUserID(userID int64, image string) error {
	imagePath := ""
	err := u.db.QueryRow(
		`SELECT "avatar" 
		FROM "user"
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
