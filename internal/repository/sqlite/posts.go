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
	SELECT p."id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.id
	WHERE u.id=?
	ORDER BY "date" DESC`

	queryGetPostByID = `
	SELECT p.*, u.username, u.avatar FROM "post" AS p
	INNER JOIN "user" AS u
	ON u.ID = p.user_id`

	queryGetPostsByCategory = `
	SELECT p."id", p."user_id", p."title", p."category", p."date", u.username
	FROM "post" AS p
	INNER JOIN "user" AS u
		ON p.user_id = u.id
	WHERE p.category=?
	ORDER BY "date" DESC`

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

func scanPostDTORows(rows *sql.Rows) ([]*domain.PostDTO, error) {
	posts := []*domain.PostDTO{}
	for rows.Next() {
		post := domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Category, &post.CreatedAt, &post.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	rows.Close()
	return posts, nil
}
