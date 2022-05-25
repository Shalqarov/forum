package sqlite

import (
	"database/sql"
	"log"

	"github.com/Shalqarov/forum/domain"
)

const (
	queryVotePostSelect = `
	SELECT "id","vote" 
	FROM "post_votes" 
	WHERE user_id = ? AND post_id = ?`
	queryVotePostExec = `
	INSERT INTO "post_votes"(
		"user_id",
		"post_id",
		"vote")
	VALUES (?,?,?)`
	queryVotePostDelete = `
	DELETE FROM "post_votes" 
	WHERE "id" = ?`

	queryGetVotedPostsByUserID = `
	SELECT p."id","author","title","category","date" 
	FROM "post" AS p
	INNER JOIN "post_votes" AS v ON p."id"=v."post_id"	
	WHERE p."user_id"=? AND v."vote"=1`

	queryGetVotesCountByPostID = `
	SELECT "vote", count("vote") 
	FROM "post_votes"
	WHERE post_id = ? 
	GROUP BY "vote"
	ORDER BY "vote" desc`
)

func (u *sqliteRepo) VotePost(postID, userID int64, vote int) error {
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

func (u *sqliteRepo) GetVotedPostsByUserID(userID int64) ([]*domain.PostDTO, error) {
	rows, err := u.db.Query(queryGetVotedPostsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) GetVotesCountByPostID(postID int64) (*domain.Vote, error) {
	rows, err := u.db.Query(queryGetVotesCountByPostID, postID)
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

func scanPostDTORows(rows *sql.Rows) ([]*domain.PostDTO, error) {
	posts := []*domain.PostDTO{}
	for rows.Next() {
		post := domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Category, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}
