package usecase

import "fmt"

func (u *postUsecase) VotePost(postID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VotePost: invalid voteType")
	}
	return u.repo.VotePost(postID, userID, vote)
}

func (u *commentUsecase) VoteComment(commentID, userID int64, vote int) error {
	if vote != -1 && vote != 1 {
		return fmt.Errorf("VoteComment: invalid voteType")
	}
	return u.repo.VoteComment(commentID, userID, vote)
}
