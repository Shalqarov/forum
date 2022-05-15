package domain

type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}

type Vote struct {
	Like    uint64
	Dislike uint64
}

type Post struct {
	ID        int64
	UserID    int64
	Author    string
	Title     string
	Content   string
	Category  string
	CreatedAt string
	Votes     Vote
}

type Comment struct {
	ID      int64
	UserID  int64
	PostID  int64
	Author  string
	Content string
	Date    string
	Votes   Vote
}

type PostDTO struct {
	ID        int64
	Author    string
	Title     string
	Category  string
	CreatedAt string
}
type CommentDTO struct {
	Author  string
	Content string
	Date    string
}
