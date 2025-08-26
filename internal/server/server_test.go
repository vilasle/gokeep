package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/huandu/go-assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/service"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

func TestCreateAccount(t *testing.T) {
	type authMockArgs struct {
		err   error
		input service.RegisterLoginUser
	}
	testCase := []struct {
		name  string
		input *pb.CreateAccountRequest
		ctx   context.Context
		authMockArgs
		expectedErr string
	}{
		{
			name: "success creating account",
			ctx:  context.Background(),
			input: &pb.CreateAccountRequest{
				Login:    "test_account",
				Password: "test_password",
			},
			authMockArgs: authMockArgs{
				err: nil,
				input: service.RegisterLoginUser{
					Username: "test_account",
					Password: "test_password",
				},
			},
			expectedErr: "",
		},
		{
			name: "creating account failed",
			ctx:  context.Background(),
			input: &pb.CreateAccountRequest{
				Login:    "test_account",
				Password: "test_password",
			},
			authMockArgs: authMockArgs{
				err: errors.New("user already exists"),
				input: service.RegisterLoginUser{
					Username: "test_account",
					Password: "test_password",
				},
			},
			expectedErr: "user already exists",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)

			authMock.EXPECT().
				Register(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.err)

			s := &Server{auth: authMock}

			resp, err := s.CreateAccount(tt.ctx, tt.input)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestLogin(t *testing.T) {
	type authMockArgs struct {
		err    error
		input  service.RegisterLoginUser
		output string
	}

	login := "test_account"
	password := "test_password"
	publicKey := []byte("test_public_key")
	token := "test_token"

	testCase := []struct {
		name   string
		input  *pb.LoginRequest
		output *pb.LoginResponse
		ctx    context.Context
		authMockArgs
		expectedErr   string
		expectedToken string
	}{
		{
			name: "success creating account",
			ctx:  context.Background(),
			input: &pb.LoginRequest{
				Login:     login,
				Password:  password,
				PublicKey: publicKey,
			},
			output: &pb.LoginResponse{
				Token: token,
				Error: "",
			},
			authMockArgs: authMockArgs{
				err: nil,
				input: service.RegisterLoginUser{
					Username:  login,
					Password:  password,
					PublicKey: publicKey,
				},
				output: token,
			},
			expectedToken: token,
			expectedErr:   "",
		},
		{
			name: "creating account failed",
			ctx:  context.Background(),
			input: &pb.LoginRequest{
				Login:     login,
				Password:  password,
				PublicKey: publicKey,
			},
			output: &pb.LoginResponse{
				Token: "",
				Error: "user not found",
			},
			authMockArgs: authMockArgs{
				err: errors.New("user not found"),
				input: service.RegisterLoginUser{
					Username:  login,
					Password:  password,
					PublicKey: publicKey,
				},
				output: "",
			},
			expectedToken: "",
			expectedErr:   "user not found",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)

			authMock.EXPECT().
				Login(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			s := &Server{auth: authMock}

			resp, err := s.Login(tt.ctx, tt.input)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedErr, resp.Error)
			assert.Equal(t, tt.expectedToken, resp.Token)
		})
	}
}

