package usecase

import (
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
	return u.repo.GetCommentsByPostID(id)
}
