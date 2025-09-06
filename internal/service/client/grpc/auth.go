package grpc

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

//GRPCAuthService - client service for registration and login on server
type GRPCAuthService struct {
	client proto.AccountServiceClient
}

//GRPCAuthService - new instance of GRPCAuthService
func NewGRPCAuthService(socket *grpc.ClientConn) *GRPCAuthService {
	return &GRPCAuthService{
		client: proto.NewAccountServiceClient(socket),
	}
}

//CreateAccount - send request to GRPC server for creating new account
func (s *GRPCAuthService) CreateAccount(ctx context.Context, accountName, password string) error {
	dto := &proto.CreateAccountRequest{
		Login:    accountName,
		Password: password,
	}

	resp, err := s.client.CreateAccount(ctx, dto)
	if err != nil {
		return err
	}

	if resp != nil && resp.Error != "" {
		return errors.New(resp.Error)
	}

	if resp == nil {
		return errors.New("response is invalid")
	}
	return nil
}

//Login - send request to GRPC server for login. Pass login, password and client public key
//server has to return JWT token which will be saved on file workspace/{account_name}.cred 
func (s *GRPCAuthService) Login(ctx context.Context, accountName, password string, publicKey []byte) ([]byte, error) {
	dto := proto.LoginRequest{
		Login:     accountName,
		Password:  password,
		PublicKey: publicKey,
	}

	resp, err := s.client.Login(ctx, &dto)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	if resp == nil {
		return nil, errors.New("response is invalid")
	}

	return []byte(resp.Token), nil
}
