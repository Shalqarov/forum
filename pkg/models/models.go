package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

// User - struct have basic user fields
type User struct {
	ID       int
	Login    string
	Email    string
	Password string
}
