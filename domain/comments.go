package domain

type CommentUsecase interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostTitle(title string) ([]*CommentDTO, error)
}
type CommentRepo interface {
	CreateComment(comm *Comment) error
	GetCommentsByPostTitle(title string) ([]*CommentDTO, error)
}
