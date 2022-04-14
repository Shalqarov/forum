package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
	"golang.org/x/crypto/bcrypt"
)

func (u *usecase) CreateUser(user *domain.User) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.Email) == "" {
		return domain.ErrBadParamInput
	}

	return u.repo.CreateUser(user)
}

func (u *usecase) GetUserByID(id int) (*domain.User, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (u *usecase) GetUserByEmail(user *domain.User) (*domain.User, error) {
	searchedUser, err := u.repo.GetUserByEmail(user)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if err = bcrypt.CompareHashAndPassword([]byte(searchedUser.Password), []byte(user.Password)); err != nil {
		return nil, domain.ErrBadParamInput
	}
	return searchedUser, nil
}

func (u *usecase) GetUserIDByUsername(username string) (int, error) {
	return u.repo.GetUserIDByUsername(username)
}