func TestSaveLoginPassword(t *testing.T) {
	type authMockArgs struct {
		input  string
		output service.SessionInfo
		err    error
	}

	type credMockArgs struct {
		input  service.AddLoginPassword
		output service.PrivateDataResponse
		err    error
	}

	login := "test_account"
	password := "test_password"
	token := "test_token"
	userID := 1
	sessionID := 1
	entityID := 1

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	testCase := []struct {
		name   string
		input  *pb.SaveLoginPasswordRequest
		output *pb.EncryptedDataResponse
		ctx    context.Context
		authMockArgs
		*credMockArgs
	}{
		{
			name: "success saving login password",
			ctx:  context.Background(),
			input: &pb.SaveLoginPasswordRequest{
				Login:    login,
				Password: password,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id: int64(entityID),
					Data: &pb.EncryptedData{
						Data: []byte("test_data"),
						View: login,
						Dek:  []byte("test_dek"),
					},
				},
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			credMockArgs: &credMockArgs{
				input: service.AddLoginPassword{
					UserID:   userID,
					Username: login,
					Password: password,
				},
				output: service.PrivateDataResponse{
					ID:   entityID,
					View: login,
					Data: "test_data",
					Key:  "test_dek",
				},
			},
		},
		{
			name: "wrong public key",
			ctx:  context.Background(),
			input: &pb.SaveLoginPasswordRequest{
				Login:    login,
				Password: password,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "public key is not valid",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: []byte("invalid_public_key"),
				},
				err: nil,
			},
			credMockArgs: nil,
		},
		{
			name: "auth service failed",
			ctx:  context.Background(),
			input: &pb.SaveLoginPasswordRequest{
				Login:    login,
				Password: password,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "auth request failed",
			},
			authMockArgs: authMockArgs{
				input:  token,
				output: service.SessionInfo{},
				err:    errors.New("auth request failed"),
			},
			credMockArgs: nil,
		},
		{
			name: "credential service failed",
			ctx:  context.Background(),
			input: &pb.SaveLoginPasswordRequest{
				Login:    login,
				Password: password,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Error: "credential request failed",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			credMockArgs: &credMockArgs{
				input: service.AddLoginPassword{
					UserID:   userID,
					Username: login,
					Password: password,
				},
				output: service.PrivateDataResponse{},
				err:    errors.New("credential request failed"),
			},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)
			credMock := NewMockLoginPasswordService(ctrl)

			authMock.EXPECT().
				GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			if tt.credMockArgs != nil {
				credMock.EXPECT().
					Add(tt.ctx, tt.credMockArgs.input, gomock.Any()).
					Return(tt.credMockArgs.output, tt.credMockArgs.err)
			}

			s := &Server{auth: authMock, cread: credMock}
			resp, err := s.SaveLoginPassword(tt.ctx, tt.input)
			require.NoError(t, err)

			if resp.Error != "" {
				assert.Equal(t, tt.output.Error, resp.Error)
				return
			}

			assert.Equal(t, tt.output.Entity.Id, resp.Entity.Id)
			assert.Equal(t, tt.output.Entity.Data.Data, resp.Entity.Data.Data)
			assert.Equal(t, tt.output.Entity.Data.Dek, resp.Entity.Data.Dek)
			assert.Equal(t, tt.output.Entity.Data.View, resp.Entity.Data.View)
		})
	}
}

