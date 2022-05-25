package sqlite

import (
	"database/sql"
	"time"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryCreateComment = `
	INSERT INTO "comment"(
		"user_id",
		"post_id",
		"author",
		"content",
		"date"
	) VALUES(?,?,?,?,?)`

	queryGetCommentsByPostID = `
	SELECT * 
	FROM "comment" 
	WHERE "post_id"=? 
	ORDER BY "date" DESC`
)

func NewSqliteCommentRepo(db *sql.DB) domain.CommentRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreateComment(comm *domain.Comment) error {
	_, err := u.db.Exec(queryCreateComment, comm.UserID, comm.PostID, comm.Author, comm.Content, time.Now().Format(time.RFC822))
	return err
}

func (u *sqliteRepo) GetCommentsByPostID(id int64) ([]*domain.Comment, error) {
	rows, err := u.db.Query(queryGetCommentsByPostID, id)
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
