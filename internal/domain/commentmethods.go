package domain

type CommentUsecase interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int64) ([]*Comment, error)
	VoteComment(commentID, userID int64, vote int) error
	GetVotesCountByCommentID(commentID int64) (*Vote, error)
}
type CommentRepo interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int64) ([]*Comment, error)
	VoteComment(commentID, userID int64, vote int) error
	GetVotesCountByCommentID(commentID int64) (*Vote, error)
}
