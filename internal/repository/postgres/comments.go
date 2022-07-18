package postgres

import (
	"log"

	"github.com/Shalqarov/forum/internal/domain"
)

const (
	queryGetVotesCountByCommentID = `
		SELECT "vote", count("vote") 
		FROM "comment_votes"
		WHERE comment_id = $1 
		GROUP BY "vote"
		ORDER BY "vote" desc`
)

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
