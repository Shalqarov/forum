package sqlite

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Shalqarov/forum/domain"
	"github.com/stretchr/testify/assert"
)

var p = &domain.Post{
	ID:        111,
	UserID:    222,
	Author:    "mangomango",
	Title:     "Title Test 1",
	Content:   " Content TESTTEST TESTTEST TESTTEST TESTTEST TESTTEST TESTTEST TESTTEST TESTTEST",
	Category:  "Category Test alem",
	CreatedAt: "01-02-2022",
	Votes: domain.Vote{
		Like:    50,
		Dislike: 50,
	},
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreatePost(t *testing.T) {
	db, mock := NewMock()
	repo := NewSqlitePostRepo(db)
	defer func() {
		db.Close()
	}()
	query := queryCreatePost

	mock.ExpectExec(query).WithArgs(p.UserID, p.Author, p.Title, p.Content, p.Category, time.Now().Format(time.RFC822)).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CreatePost(p)
	assert.NoError(t, err)
}
