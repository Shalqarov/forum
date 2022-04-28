package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
)

type userUsecase struct {
	repo domain.UserRepo
}

func NewUserUsecase(userRepo domain.UserRepo) domain.UserUsecase {
	return &userUsecase{
		repo: userRepo,
	}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		return domain.ErrBadParamInput
	}

	return u.repo.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id int) (*domain.User, error) {
	user, err := u.repo.GetUserByID(id)
	return user, err
}

func (u *userUsecase) GetUserByEmail(user *domain.User) (*domain.User, error) {
	searchedUser, err := u.repo.GetUserByEmail(user)
	return searchedUser, err
}

func (u *userUsecase) GetUserIDByUsername(username string) (int, error) {
	return u.repo.GetUserIDByUsername(username)
}
