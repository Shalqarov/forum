package domain

// User - struct have basic user fields
type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type Post struct {
	ID      int
	UserID  int
	Title   string
	Content string
}

type Usecase interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)

	CreatePost(post *Post) error
	GetPostByUserID(id int) (*Post, error)
	GetPostByTitle(title string) (*Post, error)
}

type Repo interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)

	CreatePost(post *Post) error
	GetPostByUserID(id int) (*Post, error)
	GetPostByTitle(title string) (*Post, error)
}
