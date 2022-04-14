package repository

import (
	"database/sql"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

type sqliteRepo struct {
	db *sql.DB
}

func NewSqliteRepo(db *sql.DB) domain.Repo {
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
	if err != nil {
		return domain.ErrConflict
	}

	return nil
}

func (u *sqliteRepo) GetUserByID(id int) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "id"=?`
	user := domain.User{}
	err := u.db.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &user, nil
}

func (u *sqliteRepo) GetUserByEmail(user *domain.User) (*domain.User, error) {
	stmt := `SELECT * FROM "user" WHERE "email"=?`
	searchedUser := domain.User{}
	err := u.db.QueryRow(stmt, user.Email).Scan(&searchedUser.ID, &searchedUser.Username, &searchedUser.Email, &searchedUser.Password)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &searchedUser, nil
}

func (u *sqliteRepo) CreatePost(post *domain.Post) error {
	stmt := `INSERT INTO "post"(
		"user_id",
		"title",
		"content"
		) VALUES(?,?,?)`
	_, err := u.db.Exec(stmt, post.UserID, post.Title, post.Content)
	if err != nil {
		return domain.ErrConflict
	}
	return nil
}

func (u *sqliteRepo) GetPostByUserID(id int) (*domain.Post, error) {
	stmt := `SELECT * FROM "post" WHERE "user_id" = ?`
	post := domain.Post{}
	err := u.db.QueryRow(stmt, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &post, nil
}

func (u *sqliteRepo) GetPostByTitle(title string) (*domain.Post, error) {
	stmt := `SELECT * FROM "post" WHERE "title" = ?`
	post := domain.Post{}
	err := u.db.QueryRow(stmt, title).Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &post, nil
}
