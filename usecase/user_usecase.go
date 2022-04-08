package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepo
}

func NewUserUsecase(userRepo domain.UserRepo) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Create(user *domain.User) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		return domain.ErrBadParamInput
	}

	return u.userRepo.Create(user)
}

func (u *userUsecase) GetByID(id int) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(user *domain.User) (*domain.User, error) {
	searchedUser, err := u.userRepo.GetByEmail(user)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if err = bcrypt.CompareHashAndPassword([]byte(searchedUser.Password), []byte(user.Password)); err != nil {
		return nil, domain.ErrBadParamInput
	}
	return nil, nil
}
