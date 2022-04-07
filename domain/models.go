package domain

import "context"

// User - struct have basic user fields
type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, user *User) (*User, error)
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, user *User) (*User, error)
}
