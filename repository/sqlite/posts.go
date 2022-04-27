package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Shalqarov/forum/domain"
)

func NewSqlitePostRepo(db *sql.DB) domain.PostRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreatePost(post *domain.Post) error {
	stmt := `INSERT INTO "post"(
		"user_id",
		"author",
		"title",
		"content",
		"category",
		"date"
		) VALUES(?,?,?,?,?,?)`
	_, err := u.db.Exec(stmt, post.UserID, post.Author, post.Title, post.Content, post.Category, time.Now().Format(time.RFC822))
	if err != nil {
		return domain.ErrConflict
	}
	return nil
}

func (u *sqliteRepo) GetPostsByUserID(id int) ([]*domain.Post, error) {
	stmt := `SELECT "title","content","category","date" FROM "post" WHERE "user_id" = ? ORDER BY "date" DESC`
	rows, err := u.db.Query(stmt, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	defer rows.Close()

	posts := []*domain.Post{}

	for rows.Next() {
		post := domain.Post{}
		err = rows.Scan(&post.Title, &post.Content, &post.Category, &post.CreatedAt)
		fmt.Println(post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (u *sqliteRepo) GetPostByTitle(title string) (*domain.Post, error) {
	stmt := `SELECT * FROM "post" WHERE "title" = ?`
	post := domain.Post{}
	err := u.db.QueryRow(stmt, title).Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &post, nil
}

func (u *sqliteRepo) GetPostsByCategory(category string) ([]*domain.Post, error) {
	return nil, nil
}

func (u *sqliteRepo) GetAllPosts() ([]*domain.Post, error) {
	stmt := `SELECT * FROM "post" ORDER BY "date" DESC`
	rows, err := u.db.Query(stmt)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	defer rows.Close()

	posts := []*domain.Post{}

	for rows.Next() {
		post := domain.Post{}
		err = rows.Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
		fmt.Println(post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}