func TestBankCard(t *testing.T) {
	type authMockArgs struct {
		input  string
		output service.SessionInfo
		err    error
	}

	type bankMockArgs struct {
		input  service.AddBankCard
		output service.PrivateDataResponse
		err    error
	}

	number := "1234567890"
	cvv := int64(123)
	expirationRaw := "12/31"
	expiration, err := time.Parse("01/06", expirationRaw)
	require.NoError(t, err)

	token := "test_token"
	userID := 1
	sessionID := 1
	entityID := 1

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	testCase := []struct {
		name   string
		input  *pb.SaveBankCardRequest
		output *pb.EncryptedDataResponse
		ctx    context.Context
		authMockArgs
		*bankMockArgs
	}{
		{
			name: "success saving bank card",
			ctx:  context.Background(),
			input: &pb.SaveBankCardRequest{
				Number:  number,
				Cvv:     cvv,
				Expires: expirationRaw,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id: int64(entityID),
					Data: &pb.EncryptedData{
						Data: []byte("test_data"),
						View: number,
						Dek:  []byte("test_dek"),
					},
				},
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			bankMockArgs: &bankMockArgs{
				input: service.AddBankCard{
					UserID:     userID,
					Number:     number,
					CVV:        int(cvv),
					Expiration: expiration,
				},
				output: service.PrivateDataResponse{
					ID:   entityID,
					View: number,
					Data: "test_data",
					Key:  "test_dek",
				},
			},
		},
		{
			name: "wrong public key",
			ctx:  context.Background(),
			input: &pb.SaveBankCardRequest{
				Number:  number,
				Cvv:     cvv,
				Expires: expirationRaw,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "public key is not valid",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: []byte("invalid_public_key"),
				},
				err: nil,
			},
			bankMockArgs: nil,
		},
		{
			name: "auth service failed",
			ctx:  context.Background(),
			input: &pb.SaveBankCardRequest{
				Number:  number,
				Cvv:     cvv,
				Expires: expirationRaw,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "auth request failed",
			},
			authMockArgs: authMockArgs{
				input:  token,
				output: service.SessionInfo{},
				err:    errors.New("auth request failed"),
			},
			bankMockArgs: nil,
		},
		{
			name: "bank service failed",
			ctx:  context.Background(),
			input: &pb.SaveBankCardRequest{
				Number:  number,
				Cvv:     cvv,
				Expires: expirationRaw,
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Error: "bank request failed",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			bankMockArgs: &bankMockArgs{
				input: service.AddBankCard{
					UserID:     userID,
					Number:     number,
					CVV:        int(cvv),
					Expiration: expiration,
				},
				output: service.PrivateDataResponse{},
				err:    errors.New("bank request failed"),
			},
		},
		{
			name: "wrong date",
			ctx:  context.Background(),
			input: &pb.SaveBankCardRequest{
				Number:  number,
				Cvv:     cvv,
				Expires: "wrong_date",
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil},
				Error: "parsing time \"wrong_date\" as \"01/06\": cannot parse \"wrong_date\" as \"01\"",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			bankMockArgs: nil,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)
			bankMock := NewMockBankCardService(ctrl)

			authMock.EXPECT().
				GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			if tt.bankMockArgs != nil {
				bankMock.EXPECT().
					Add(tt.ctx, tt.bankMockArgs.input, gomock.Any()).
					Return(tt.bankMockArgs.output, tt.bankMockArgs.err)
			}

			s := &Server{auth: authMock, bank: bankMock}
			resp, err := s.SaveBankCard(tt.ctx, tt.input)
			require.NoError(t, err)

			if resp.Error != "" {
				assert.Equal(t, tt.output.Error, resp.Error)
				return
			}

			assert.Equal(t, tt.output.Entity.Id, resp.Entity.Id)
			assert.Equal(t, tt.output.Entity.Data.Data, resp.Entity.Data.Data)
			assert.Equal(t, tt.output.Entity.Data.Dek, resp.Entity.Data.Dek)
			assert.Equal(t, tt.output.Entity.Data.View, resp.Entity.Data.View)
		})
	}
}

