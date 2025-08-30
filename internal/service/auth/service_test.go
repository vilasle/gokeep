package auth

import (
	context "context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	model "github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/server"
	"github.com/vilasle/gokeep/internal/service"
)

func Test_AuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username: "test",
			Password: "password",
		}

		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{}, model.ErrUserNotFound)
		user.EXPECT().Add(ctx, model.UserAdd{
			Login:    req.Username,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		})

		err := svc.Register(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username: "test",
			Password: "password",
		}

		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{}, errors.New("repository error"))

		err := svc.Register(ctx, req)
		assert.Error(t, err)
	})

	t.Run("user is existed already", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username: "test",
			Password: "password",
		}

		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{}, nil)

		err := svc.Register(ctx, req)
		assert.Error(t, err)
	})
}

func Test_AuthService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username:  "test",
			Password:  "password",
			PublicKey: []byte("public key"),
		}
		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Create(ctx, repository.CredentialCreate{
			UserID:    1,
			PublicKey: []byte("public key"),
		}).Return(1, nil)

		token, err := svc.Login(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username:  "test",
			Password:  "password",
			PublicKey: []byte("public key"),
		}

		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{}, errors.New("user not found"))

		token, err := svc.Login(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("invalid password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username:  "test",
			Password:  "password",
			PublicKey: []byte("public key"),
		}
		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{
			ID:       1,
			Password: "5e88",
		}, nil)

		token, err := svc.Login(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("create session failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := service.RegisterLoginUser{
			Username:  "test",
			Password:  "password",
			PublicKey: []byte("public key"),
		}
		user.EXPECT().Find(ctx, req.Username).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Create(ctx, repository.CredentialCreate{
			UserID:    1,
			PublicKey: []byte("public key"),
		}).Return(0, errors.New("create session failed"))

		token, err := svc.Login(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func Test_AuthService_GetSessionByCredentialToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		sessionID := 1
		publicKey := []byte("public key")
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Get(ctx, sessionID).Return(repository.CredentialInfo{
			ID:        1,
			UserID:    1,
			PublicKey: publicKey,
		}, nil)

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, userID, info.UserID)
		assert.Equal(t, sessionID, info.SessionID)
		assert.Equal(t, publicKey, info.PublicKey)
	})

	t.Run("invalid secret token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("another secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, info)
	})

	t.Run("empty userID and sessionID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MCwic2Vzc2lvbklEIjowfQ.vecNKE02SXlt-cCUxi26kgZST2ylKPRAljzqU61wcTI"

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.ErrorIs(t, err, service.ErrInvalidToken)
		assert.Empty(t, info)
	})

	t.Run("user repository failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, errors.New("user repository failed"))

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, info)
	})

	t.Run("user not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, model.ErrUserNotFound)

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.ErrorIs(t, err, service.ErrUserNotFound)
		assert.Empty(t, info)
	})

	t.Run("getting session failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		sessionID := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Get(ctx, sessionID).Return(repository.CredentialInfo{}, errors.New("getting session failed"))

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, info)
	})

	t.Run("session not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		sessionID := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Get(ctx, sessionID).Return(repository.CredentialInfo{}, repository.ErrNotFound)

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.ErrorIs(t, err, service.ErrSessionNotFound)
		assert.Empty(t, info)
	})

	t.Run("session is not connected with user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := NewMockUserRepository(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(nil)

		manager := model.NewModelManager(collector)

		session := NewMockSessionRepository(ctrl)

		jwt := []byte("secret")

		svc := NewAuthService(manager, session, jwt)

		ctx := context.Background()
		req := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwic2Vzc2lvbl9pZCI6MX0.3bKrUUKhFWJocpAKZsQGS5eMyucKfIWIRig29EM-5m8"
		userID := 1
		sessionID := 1
		publicKey := []byte("public key")
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		}, nil)

		session.EXPECT().Get(ctx, sessionID).Return(repository.CredentialInfo{
			ID:        1,
			UserID:    2,
			PublicKey: publicKey,
		}, nil)

		info, err := svc.GetSessionByCredentialToken(ctx, req)
		assert.ErrorIs(t, err, service.ErrSessionNotConnectedWithUser)
		assert.Empty(t, info)
	})
}
