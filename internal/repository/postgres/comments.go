package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryCreateComment = `
	INSERT INTO "comment"(
		"user_id",
		"post_id",
		"content",
		"date"
	) VALUES(?,?,?,?)`

	queryGetCommentsByPostID = `
	SELECT c.*, u.username, u.avatar FROM "comment" AS c
	INNER JOIN "user" AS u
	ON u.ID = c.user_id
	WHERE "post_id"=?
	ORDER BY "date" DESC`
	queryGetVotesCountByCommentID = `
	SELECT "vote", count("vote") 
	FROM "comment_votes"
	WHERE comment_id = $1 
	GROUP BY "vote"
	ORDER BY "vote" desc`
)

func NewSqliteCommentRepo(db *sql.DB) domain.CommentRepo {
	return &repo{db: db}
}

func (u *repo) CreateComment(comm *domain.Comment) error {
	_, err := u.db.Exec(queryCreateComment, comm.UserID, comm.PostID, comm.Content, time.Now())
	return err
}

func (u *repo) GetCommentsByPostID(id int64) ([]*domain.Comment, error) {
	rows, err := u.db.Query(queryGetCommentsByPostID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*domain.Comment{}
	for rows.Next() {
		date := time.Time{}
		comment := domain.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.PostID,
			&comment.Content,
			&date,
			&comment.Author,
			&comment.UserAvatar,
		)
		if err != nil {
			return nil, err
		}
		comment.Date = date.Format(time.RFC822)
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (u *repo) GetVotesCountByCommentID(commentID int64) (*domain.Vote, error) {
	rows, err := u.db.Query(queryGetVotesCountByCommentID, commentID)
	if err != nil {
		return nil, err
	}
	votes := &domain.Vote{
		Like:    0,
		Dislike: 0,
	}
	for rows.Next() {
		var voteType int64
		var cnt int64
		err := rows.Scan(&voteType, &cnt)
		if err != nil {
			return nil, err
		}
		switch voteType {
		case 1:
			votes.Like = uint64(cnt)
		case -1:
			votes.Dislike = uint64(cnt)
		default:
			log.Println("Get Votes count bug:", voteType, cnt)
		}
	}
	return votes, nil
}
