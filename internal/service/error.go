package service

import "errors"

var ErrInvalidCredential = errors.New("invalid token")
var ErrSessionNotFound = errors.New("not found session")
var ErrUserNotFound = errors.New("not found user")
var ErrSessionNotConnectedWithUser = errors.New("session not connected with user")

