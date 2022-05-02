package usecase

import "github.com/Shalqarov/forum/domain"

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

func (u *commentUsecase) GetCommentsByPostTitle(title string) ([]*domain.CommentDTO, error) {
	return u.repo.GetCommentsByPostTitle(title)
}
