package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Shalqarov/forum/domain"
)

func NewSqlitePostRepo(db *sql.DB) domain.PostRepo {
	return &sqliteRepo{db: db}
}

func (u *sqliteRepo) CreatePost(post *domain.Post) error {
	stmt := `INSERT INTO "post"(
		"user_id",
		"author",
		"title",
		"content",
		"category",
		"date"
		) VALUES(?,?,?,?,?,?)`
	_, err := u.db.Exec(stmt, post.UserID, post.Author, post.Title, post.Content, post.Category, time.Now().Format(time.RFC822))
	fmt.Println(err)
	return fmt.Errorf("Ошабка")
}

func (u *sqliteRepo) GetPostsByUserID(id int64) ([]*domain.PostDTO, error) {
	stmt := `SELECT "id","title","category","date" FROM "post" WHERE "user_id" = ? ORDER BY "id" DESC`
	rows, err := u.db.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []*domain.PostDTO{}
	for rows.Next() {
		post := domain.PostDTO{}
		err := rows.Scan(&post.ID, &post.Title, &post.Category, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (u *sqliteRepo) GetPostByID(id int64) (*domain.Post, error) {
	stmt := `SELECT * FROM "post" WHERE "id" = ?`
	post := domain.Post{}
	err := u.db.QueryRow(stmt, id).Scan(&post.ID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (u *sqliteRepo) GetPostsByCategory(category string) ([]*domain.PostDTO, error) {
	stmt := `SELECT "id","author","title","category","date" FROM "post" WHERE "category"=? ORDER BY "id" DESC`
	rows, err := u.db.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) GetAllPosts() ([]*domain.PostDTO, error) {
	stmt := `SELECT "id","author","title","category","date" FROM "post" ORDER BY "date" DESC`
	rows, err := u.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) VotePost(postID, userID int64, vote int) error {
	stmtSelect := `SELECT "id","vote" FROM "post_votes" WHERE user_id = ? AND post_id = ?`
	stmtExec := `INSERT INTO "post_votes"(
					"user_id",
					"post_id",
					"vote")
					VALUES (?,?,?)`
	stmtDelete := `DELETE FROM "post_votes" WHERE "id" = ?`
	var voteID int64
	var voteInDB int
	row := u.db.QueryRow(stmtSelect, userID, postID)
	err := row.Scan(&voteID, &voteInDB)
	if err == sql.ErrNoRows {
		_, err := u.db.Exec(stmtExec, userID, postID, vote)
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
	_, err = u.db.Exec(stmtExec, userID, postID, vote)
	if err != nil {
		return err
	}
	return nil
}

func (u *sqliteRepo) GetVotedPostsByUserID(userID int64) ([]*domain.PostDTO, error) {
	stmt := `SELECT p."id","author","title","category","date" 
				FROM "post" AS p
				INNER JOIN "post_votes" AS v ON p."id"=v."post_id"	
				WHERE p."user_id"=? AND v."vote"=1`
	rows, err := u.db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPostDTORows(rows)
}

func (u *sqliteRepo) GetVotesCountByPostID(postID int64) (*domain.Vote, error) {
	stmt := `SELECT "vote", count("vote") FROM "post_votes"
			WHERE post_id = ? 
			GROUP BY "vote"
			ORDER BY "vote" desc`
	rows, err := u.db.Query(stmt, postID)
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
