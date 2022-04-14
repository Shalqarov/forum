package usecase

import "github.com/Shalqarov/forum/domain"

type usecase struct {
	repo domain.Repo
}

func NewUsecase(userRepo domain.Repo) domain.Usecase {
	return &usecase{
		repo: userRepo,
	}
}
