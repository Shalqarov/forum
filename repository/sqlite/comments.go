package sqlite

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/domain"
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

func (u *sqliteRepo) GetCommentsByPostID(id int64) ([]*domain.Comment, error) {
	stmt := `SELECT * FROM "comment" WHERE "post_id"=? ORDER BY "date" DESC`
	rows, err := u.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*domain.Comment{}
	for rows.Next() {
		comment := domain.Comment{}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Author, &comment.Content, &comment.Date)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
