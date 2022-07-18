package domain

type VoteUsecase interface {
	VotePost(postID, userID int64, vote int) error
	VoteComment(commentID, userID int64, vote int) error
}

type VoteRepo interface {
	VotePost(postID, userID int64, vote int) error
	VoteComment(commentID, userID int64, vote int) error
}
