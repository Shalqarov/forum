package domain

// User - struct have basic user fields
type User struct {
	ID       int
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserUsecase interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(user *User) (*User, error)
}

type UserRepo interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(user *User) (*User, error)
}
