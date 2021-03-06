package usecase

import (
	"database/sql"
	"fmt"

	"github.com/Shalqarov/forum/internal/domain"
)

type commentUsecase struct {
	repo domain.CommentRepo
}

func NewCommentUsecase(commRepo domain.CommentRepo) domain.CommentUsecase {
	return &commentUsecase{
		repo: commRepo,
	}
}

func (u *commentUsecase) CreateComment(comm *domain.Comment) error {
	return u.repo.CreateComment(comm)
}

func (u *commentUsecase) GetCommentsByPostID(id int64) ([]*domain.Comment, error) {
	comments, err := u.repo.GetCommentsByPostID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("GetCommentsByPostID: %w", err)
	}

	for _, comment := range comments {
		votes, err := u.GetVotesCountByCommentID(comment.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}
			return nil, fmt.Errorf("GetVotesCountByCommentID: %w", err)
		}
		comment.Votes = *votes
	}
	return comments, nil
}

func (u *commentUsecase) GetVotesCountByCommentID(commentID int64) (*domain.Vote, error) {
	return u.repo.GetVotesCountByCommentID(commentID)
}
