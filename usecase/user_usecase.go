package usecase

import (
	"context"

	"github.com/Shalqarov/forum/domain"
)

type userUsecase struct {
	userRepo domain.UserRepo
}

func NewUserUsecase(userRepo domain.UserRepo) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (u *userUsecase) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return nil, nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
