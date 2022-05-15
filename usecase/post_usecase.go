package usecase

import (
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
	return u.repo.CreatePost(post)
}

func (u *postUsecase) GetAllPostsByUserID(id int64) ([]*domain.PostDTO, error) {
	return u.repo.GetAllPostsByUserID(id)
}

func (u *postUsecase) GetPostByID(id int64) (*domain.Post, error) {
	return u.repo.GetPostByID(id)
}

func (u *postUsecase) GetPostsByCategory(category string) ([]*domain.Post, error) {
	return nil, nil
}

func (u *postUsecase) GetAllPosts() ([]*domain.PostDTO, error) {
	return u.repo.GetAllPosts()
}

func (u *postUsecase) VotePost(postID, userID int64, vote int) error {
	return nil
}
