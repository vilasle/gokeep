package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyID = errors.New("empty id")
	ErrNotFound = errors.New("entity not found")
	ErrEncryption = errors.New("encryption failed")
	
)
