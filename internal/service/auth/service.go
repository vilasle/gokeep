package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/server"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.AuthService = (*AuthService)(nil)

// AuthService - creates new users and check existing users, creates and checks JWT tokens
type AuthService struct {
	jwtKey  []byte
	manager model.ModelManager
	session repository.SessionRepository
}

// NewAuthService creates new AuthService
func NewAuthService(modelManager model.ModelManager, session repository.SessionRepository, jwtKey []byte) *AuthService {
	return &AuthService{manager: modelManager, jwtKey: jwtKey, session: session}
}

// Login check user's login and password and return credential token
func (s *AuthService) Login(ctx context.Context, req service.RegisterLoginUser) (string, error) {
	u, err := s.manager.Users.FindByLogin(ctx, req.Username)
	if err != nil {
		//TODO repository error
		return "", err
	}

	if !u.PasswordIsValid(req.Password) {
		//TODO wrap error
		return "", errors.New("invalid password")
	}

	//create session
	sessionID, err := s.session.Create(ctx, repository.CredentialCreate{
		UserId:    u.ID(),
		PublicKey: req.PublicKey,
	})
	if err != nil {
		//TODO repository error
		return "", err
	}

	return s.createToken(u, sessionID)

}

// Register create user if it does not exist	and return credential token
func (s *AuthService) Register(ctx context.Context, req service.RegisterLoginUser) error {
	if _, err := s.manager.Users.FindByLogin(ctx, req.Username); !errors.Is(err, model.ErrUserNotFound) {
		if err != nil {
			//TODO repository error
			return err
		}
		return errors.New("user already exists")
	}

	u := s.manager.Users.New(req.Username, req.Password)
	if err := u.Save(ctx); err != nil {
		//TODO storage error
		return err
	}
	return nil
}

// GetSessionByCredentialToken parse token, and try to find user and session on database, if it is found return nil, else return error
func (s *AuthService) GetSessionByCredentialToken(ctx context.Context, token string) (service.SessionInfo, error) {
	return s.checkToken(ctx, token)
}

// createToken - create token for user
func (s *AuthService) createToken(user *model.User, sessionID int) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID(),
		"session_id": sessionID,
	})
	token, err := t.SignedString(s.jwtKey)
	if err != nil {
		//TODO create token error, wrap error
		return "", err
	}

	return token, nil
}

// checkToken - check token and try to get user by id from database
func (s *AuthService) checkToken(ctx context.Context, token string) (resp service.SessionInfo, err error) {
	decToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil {
		//TODO wrap error
		return resp, err
	}

	var userID, sessionID int
	if claims, ok := decToken.Claims.(jwt.MapClaims); ok {
		fid, ok := claims["id"].(float64)
		if !ok {
			//TODO wrap error
			return resp, errors.New("invalid token")
		}
		userID = int(fid)

		fses, ok := claims["session_id"].(float64)
		if !ok {
			//TODO wrap error
			return resp, errors.New("invalid token")
		}
		sessionID = int(fses)
	} else {
		//TODO wrap error
		return resp, errors.New("invalid token")
	}

	_, err = s.manager.Users.Get(ctx, int(userID))
	if err != nil {
		if err == model.ErrUserNotFound {
			return resp, service.ErrUserNotFound
		}
		return resp, err
	}

	cred, err := s.session.Get(ctx, sessionID)
	if err != nil {
		if err == repository.ErrNotFound {
			return resp, service.ErrSessionNotFound
		}
		return resp, err
	}

	if cred.UserId != userID {
		return resp, service.ErrSessionNotConnectedWithUser
	}

	resp.UserID = userID
	resp.SessionID = sessionID
	resp.PublicKey = cred.PublicKey

	return resp, nil
}
