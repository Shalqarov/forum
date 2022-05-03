package repository

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/tree/master/domain"
)

func NewSqliteCommentRepo(db *sql.DB) domain.CommentRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreateComment(comm *domain.Comment) error {
	stmt := `INSERT INTO "comment"(
		"user_id",
		"post_id",
		"author",
		"content",
		"date"
		) VALUES(?,?,?,?,?)`
	_, err := u.db.Exec(stmt, comm.UserID, comm.PostID, comm.Author, comm.Content, time.Now().Format(time.RFC822))
	return err
}

func (u *sqliteRepo) GetCommentsByPostTitle(title string) ([]*domain.CommentDTO, error) {
	stmt := `SELECT "author","content","date" FROM "comment" as comm INNER JOIN "post" as p WHERE p.title=? ORDER BY "date" DESC`
	rows, err := u.db.Query(stmt, title)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	defer rows.Close()
	comments := []*domain.CommentDTO{}
	for rows.Next() {
		comment := domain.CommentDTO{}
		err := rows.Scan(&comment.Author, &comment.Content, &comment.Date)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
