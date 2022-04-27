package usecase

import (
	"strings"

	"github.com/Shalqarov/forum/domain"
)

type postUsecase struct {
	repo domain.PostRepo
}

func NewPostUsecase(postRepo domain.PostRepo) domain.PostUsecase {
	return &postUsecase{
		repo: postRepo,
	}
}

func (u *postUsecase) CreatePost(post *domain.Post) error {
	if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Content) == "" {
		return domain.ErrBadParamInput
	}
	return u.repo.CreatePost(post)
}

func (u *postUsecase) GetPostsByUserID(id int) ([]*domain.Post, error) {
	return u.repo.GetPostsByUserID(id)
}

func (u *postUsecase) GetPostByTitle(title string) (*domain.Post, error) {
	return nil, nil
}

func (u *postUsecase) GetPostsByCategory(category string) ([]*domain.Post, error) {
	return nil, nil
}

func (u *postUsecase) GetAllPosts() ([]*domain.Post, error) {
	return u.repo.GetAllPosts()
}
