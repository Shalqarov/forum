package postgres

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/internal/domain"
)

func NewPostgresPostRepo(db *sql.DB) domain.PostRepo {
	return &repo{db: db}
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
	) VALUES($1, $2, $3, $4, $5, $6)`

	queryGetPostsByUserID = `
	SELECT p."post_id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.user_id
	WHERE u.user_id=$1
	ORDER BY "date" DESC`

	queryGetPostByID = `
	SELECT p.*, u.username, u.avatar FROM "post" AS p
	INNER JOIN "user" AS u
	ON u.user_id = p.user_id
	WHERE p.post_id=$1`

	queryGetPostsByCategory = `
	SELECT p."post_id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.user_id
	WHERE p.category=$1
	ORDER BY "date" DESC`

	queryGetAllPosts = `
	SELECT p."post_id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.user_id
	ORDER BY "date" DESC`
)

func (r *repo) CreatePost(post *domain.Post) (int64, error) {
	var lastInsertId int64
	if err := r.db.QueryRow(
		queryCreatePost,
		post.UserID,
		post.Title,
		post.Content,
		post.Category,
		time.Now(),
		post.Image,
	).Scan(&lastInsertId); err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (u *repo) GetPostsByUserID(id int64) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetPostsByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *repo) GetPostByID(id int64) (*domain.Post, error) {
	post := domain.Post{}
	time := time.Time{}
	err := u.db.QueryRow(queryGetPostByID, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Category,
		&time,
		&post.Image,
		&post.Author,
		&post.UserAvatar,
	)
	post.CreatedAt = time.Format("01-02-2006 15:04:05")
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (u *repo) GetPostsByCategory(category string) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetPostsByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *repo) GetAllPosts() ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func scanPostDTORows(rows *sql.Rows) ([]*domain.PostDTO, error) {
	posts := []*domain.PostDTO{}
	time := time.Time{}
	for rows.Next() {
		post := domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Category, &time, &post.Author)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = time.Format("01-02-2006 15:04:05")
		posts = append(posts, &post)
	}
	rows.Close()
	return posts, nil
}