func TestTextData(t *testing.T) {
	type authMockArgs struct {
		input  string
		output service.SessionInfo
		err    error
	}

	type textMockArgs struct {
		input  service.AddTextData
		output service.PrivateDataResponse
		err    error
	}
	name := "test"
	text := "some test text"
	token := "test_token"
	userID := 1
	sessionID := 1
	entityID := 1

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	testCase := []struct {
		name   string
		input  *pb.SaveTextDataRequest
		output *pb.EncryptedDataResponse
		ctx    context.Context
		authMockArgs
		*textMockArgs
	}{
		{
			name: "success saving text data",
			ctx:  context.Background(),
			input: &pb.SaveTextDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id: int64(entityID),
					Data: &pb.EncryptedData{
						Data: []byte("test_data"),
						View: name,
						Dek:  []byte("test_dek"),
					},
				},
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			textMockArgs: &textMockArgs{
				input: service.AddTextData{
					UserID: userID,
					Name:   name,
					Text:   []byte(text),
				},
				output: service.PrivateDataResponse{
					ID:   entityID,
					View: name,
					Data: "test_data",
					Key:  "test_dek",
				},
			},
		},
		{
			name: "wrong public key",
			ctx:  context.Background(),
			input: &pb.SaveTextDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "public key is not valid",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: []byte("invalid_public_key"),
				},
				err: nil,
			},
			textMockArgs: nil,
		},
		{
			name: "auth service failed",
			ctx:  context.Background(),
			input: &pb.SaveTextDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},

				Error: "auth request failed",
			},
			authMockArgs: authMockArgs{
				input:  token,
				output: service.SessionInfo{},
				err:    errors.New("auth request failed"),
			},
			textMockArgs: nil,
		},
		{
			name: "credential service failed",
			ctx:  context.Background(),
			input: &pb.SaveTextDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Error: "credential request failed",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			textMockArgs: &textMockArgs{
				input: service.AddTextData{
					UserID: userID,
					Name:   name,
					Text:   []byte(text),
				},
				output: service.PrivateDataResponse{},
				err:    errors.New("credential request failed"),
			},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)
			textMock := NewMockTextDataService(ctrl)

			authMock.EXPECT().
				GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			if tt.textMockArgs != nil {
				textMock.EXPECT().
					Add(tt.ctx, tt.textMockArgs.input, gomock.Any()).
					Return(tt.textMockArgs.output, tt.textMockArgs.err)
			}

			s := &Server{auth: authMock, text: textMock}
			resp, err := s.SaveTextData(tt.ctx, tt.input)
			require.NoError(t, err)

			if resp.Error != "" {
				assert.Equal(t, tt.output.Error, resp.Error)
				return
			}

			assert.Equal(t, tt.output.Entity.Id, resp.Entity.Id)
			assert.Equal(t, tt.output.Entity.Data.Data, resp.Entity.Data.Data)
			assert.Equal(t, tt.output.Entity.Data.Dek, resp.Entity.Data.Dek)
			assert.Equal(t, tt.output.Entity.Data.View, resp.Entity.Data.View)
		})
	}
}

func TestBinaryData(t *testing.T) {
	type authMockArgs struct {
		input  string
		output service.SessionInfo
		err    error
	}

	type binaryMockArgs struct {
		input  service.AddBinaryData
		output service.PrivateDataResponse
		err    error
	}
	name := "test"
	text := "some test text"
	token := "test_token"
	userID := 1
	sessionID := 1
	entityID := 1

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	testCase := []struct {
		name   string
		input  *pb.SaveBinaryDataRequest
		output *pb.EncryptedDataResponse
		ctx    context.Context
		authMockArgs
		*binaryMockArgs
	}{
		{
			name: "success saving text data",
			ctx:  context.Background(),
			input: &pb.SaveBinaryDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id: int64(entityID),
					Data: &pb.EncryptedData{
						Data: []byte("test_data"),
						View: name,
						Dek:  []byte("test_dek"),
					},
				},
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			binaryMockArgs: &binaryMockArgs{
				input: service.AddBinaryData{
					UserID: userID,
					Name:   name,
					Data:   []byte(text),
				},
				output: service.PrivateDataResponse{
					ID:   entityID,
					View: name,
					Data: "test_data",
					Key:  "test_dek",
				},
			},
		},
		{
			name: "wrong public key",
			ctx:  context.Background(),
			input: &pb.SaveBinaryDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "public key is not valid",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: []byte("invalid_public_key"),
				},
				err: nil,
			},
			binaryMockArgs: nil,
		},
		{
			name: "auth service failed",
			ctx:  context.Background(),
			input: &pb.SaveBinaryDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Entity: &pb.EncryptedEntity{
					Id:   0,
					Data: nil,
				},
				Error: "auth request failed",
			},
			authMockArgs: authMockArgs{
				input:  token,
				output: service.SessionInfo{},
				err:    errors.New("auth request failed"),
			},
			binaryMockArgs: nil,
		},
		{
			name: "credential service failed",
			ctx:  context.Background(),
			input: &pb.SaveBinaryDataRequest{
				Name: name,
				Data: []byte(text),
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
			},
			output: &pb.EncryptedDataResponse{
				Error: "credential request failed",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    userID,
					SessionID: sessionID,
					PublicKey: publicKey,
				},
				err: nil,
			},
			binaryMockArgs: &binaryMockArgs{
				input: service.AddBinaryData{
					UserID: userID,
					Name:   name,
					Data:   []byte(text),
				},
				output: service.PrivateDataResponse{},
				err:    errors.New("credential request failed"),
			},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)
			binaryMock := NewMockBinaryDataService(ctrl)

			authMock.EXPECT().
				GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			if tt.binaryMockArgs != nil {
				binaryMock.EXPECT().
					Add(tt.ctx, tt.binaryMockArgs.input, gomock.Any()).
					Return(tt.binaryMockArgs.output, tt.binaryMockArgs.err)
			}

			s := &Server{auth: authMock, binary: binaryMock}
			resp, err := s.SaveBinaryData(tt.ctx, tt.input)
			require.NoError(t, err)

			if resp.Error != "" {
				assert.Equal(t, tt.output.Error, resp.Error)
				return
			}

			assert.Equal(t, tt.output.Entity.Id, resp.Entity.Id)
			assert.Equal(t, tt.output.Entity.Data.Data, resp.Entity.Data.Data)
			assert.Equal(t, tt.output.Entity.Data.Dek, resp.Entity.Data.Dek)
			assert.Equal(t, tt.output.Entity.Data.View, resp.Entity.Data.View)
		})
	}
}

