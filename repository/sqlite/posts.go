package sqlite

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/domain"
)

func NewSqlitePostRepo(db *sql.DB) domain.PostRepo {
	return &sqliteRepo{db: db}
}

const (
	queryCreatePost = `
	INSERT INTO "post"(
	"user_id",
	"author",
	"title",
	"content",
	"category",
	"date"
	) VALUES(?,?,?,?,?,?)`

	queryGetPostsByUserID = `
	SELECT "id","title","category","date" 
	FROM "post" 
	WHERE "user_id" = ? 
	ORDER BY "id" DESC`

	queryGetPostByID = `
	SELECT * 
	FROM "post" 
	WHERE "id" = ?`

	queryGetPostsByCategory = `
	SELECT "id","author","title","category","date" 
	FROM "post" 
	WHERE "category"=? 
	ORDER BY "id" DESC`

	queryGetAllPosts = `
	SELECT "id","author","title","category","date" 
	FROM "post" 
	ORDER BY "date" DESC`
)

func (u *sqliteRepo) CreatePost(post *domain.Post) (int64, error) {
	result, err := u.db.Exec(queryCreatePost, post.UserID, post.Author, post.Title, post.Content, post.Category, time.Now().Format(time.RFC822))
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (u *sqliteRepo) GetPostsByUserID(id int64) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetPostsByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*domain.PostDTO{}
	for rows.Next() {
		post := domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.Title, &post.Category, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (u *sqliteRepo) GetPostByID(id int64) (*domain.Post, error) {
	post := domain.Post{}
	err := u.db.QueryRow(queryGetPostByID, id).Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (u *sqliteRepo) GetPostsByCategory(category string) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetPostsByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) GetAllPosts() ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}
