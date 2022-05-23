package sqlite

import (
	"database/sql"
	"log"

	"github.com/Shalqarov/forum/domain"
)

func (u *sqliteRepo) VoteComment(commentID, userID int64, vote int) error {
	stmtSelect := `SELECT "id","vote" FROM "comment_votes" WHERE user_id = ? AND comment_id = ?`
	stmtExec := `INSERT INTO "comment_votes"(
		"user_id",
		"comment_id",
		"vote")
		VALUES (?,?,?)`
	stmtDelete := `DELETE FROM "comment_votes" WHERE "id" = ?`
	var voteID int64
	var voteInDB int
	row := u.db.QueryRow(stmtSelect, userID, commentID)
	err := row.Scan(&voteID, &voteInDB)
	if err == sql.ErrNoRows {
		_, err := u.db.Exec(stmtExec, userID, commentID, vote)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = u.db.Exec(stmtDelete, voteID)
	if err != nil {
		return err
	}
	if vote == voteInDB {
		return nil
	}
	_, err = u.db.Exec(stmtExec, userID, commentID, vote)
	if err != nil {
		return err
	}
	return nil
}

func (u *sqliteRepo) GetVotesCountByCommentID(commentID int64) (*domain.Vote, error) {
	stmt := `SELECT "vote", count("vote") FROM "comment_votes"
		WHERE comment_id = ? 
		GROUP BY "vote"
		ORDER BY "vote" desc`
	rows, err := u.db.Query(stmt, commentID)
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
