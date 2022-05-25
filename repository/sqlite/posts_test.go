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

var post = &domain.Post{
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

func TestGetPostsByUserID(t *testing.T) {
	db, mock := NewMock()
	repo := NewSqlitePostRepo(db)
	query := queryGetPostsByUserID
	rows := sqlmock.NewRows([]string{"id", "title", "category", "date"}).AddRow(post.ID, post.Title, post.Category, time.Now().Format(time.RFC822))

	mock.ExpectQuery(query).WithArgs(post.UserID).WillReturnRows(rows)

	posts, err := repo.GetPostsByUserID(post.UserID)
	assert.NotNil(t, posts)
	assert.NoError(t, err)
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
	query := queryCreatePost

	mock.ExpectExec(query).WithArgs(post.UserID, post.Author, post.Title, post.Content, post.Category, time.Now().Format(time.RFC822)).WillReturnResult(sqlmock.NewResult(111, 1))

	id, err := repo.CreatePost(post)
	assert.NoError(t, err)
	assert.Equal(t, int64(111), id)
}
