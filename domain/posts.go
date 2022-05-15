package domain

type PostUsecase interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int64) ([]*PostDTO, error)
	GetPostByID(id int64) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
	VotePost(postID, userID int64, vote int) error
}

type PostRepo interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int64) ([]*PostDTO, error)
	GetPostByID(id int64) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
	GetAllPosts() ([]*PostDTO, error)
	VotePost(postID, userID int64, vote int) error
}
