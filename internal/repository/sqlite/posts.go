package sqlite

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/internal/domain"
)

func NewSqlitePostRepo(db *sql.DB) domain.PostRepo {
	return &sqliteRepo{db: db}
}

const (
	queryCreatePost = `
	INSERT INTO "post"(
	"user_id",
	"title",
	"content",
	"category",
	"date",
	"image"
	) VALUES(?,?,?,?,?,?)`

	queryGetPostsByUserID = `
	SELECT "id","user_id","title","category","date" 
	FROM "post" 
	WHERE "user_id" = ? 
	ORDER BY "id" DESC`

	queryGetPostByID = `
	SELECT p.*, u.username, u.avatar FROM "post" AS p
	INNER JOIN "user" AS u
	ON u.ID = p.user_id`

	queryGetPostsByCategory = `
	SELECT "id","user_id","author","title","category","date" 
	FROM "post" 
	WHERE "category"=? 
	ORDER BY "id" DESC`

	queryGetAllPosts = `
	SELECT p."id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.id
	ORDER BY "date" DESC`
)

func (u *sqliteRepo) CreatePost(post *domain.Post) (int64, error) {
	result, err := u.db.Exec(
		queryCreatePost,
		post.UserID,
		post.Title,
		post.Content,
		post.Category,
		time.Now().Format(time.RFC822),
		post.Image,
	)
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
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) GetPostByID(id int64) (*domain.Post, error) {
	post := domain.Post{}
	err := u.db.QueryRow(queryGetPostByID, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Category,
		&post.CreatedAt,
		&post.Image,
		&post.Author,
		&post.UserAvatar,
	)
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
