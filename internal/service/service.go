package service

import "context"

// AuthService - create new users and check existing users, and check tokens
type AuthService interface {
	//Login check user's login and password and return credential token
	Login(ctx context.Context, username, password string) (string, error)
	//Register create user if it does not exist	and return credential token
	Register(ctx context.Context, username, password string) (string, error)
	Valid(ctx context.Context, token string) error
}