func TestDelete(t *testing.T) {
	type authMockArgs struct {
		input  string
		output service.SessionInfo
		err    error
	}

	type deleteMockArgs struct {
		input service.DeletePrivateData
		err   error
	}
	token := "test_token"
	entityID := 1

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKey := pem.EncodeToMemory(publicKeyBlock)

	testCase := []struct {
		name                        string
		ctx                         context.Context
		input                       *pb.DeleteDataRequest
		output                      *pb.DeleteDataResponse
		authMockArgs                authMockArgs
		deleteLoginPasswordMockArgs *deleteMockArgs
		deleteBankCardMockArgs      *deleteMockArgs
		deleteTextDataMockArgs      *deleteMockArgs
		deleteBinaryMockArgs        *deleteMockArgs
	}{
		{
			name: "delete login password",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 2,
			},
			output: &pb.DeleteDataResponse{
				Error: "",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: &deleteMockArgs{
				input: service.DeletePrivateData{
					UserID: 1,
					ID:     entityID,
				},
				err: nil,
			},
			deleteBankCardMockArgs: nil,
			deleteTextDataMockArgs: nil,
			deleteBinaryMockArgs:   nil,
		},
		{
			name: "delete bank card",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 3,
			},
			output: &pb.DeleteDataResponse{
				Error: "",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs: &deleteMockArgs{
				input: service.DeletePrivateData{
					UserID: 1,
					ID:     entityID,
				},
				err: nil,
			},
			deleteTextDataMockArgs: nil,
			deleteBinaryMockArgs:   nil,
		},
		{
			name: "delete text data",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 4,
			},
			output: &pb.DeleteDataResponse{
				Error: "",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs:      nil,
			deleteTextDataMockArgs: &deleteMockArgs{
				input: service.DeletePrivateData{
					UserID: 1,
					ID:     entityID,
				},
				err: nil,
			},
			deleteBinaryMockArgs: nil,
		},
		{
			name: "delete binary data",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 5,
			},
			output: &pb.DeleteDataResponse{
				Error: "",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs:      nil,
			deleteTextDataMockArgs:      nil,
			deleteBinaryMockArgs: &deleteMockArgs{
				input: service.DeletePrivateData{
					UserID: 1,
					ID:     entityID,
				},
				err: nil,
			},
		},
		{
			name: "auth service failed",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 5,
			},
			output: &pb.DeleteDataResponse{
				Error: "auth failed",
			},
			authMockArgs: authMockArgs{
				input:  token,
				output: service.SessionInfo{},
				err:    errors.New("auth failed"),
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs:      nil,
			deleteTextDataMockArgs:      nil,
			deleteBinaryMockArgs:        nil,
		},
		{
			name: "delete failed",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 5,
			},
			output: &pb.DeleteDataResponse{
				Error: "delete failed",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs:      nil,
			deleteTextDataMockArgs:      nil,
			deleteBinaryMockArgs: &deleteMockArgs{
				input: service.DeletePrivateData{
					UserID: 1,
					ID:     entityID,
				},
				err: errors.New("delete failed"),
			},
		},
		{
			name: "unknown type of data",
			ctx:  context.Background(),
			input: &pb.DeleteDataRequest{
				Credential: &pb.ConfirmAssess{
					Token: token,
				},
				Id:   int64(entityID),
				Type: 6,
			},
			output: &pb.DeleteDataResponse{
				Error: "unknown type",
			},
			authMockArgs: authMockArgs{
				input: token,
				output: service.SessionInfo{
					UserID:    1,
					PublicKey: publicKey,
				},
				err: nil,
			},
			deleteLoginPasswordMockArgs: nil,
			deleteBankCardMockArgs:      nil,
			deleteTextDataMockArgs:      nil,
			deleteBinaryMockArgs:        nil,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authMock := NewMockAuthService(ctrl)
			credMock := NewMockLoginPasswordService(ctrl)
			bankMock := NewMockBankCardService(ctrl)
			textMock := NewMockTextDataService(ctrl)
			binaryMock := NewMockBinaryDataService(ctrl)

			authMock.EXPECT().
				GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
				Return(tt.authMockArgs.output, tt.authMockArgs.err)

			if tt.deleteLoginPasswordMockArgs != nil {
				credMock.EXPECT().
					Delete(tt.ctx, tt.deleteLoginPasswordMockArgs.input).
					Return(tt.deleteLoginPasswordMockArgs.err)
			}

			if tt.deleteBankCardMockArgs != nil {
				bankMock.EXPECT().
					Delete(tt.ctx, tt.deleteBankCardMockArgs.input).
					Return(tt.deleteBankCardMockArgs.err)
			}

			if tt.deleteTextDataMockArgs != nil {
				textMock.EXPECT().
					Delete(tt.ctx, tt.deleteTextDataMockArgs.input).
					Return(tt.deleteTextDataMockArgs.err)
			}

			if tt.deleteBinaryMockArgs != nil {
				binaryMock.EXPECT().
					Delete(tt.ctx, tt.deleteBinaryMockArgs.input).
					Return(tt.deleteBinaryMockArgs.err)
			}

			s := &Server{auth: authMock, cread: credMock, bank: bankMock, text: textMock, binary: binaryMock}
			resp, err := s.Delete(tt.ctx, tt.input)
			require.NoError(t, err)

			if resp.Error != "" {
				assert.Equal(t, tt.output.Error, resp.Error)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	// type authMockArgs struct {
	// 	input  string
	// 	output service.SessionInfo
	// 	err    error
	// }

	// type getMockArgs struct {
	// 	input service.GetPrivateData
	// 	err   error
	// }
	// token := "test_token"
	// entityID := 1

	// pk, err := rsa.GenerateKey(rand.Reader, 2048)
	// require.NoError(t, err)

	// publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	// require.NoError(t, err)

	// publicKeyBlock := &pem.Block{
	// 	Type:  "RSA PUBLIC KEY",
	// 	Bytes: publicKeyBytes,
	// }
	// publicKey := pem.EncodeToMemory(publicKeyBlock)

	// testCase := []struct {
	// 	name         string
	// 	ctx          context.Context
	// 	input        *pb.GetDataRequest
	// 	output       *pb.GetDataResponse
	// 	authMockArgs authMockArgs
	// 	creadMockArgs  *getMockArgs
	// 	bankMockArgs  *getMockArgs
	// 	textMockArgs  *getMockArgs
	// 	binaryMockArgs  *getMockArgs
	// }{
	// 	{
	// 		name: "get specific login password",
	// 		ctx:  context.Background(),
	// 		input: &pb.GetDataRequest{
	// 			Credential: &pb.ConfirmAssess{
	// 				Token: token,
	// 			},
	// 			Id:   int64(entityID),
	// 			Type: 2,
	// 		},
	// 		output: &pb.GetDataResponse{},
	// 		authMockArgs: authMockArgs{
	// 			input: token,
	// 			output: service.SessionInfo{
	// 				UserID:    1,
	// 				PublicKey: publicKey,
	// 			},
	// 			err: nil,
	// 		},
	// 		creadMockArgs: &getMockArgs{
	// 			input: service.DeletePrivateData{
	// 				UserID: 1,
	// 				ID:     entityID,
	// 			},
	// 			err: nil,
	// 		},
	// 		bankMockArgs: nil,
	// 		textMockArgs: nil,
	// 		binaryMockArgs:   nil,
	// 	},
	// }

	// for _, tt := range testCase {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		ctrl := gomock.NewController(t)
	// 		defer ctrl.Finish()

	// 		authMock := NewMockAuthService(ctrl)
	// 		credMock := NewMockLoginPasswordService(ctrl)
	// 		bankMock := NewMockBankCardService(ctrl)
	// 		textMock := NewMockTextDataService(ctrl)
	// 		binaryMock := NewMockBinaryDataService(ctrl)

	// 		authMock.EXPECT().
	// 			GetSessionByCredentialToken(tt.ctx, tt.authMockArgs.input).
	// 			Return(tt.authMockArgs.output, tt.authMockArgs.err)

	// 		if tt.deleteLoginPasswordMockArgs != nil {
	// 			credMock.EXPECT().
	// 				Delete(tt.ctx, tt.deleteLoginPasswordMockArgs.input).
	// 				Return(tt.deleteLoginPasswordMockArgs.err)
	// 		}

	// 		if tt.deleteBankCardMockArgs != nil {
	// 			bankMock.EXPECT().
	// 				Delete(tt.ctx, tt.deleteBankCardMockArgs.input).
	// 				Return(tt.deleteBankCardMockArgs.err)
	// 		}

	// 		if tt.deleteTextDataMockArgs != nil {
	// 			textMock.EXPECT().
	// 				Delete(tt.ctx, tt.deleteTextDataMockArgs.input).
	// 				Return(tt.deleteTextDataMockArgs.err)
	// 		}

	// 		if tt.deleteBinaryMockArgs != nil {
	// 			binaryMock.EXPECT().
	// 				Delete(tt.ctx, tt.deleteBinaryMockArgs.input).
	// 				Return(tt.deleteBinaryMockArgs.err)
	// 		}

	// 		s := &Server{auth: authMock, cread: credMock, bank: bankMock, text: textMock, binary: binaryMock}
	// 		resp, err := s.Delete(tt.ctx, tt.input)
	// 		require.NoError(t, err)

	// 		if resp.Error != "" {
	// 			assert.Equal(t, tt.output.Error, resp.Error)
	// 			return
	// 		}
	// 	})
	// }
}

func TestPrepareServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := NewMockAuthService(ctrl)
	credMock := NewMockLoginPasswordService(ctrl)
	bankMock := NewMockBankCardService(ctrl)
	textMock := NewMockTextDataService(ctrl)
	binaryMock := NewMockBinaryDataService(ctrl)

	server, err := NewServer(Config{
		AuthService:          authMock,
		LoginPasswordService: credMock,
		BankCardService:      bankMock,
		TextDataService:      textMock,
		BinaryDataService:    binaryMock,
	}, WithLogger)

	require.NoError(t, err)

	go func() {
		err := server.Listen()
		if err != nil {
			require.ErrorIs(t, err, grpc.ErrServerStopped)
		}
	}()

	time.Sleep(1 * time.Second)
	server.Stop()
	time.Sleep(2 * time.Second)

}
