package domain

type CommentUsecase interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int) ([]*Comment, error)
}
type CommentRepo interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostID(id int) ([]*Comment, error)
}
