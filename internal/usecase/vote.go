package usecase

import (
	"fmt"

	"github.com/Shalqarov/forum/internal/domain"
)

type voteUsecase struct {
	repo domain.VoteRepo
}

func NewVoteUsecase(voteRepo domain.VoteRepo) domain.VoteUsecase {
	return &voteUsecase{
		repo: voteRepo,
	}
}

func (u *voteUsecase) VotePost(postID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VotePost: invalid voteType")
	}
	return u.repo.VotePost(postID, userID, vote)
}

func (u *voteUsecase) VoteComment(commentID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VoteComment: invalid voteType")
	}
	return u.repo.VoteComment(commentID, userID, vote)
}
