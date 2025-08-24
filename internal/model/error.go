package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyID = errors.New("empty id")
)
