package domain

type UserUsecase interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(user *User) (*User, error)
}
