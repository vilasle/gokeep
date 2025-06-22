package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserNotExists = errors.New("user not exists")
)
