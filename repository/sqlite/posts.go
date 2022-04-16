package repository

import "github.com/Shalqarov/forum/domain"

func (u *sqliteRepo) CreatePost(post *domain.Post) error {
	stmt := `INSERT INTO "post"(
		"user_id",
		"title",
		"content",
		"category"
		) VALUES(?,?,?,?)`
	_, err := u.db.Exec(stmt, post.UserID, post.Title, post.Content, post.Category)
	if err != nil {
		return domain.ErrConflict
	}
	return nil
}

func (u *sqliteRepo) GetPostsByUserID(id int) ([]*domain.Post, error) {
	stmt := `SELECT * FROM "post" WHERE "user_id" = ?`
	rows, err := u.db.Query(stmt)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	defer rows.Close()
	post := domain.Post{}
	posts := []domain.Post{}

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category)
		if err != nil {
		}

	}

	return nil, nil
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

func (u *sqliteRepo) GetPostsByCategory(category string) ([]*domain.Post, error) {
	return nil, nil
}
