package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.AuthService = (*AuthService)(nil)

// AuthService - creates new users and check existing users, creates and checks JWT tokens
type AuthService struct {
	jwtKey  []byte
	manager model.ModelManager
}

// NewAuthService creates new AuthService
func NewAuthService(modelManager model.ModelManager, jwtKey []byte) *AuthService {
	return &AuthService{manager: modelManager, jwtKey: jwtKey}
}

// Login check user's login and password and return credential token
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.manager.Users.FindByLogin(ctx, username)
	if err != nil {
		//TODO repository error
		return "", err
	}

	if !u.PasswordIsValid(password) {
		//TODO wrap error
		return "", errors.New("invalid password")
	}

	return s.createToken(u)

}

// Register create user if it does not exist	and return credential token
func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if _, err := s.manager.Users.FindByLogin(ctx, username); !errors.Is(err, model.ErrUserNotFound) {
		if err != nil {
			//TODO repository error
			return "", err
		}
		return "", errors.New("user already exists")
	}

	u := s.manager.Users.New(username, password)
	if err := u.Save(ctx); err != nil {
		//TODO storage error
		return "", err
	}
	return s.createToken(u)
}

// Valid parse token, and try to find user on database, if it is found return nil, else return error
func (s *AuthService) Valid(ctx context.Context, token string) error {
	return s.checkToken(ctx, token)
}

// createToken - create token for user
func (s *AuthService) createToken(user *model.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID(),
	})
	token, err := t.SignedString(s.jwtKey)
	if err != nil {
		//TODO create token error, wrap error
		return "", err
	}

	return token, nil
}

// checkToken - check token and try to get user by id from database
func (s *AuthService) checkToken(ctx context.Context, token string) error {
	decToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		//TODO wrap error
		return err
	}

	var id int64
	if claims, ok := decToken.Claims.(jwt.MapClaims); ok {
		fid, ok := claims["id"].(float64)
		if !ok {
			//TODO wrap error
			return errors.New("invalid token")
		}
		id = int64(fid)
	} else {
		//TODO wrap error
		return errors.New("invalid token")
	}

	if _, err := s.manager.Users.GetByID(ctx, id); err != nil {
		//user not found or we get storage error
		//TODO wrap error
		return err
	}
	return nil
}
