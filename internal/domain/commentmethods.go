package domain

type CommentUsecase interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int64) ([]*Comment, error)
	GetVotesCountByCommentID(commentID int64) (*Vote, error)
}
type CommentRepo interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int64) ([]*Comment, error)
	GetVotesCountByCommentID(commentID int64) (*Vote, error)
}
