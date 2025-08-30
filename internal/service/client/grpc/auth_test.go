package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCAuthService_New(t *testing.T) {
	conn, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer conn.Close()

	srv := grpc.NewServer()
	defer srv.Stop()

	go srv.Serve(conn)

	clientConn, err := grpc.NewClient(conn.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()

	_ = NewGRPCAuthService(clientConn)
}
func TestGRPCAuthService_CreateAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"

		mock.EXPECT().CreateAccount(ctx, &proto.CreateAccountRequest{
			Login:    login,
			Password: password,
		}).Return(&proto.CreateAccountResponse{
			Error: "",
		}, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		err := svc.CreateAccount(ctx, login, password)
		assert.NoError(t, err)
	})

	t.Run("error client", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"

		mock.EXPECT().CreateAccount(ctx, &proto.CreateAccountRequest{
			Login:    login,
			Password: password,
		}).Return(&proto.CreateAccountResponse{}, errors.New("error"))

		svc := GRPCAuthService{
			client: mock,
		}

		err := svc.CreateAccount(ctx, login, password)
		assert.Error(t, err)
	})

	t.Run("error response", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"

		mock.EXPECT().CreateAccount(ctx, &proto.CreateAccountRequest{
			Login:    login,
			Password: password,
		}).Return(&proto.CreateAccountResponse{Error: "error"}, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		err := svc.CreateAccount(ctx, login, password)
		assert.Error(t, err)
	})
	t.Run("invalid response", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"

		mock.EXPECT().CreateAccount(ctx, &proto.CreateAccountRequest{
			Login:    login,
			Password: password,
		}).Return(nil, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		err := svc.CreateAccount(ctx, login, password)
		assert.Error(t, err)
	})
}

func TestGRPCAuthService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"
		publicKey := []byte("test_public_key")

		mock.EXPECT().Login(ctx, &proto.LoginRequest{
			Login:     login,
			Password:  password,
			PublicKey: publicKey,
		}).Return(&proto.LoginResponse{
			Token: "test_token",
			Error: "",
		}, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		token, err := svc.Login(ctx, login, password, publicKey)
		assert.NoError(t, err)
		assert.Equal(t, []byte("test_token"), token)
	})

	t.Run("client error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"
		publicKey := []byte("test_public_key")

		mock.EXPECT().Login(ctx, &proto.LoginRequest{
			Login:     login,
			Password:  password,
			PublicKey: publicKey,
		}).Return(nil, errors.New("error"))

		svc := GRPCAuthService{
			client: mock,
		}

		token, err := svc.Login(ctx, login, password, publicKey)
		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("response error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"
		publicKey := []byte("test_public_key")

		mock.EXPECT().Login(ctx, &proto.LoginRequest{
			Login:     login,
			Password:  password,
			PublicKey: publicKey,
		}).Return(&proto.LoginResponse{Error: "error"}, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		token, err := svc.Login(ctx, login, password, publicKey)
		assert.Error(t, err)
		assert.Nil(t, token)
	})
	t.Run("response is invalid", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockAccountServiceClient(ctrl)

		login, password := "test", "test_password"
		publicKey := []byte("test_public_key")

		mock.EXPECT().Login(ctx, &proto.LoginRequest{
			Login:     login,
			Password:  password,
			PublicKey: publicKey,
		}).Return(nil, nil)

		svc := GRPCAuthService{
			client: mock,
		}

		token, err := svc.Login(ctx, login, password, publicKey)
		assert.Error(t, err)
		assert.Nil(t, token)
	})
}
