package domain

// User - struct have basic user fields
type User struct {
	ID       int
	Username string
	Email    string
	Password string
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
