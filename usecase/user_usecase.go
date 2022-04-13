package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.Repo
}

func NewUserUsecase(userRepo domain.Repo) domain.Usecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		return domain.ErrBadParamInput
	}

	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id int) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (u *userUsecase) GetUserByEmail(user *domain.User) (*domain.User, error) {
	searchedUser, err := u.userRepo.GetUserByEmail(user)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if err = bcrypt.CompareHashAndPassword([]byte(searchedUser.Password), []byte(user.Password)); err != nil {
		return nil, domain.ErrBadParamInput
	}
	return searchedUser, nil
}

func (u *userUsecase) CreatePost(post *domain.Post) error {
	return nil
}

func (u *userUsecase) GetPostByID(id int) (*domain.Post, error) {
	return nil, nil
}

func (u *userUsecase) GetPostByTitle(title string) (*domain.Post, error) {
	return nil, nil
}
