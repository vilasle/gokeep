package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/model"
)

func Test_AuthService_Register(t *testing.T) {
	behaviorFind := func(m *MockUserRepository, ctx context.Context, username string, err error) {
		m.EXPECT().Find(ctx, username).Return(nil, err).Times(1)
	}

	behaviorAdd := func(m *MockUserRepository, ctx context.Context, err error) {
		m.EXPECT().Add(ctx, gomock.Any()).Return(err).Times(1)
	}

	testCases := []struct {
		name          string
		username      string
		password      string
		key           []byte
		expectedToken string
		addErr        error
		findErr       error
		wantErr       bool
	}{
		{
			name:          "success adding new user",
			username:      "login",
			password:      "password",
			key:           []byte("secret-key"),
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			findErr:       model.ErrUserNotFound,
		},
		{
			name:          "error on adding new user",
			username:      "login",
			password:      "password",
			key:           []byte("secret-key"),
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			findErr:       model.ErrUserNotFound,
			addErr:        errors.New("error"),
			wantErr:       true,
		},
		{
			name:          "error on finding user by username",
			username:      "login",
			password:      "password",
			key:           []byte("secret-key"),
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			findErr:       errors.New("errors"),
			wantErr:       true,
		},
		{
			name:          "user already exists",
			username:      "login",
			password:      "password",
			key:           []byte("secret-key"),
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			wantErr:       true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockRepositoryCollector(ctrl)
			userRep := NewMockUserRepository(ctrl)
			rep.EXPECT().User().Return(userRep).Times(1)
			rep.EXPECT().Private().Return(nil).Times(3)
			rep.EXPECT().Encryption().Return(nil).Times(3)

			ctx := context.Background()
			behaviorFind(userRep, ctx, tt.username, tt.findErr)
			if tt.findErr == model.ErrUserNotFound {
				behaviorAdd(userRep, ctx, tt.addErr)
			}
			manager := model.NewModelManager(rep)

			authService := NewAuthService(manager, tt.key)

			token, err := authService.Register(ctx, tt.username, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedToken, token)
		})
	}
}

func Test_AuthService_Login(t *testing.T) {
	key := []byte("secret-key")

	type argsFind struct {
		username string
		password string
		err      error
	}

	behaviorFind := func(m *MockUserRepository, ctx context.Context, args argsFind, user *model.User) {
		m.EXPECT().Find(ctx, args.username).Return(user, args.err).Times(1)
	}
	testCases := []struct {
		name     string
		username string
		password string
		argsFind
		expectedToken string
		wantErr       bool
	}{
		{
			name:          "success login",
			username:      "login",
			password:      "password",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			argsFind: argsFind{
				username: "login",
				password: "password",
				err:      nil,
			},
		},
		{
			name:          "not found user",
			username:      "login",
			password:      "password",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			argsFind: argsFind{
				username: "login",
				password: "password",
				err:      model.ErrUserNotFound,
			},
			wantErr: true,
		},
		{
			name:          "wrong password",
			username:      "login",
			password:      "password",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.C8JtruNY6MFldmT-5cb17sogzE-xkXokRcQLnLE-9qM",
			argsFind: argsFind{
				username: "login",
				password: "another password",
				err:      nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockRepositoryCollector(ctrl)
			userRep := NewMockUserRepository(ctrl)
			rep.EXPECT().User().Return(userRep).Times(1)
			rep.EXPECT().Private().Return(nil).Times(3)
			rep.EXPECT().Encryption().Return(nil).Times(3)

			ctx := context.Background()

			manager := model.NewModelManager(rep)

			user := manager.Users.New(tt.argsFind.username, tt.argsFind.password)

			behaviorFind(userRep, ctx, tt.argsFind, user)

			authService := NewAuthService(manager, key)

			token, err := authService.Login(ctx, tt.username, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedToken, token)
		})
	}
}

func Test_AuthService_Valid(t *testing.T) {
	type argsGet struct {
		id  int64
		err error
	}
	behaviorGet := func(m *MockUserRepository, ctx context.Context, args argsGet, user *model.User) {
		m.EXPECT().Get(ctx, args.id).Return(user, args.err).Times(1)
	}

	testCases := []struct {
		name       string
		key        []byte
		token      string
		validToken bool
		argsGet
		wantErr bool
	}{
		{
			name: "token is valid",
			key:  []byte("secret-key"),
			//id = 1
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.Ye0qYHbElGyNpRhrKfpdDlcq3hT23hob2viGYG-fZZI",
			validToken: true,
			argsGet: argsGet{
				id:  1,
				err: nil,
			},
		},
		{
			name: "user not found",
			key:  []byte("secret-key"),
			//id = 1
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.Ye0qYHbElGyNpRhrKfpdDlcq3hT23hob2viGYG-fZZI",
			validToken: true,
			argsGet: argsGet{
				id:  1,
				err: model.ErrUserNotFound,
			},
			wantErr: true,
		},
		{
			name: "invalid token",
			key:  []byte("secret-key"),
			//id = 1
			token:      "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NTExMTI0NTksImV4cCI6MTc4MjY0ODQ1OSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIk5hbWUiOiJUZXN0In0.ri0VNHXivGbsIZbRBdk24duVExopIWR4zQ8QgG5yAKU",
			validToken: false,
			argsGet: argsGet{
				id:  1,
				err: nil,
			},
			wantErr: true,
		},
		{
			name: "valid token but invalid content",
			key:  []byte("secret-key"),
			//id = 1
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IlNvbWVUZXN0VG9rZW4ifQ.bOZSn0b2xlm-947d1DMrapn1TcNlWmS4JvvsXwTA3WY",
			validToken: false,
			argsGet: argsGet{
				id:  1,
				err: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rep := NewMockRepositoryCollector(ctrl)
			userRep := NewMockUserRepository(ctrl)
			rep.EXPECT().User().Return(userRep).Times(1)
			rep.EXPECT().Private().Return(nil).Times(3)
			rep.EXPECT().Encryption().Return(nil).Times(3)

			ctx := context.Background()

			manager := model.NewModelManager(rep)

			if tt.validToken {
				user := manager.Users.New("test", "test")
				behaviorGet(userRep, ctx, tt.argsGet, user)
			}

			authService := NewAuthService(manager, tt.key)

			err := authService.Valid(ctx, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

		})
	}
}
