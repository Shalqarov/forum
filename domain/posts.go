package domain

type PostUsecase interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int) ([]*PostDTO, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
}

type PostRepo interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int) ([]*PostDTO, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
}
