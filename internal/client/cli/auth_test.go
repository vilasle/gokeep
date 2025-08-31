package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthService(ctrl)

	account := "test"
	password := "test"
	ctx := context.Background()

	authMock.EXPECT().CreateAccount(ctx, account, password).Return(nil)

	client := Client{
		auth: authMock,
	}

	err := client.CreateAccount(ctx, account, password)
	assert.NoError(t, err)
}

func TestClient_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMock := NewMockAuthService(ctrl)
		localStorage := NewMockClientRepository(ctrl)

		pwd, err := os.Getwd()
		require.NoError(t, err)

		account := "test"
		password := "test"
		publicKey := []byte("public key")
		ctx := context.Background()

		creadPath := filepath.Join(pwd, (account + creadExt))

		authMock.EXPECT().Login(ctx, account, password, publicKey).Return([]byte("cread token"), nil)
		localStorage.EXPECT().CreateScheme(ctx).Return(nil)

		client := Client{
			auth:             authMock,
			localStorage:     localStorage,
			publicKeyContent: publicKey,
			workspace: WorkplaceConfig{
				Credentials: PathInfo{Path: pwd},
			},
		}

		err = client.Login(ctx, account, password)
		assert.NoError(t, err)

		_, err = os.Stat(creadPath)
		assert.NoError(t, err)

		os.Remove(creadPath)
	})

	t.Run("creating schema failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMock := NewMockAuthService(ctrl)
		localStorage := NewMockClientRepository(ctrl)

		pwd, err := os.Getwd()
		require.NoError(t, err)

		account := "test"
		password := "test"
		publicKey := []byte("public key")
		ctx := context.Background()

		localStorage.EXPECT().CreateScheme(ctx).Return(errors.New("error"))

		client := Client{
			auth:             authMock,
			localStorage:     localStorage,
			publicKeyContent: publicKey,
			workspace: WorkplaceConfig{
				Credentials: PathInfo{Path: pwd},
			},
		}

		err = client.Login(ctx, account, password)
		assert.Error(t, err)
	})

	t.Run("login is failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMock := NewMockAuthService(ctrl)
		localStorage := NewMockClientRepository(ctrl)

		pwd, err := os.Getwd()
		require.NoError(t, err)

		account := "test"
		password := "test"
		publicKey := []byte("public key")
		ctx := context.Background()

		authMock.EXPECT().Login(ctx, account, password, publicKey).Return(nil, errors.New("error"))
		localStorage.EXPECT().CreateScheme(ctx).Return(nil)

		client := Client{
			auth:             authMock,
			localStorage:     localStorage,
			publicKeyContent: publicKey,
			workspace: WorkplaceConfig{
				Credentials: PathInfo{Path: pwd},
			},
		}

		err = client.Login(ctx, account, password)
		assert.Error(t, err)
	})

	t.Run("wrong credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authMock := NewMockAuthService(ctrl)
		localStorage := NewMockClientRepository(ctrl)

		pwd, err := os.Getwd()
		require.NoError(t, err)

		account := "test"
		password := "test"
		publicKey := []byte("public key")
		ctx := context.Background()

		authMock.EXPECT().Login(ctx, account, password, publicKey).Return([]byte{}, nil)
		localStorage.EXPECT().CreateScheme(ctx).Return(nil)

		client := Client{
			auth:             authMock,
			localStorage:     localStorage,
			publicKeyContent: publicKey,
			workspace: WorkplaceConfig{
				Credentials: PathInfo{Path: pwd},
			},
		}

		err = client.Login(ctx, account, password)
		assert.Error(t, err)
	})

}
