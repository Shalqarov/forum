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
