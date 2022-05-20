package domain

type PostUsecase interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int64) ([]*PostDTO, error)
	GetPostByID(id int64) (*Post, error)
	GetPostsByCategory(category string) ([]*PostDTO, error)
	GetAllPosts() ([]*PostDTO, error)
	VotePost(postID, userID int64, vote int) error
	GetVotedPostsByUserID(userID int64) ([]*PostDTO, error)
	GetVotesCountByPostID(postID int64) (*Vote, error)
}

type PostRepo interface {
	CreatePost(post *Post) error
	GetAllPostsByUserID(id int64) ([]*PostDTO, error)
	GetPostByID(id int64) (*Post, error)
	GetPostsByCategory(category string) ([]*PostDTO, error)
	GetAllPosts() ([]*PostDTO, error)
	VotePost(postID, userID int64, vote int) error
	GetVotedPostsByUserID(userID int64) ([]*PostDTO, error)
	GetVotesCountByPostID(postID int64) (*Vote, error)
}
