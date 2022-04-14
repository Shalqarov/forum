package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
)

func (u *usecase) CreatePost(post *domain.Post) error {
	if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Content) == "" {
		return domain.ErrBadParamInput
	}
	return u.repo.CreatePost(post)
}

func (u *usecase) GetPostByUserID(id int) (*domain.Post, error) {
	return nil, nil
}

func (u *usecase) GetPostByTitle(title string) (*domain.Post, error) {
	return nil, nil
}

func (u *usecase) GetPostsByCategory(category string) ([]*domain.Post, error) {
	return nil, nil
}
