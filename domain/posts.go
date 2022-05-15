package domain

type PostUsecase interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int) ([]*PostDTO, error)
	GetPostByID(id int) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
	// VotePost()
}

type PostRepo interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int) ([]*PostDTO, error)
	GetPostByID(id int) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
}
