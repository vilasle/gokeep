package grpc

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

type GRPCAuthService struct {
	client proto.AccountServiceClient
}

func NewGRPCAuthService(socket *grpc.ClientConn) *GRPCAuthService {
	return &GRPCAuthService{
		client: proto.NewAccountServiceClient(socket),
	}
}

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
