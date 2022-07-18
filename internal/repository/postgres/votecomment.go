package postgres

import (
	"database/sql"
	"log"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryVoteCommentSelect = `
	SELECT "id","vote" 
	FROM "comment_votes" 
	WHERE user_id = $1 AND comment_id = $2`
	queryVoteCommentExec = `
	INSERT INTO "comment_votes"(
		"user_id",
		"comment_id",
		"vote")
	VALUES ($1, $2, $3)`
	queryVoteCommentDelete = `
	DELETE FROM "comment_votes" 
	WHERE "id" = $1`

	queryGetVotesCountByCommentID = `
	SELECT "vote", count("vote") 
	FROM "comment_votes"
	WHERE comment_id = $1 
	GROUP BY "vote"
	ORDER BY "vote" desc`
)

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
