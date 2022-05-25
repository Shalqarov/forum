package usecase

import (
	"fmt"

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

func (u *postUsecase) CreatePost(post *domain.Post) (int64, error) {
	return u.repo.CreatePost(post)
}

func (u *postUsecase) GetPostsByUserID(id int64) ([]*domain.PostDTO, error) {
	return u.repo.GetPostsByUserID(id)
}

func (u *postUsecase) GetPostByID(id int64) (*domain.Post, error) {
	post, err := u.repo.GetPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetPostByID error: %w", err)
	}
	votes, err := u.GetVotesCountByPostID(id)
	if err != nil {
		return nil, fmt.Errorf("GetVotesCountByPostID error: %w", err)
	}
	post.Votes = *votes
	return post, nil
}

func (u *postUsecase) GetPostsByCategory(category string) ([]*domain.PostDTO, error) {
	if !containsCategory(category) {
		return nil, fmt.Errorf("postCategory: entered category doesn't exists")
	}
	return u.repo.GetPostsByCategory(category)
}

var categories = map[string]bool{
	"alem":     true,
	"Study":    true,
	"Teamalem": true,
	"Linkedin": true,
	"Offtop":   true,
}

func containsCategory(category string) bool {
	_, ok := categories[category]
	return ok
}

func (u *postUsecase) GetAllPosts() ([]*domain.PostDTO, error) {
	return u.repo.GetAllPosts()
}

func (u *postUsecase) VotePost(postID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VotePost: invalid voteType")
	}
	return u.repo.VotePost(postID, userID, vote)
}

func (u *postUsecase) GetVotedPostsByUserID(userID int64) ([]*domain.PostDTO, error) {
	return u.repo.GetVotedPostsByUserID(userID)
}

func (u *postUsecase) GetVotesCountByPostID(postID int64) (*domain.Vote, error) {
	return u.repo.GetVotesCountByPostID(postID)
}
