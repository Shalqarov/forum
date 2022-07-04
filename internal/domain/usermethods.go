package domain

type UserUsecase interface {
	CreateUser(user *User) (int64, error)
	GetUserIDByUsername(username string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ChangeAvatarByUserID(userID int64, image string) error
	ChangePassword(newPassword string, userID int64) error
}

type UserRepo interface {
	CreateUser(user *User) (int64, error)
	GetUserIDByUsername(username string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ChangeAvatarByUserID(userID int64, image string) error
	ChangePassword(newPassword string, userID int64) error
}
