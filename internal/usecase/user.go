package usecase

import (
	"github.com/Shalqarov/forum/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo domain.UserRepo
}

func NewUserUsecase(userRepo domain.UserRepo) domain.UserUsecase {
	return &userUsecase{
		repo: userRepo,
	}
}

func (u *userUsecase) CreateUser(user *domain.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashedPassword)
	return u.repo.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id int64) (*domain.User, error) {
	user, err := u.repo.GetUserByID(id)
	return user, err
}

func (u *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	searchedUser, err := u.repo.GetUserByEmail(email)
	return searchedUser, err
}

func (u *userUsecase) GetUserIDByUsername(username string) (int64, error) {
	return u.repo.GetUserIDByUsername(username)
}

func (u *userUsecase) ChangeAvatarByUserID(userID int64, image string) error {
	return u.repo.ChangeAvatarByUserID(userID, image)
}

func (u *userUsecase) ChangePassword(newPassword string, userID int64) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return err
	}
	return u.repo.ChangePassword(string(hashedPassword), userID)
}
