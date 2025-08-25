package grpc

import (
	"fmt"

	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

type GRPCAuthService struct {
	socket pb.AccountServiceClient
}

// TODO implement me
func NewGRPCAuthService(socket *grpc.ClientConn) *GRPCAuthService {
	return &GRPCAuthService{
		socket: pb.NewAccountServiceClient(socket),
	}
}

func (s *GRPCAuthService) CreateAccount(accountName, password string, publicKey []byte) error {
	fmt.Println("created account")
	return nil
}

func (s *GRPCAuthService) Login(accountName, password string) ([]byte, error) {
	fmt.Printf("%s login\n", accountName)
	return []byte("creadData"), nil
}
