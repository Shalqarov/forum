package domain

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type Post struct {
	ID        int
	UserID    int
	Author    string
	Title     string
	Content   string
	Category  string
	CreatedAt string
}

type Comment struct {
	ID      int
	UserID  int
	PostID  int
	Author  string
	Content string
	Date    string
}

type PostDTO struct {
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
