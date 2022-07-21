package postgres

import (
	"database/sql"

	"github.com/Shalqarov/forum/internal/domain"
)

func NewPostgresVoteRepo(db *sql.DB) domain.VoteRepo {
	return &repo{db: db}
}

const (
	queryVotePostSelect = `
	SELECT "post_vote_id", "vote" 
	FROM "post_vote" 
	WHERE user_id = $1 AND post_id = $2`
	queryVotePostExec = `
	INSERT INTO "post_vote"(
		"user_id",
		"post_id",
		"vote")
	VALUES ($1, $2, $3) RETURNING "post_vote_id"`
	queryVotePostDelete = `
	DELETE FROM "post_vote" 
	WHERE "post_vote_id" = $1`
	queryVoteCommentSelect = `
	SELECT "comment_vote_id", "vote" 
	FROM "comment_vote" 
	WHERE user_id = $1 AND comment_id = $2`
	queryVoteCommentExec = `
	INSERT INTO "comment_vote"(
		"user_id",
		"comment_id",
		"vote")
	VALUES ($1, $2, $3) RETURNING "comment_vote_id"`
	queryVoteCommentDelete = `
	DELETE FROM "comment_vote" 
	WHERE "comment_vote_id" = $1`
)

func (u *repo) VotePost(postID, userID int64, vote int) error {
	var voteID int64
	var voteInDB int
	row := u.db.QueryRow(queryVotePostSelect, userID, postID)
	err := row.Scan(&voteID, &voteInDB)
	if err == sql.ErrNoRows {
		_, err := u.db.Exec(queryVotePostExec, userID, postID, vote)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = u.db.Exec(queryVotePostDelete, voteID)
	if err != nil {
		return err
	}
	if vote == voteInDB {
		return nil
	}
	_, err = u.db.Exec(queryVotePostExec, userID, postID, vote)
	if err != nil {
		return err
	}
	return nil
}

func (u *repo) VoteComment(commentID, userID int64, vote int) error {
	var voteID int64
	var voteInDB int
	row := u.db.QueryRow(queryVoteCommentSelect, userID, commentID)
	err := row.Scan(&voteID, &voteInDB)
	if err == sql.ErrNoRows {
		_, err := u.db.Exec(queryVoteCommentExec, userID, commentID, vote)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = u.db.Exec(queryVoteCommentDelete, voteID)
	if err != nil {
		return err
	}
	if vote == voteInDB {
		return nil
	}
	_, err = u.db.Exec(queryVoteCommentExec, userID, commentID, vote)
	if err != nil {
		return err
	}
	return nil
}
