package domain

type UserUsecase interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(user *User) (*User, error)
}

type UserRepo interface {
	CreateUser(user *User) error
	GetUserIDByUsername(username string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(user *User) (*User, error)
}
