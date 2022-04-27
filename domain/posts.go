package domain

type PostUsecase interface {
	CreatePost(post *Post) error
	GetPostsByUserID(id int) ([]*Post, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*Post, error)
}

type PostRepo interface {
	CreatePost(post *Post) error
	GetPostsByUserID(id int) ([]*Post, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*Post, error)
}
