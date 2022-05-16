package usecase

import (
	"fmt"

	"github.com/Shalqarov/forum/domain"
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
		return nil, fmt.Errorf("GetPostByID error: %w", err)
	}

	for _, comment := range comments {
		votes, err := u.GetVotesCountByCommentID(comment.ID)
		if err != nil {
			return nil, fmt.Errorf("GetVotesCountByCommentID error: %w", err)
		}
		comment.Votes = *votes
	}
	return comments, nil
}

func (u *commentUsecase) VoteComment(commentID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VoteComment: invalid voteType")
	}
	return u.repo.VoteComment(commentID, userID, vote)
}

func (u *commentUsecase) GetVotesCountByCommentID(commentID int64) (*domain.Vote, error) {
	return u.repo.GetVotesCountByCommentID(commentID)
}
