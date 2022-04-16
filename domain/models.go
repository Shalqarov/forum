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
	Title     string
	Content   string
	Category  string
	CreatedAt string
}

type Usecase interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)

	CreatePost(post *Post) error
	GetPostsByUserID(id int) ([]*Post, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
}

type Repo interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)

	CreatePost(post *Post) error
	GetPostsByUserID(id int) ([]*Post, error)
	GetPostByTitle(title string) (*Post, error)
	GetPostsByCategory(category string) ([]*Post, error)
}
