package grpc

import "fmt"

type GRPCAuthService struct {
	socket string
}
//TODO implement me
func NewGRPCAuthService(socket string) *GRPCAuthService {
	return &GRPCAuthService{
		socket: socket,
	}
}

func (s *GRPCAuthService) CreateAccount(accountName, password string, publicKey []byte) error {
	fmt.Println("created account")
	return  nil
}

func (s *GRPCAuthService) Login(accountName, password string) ([]byte, error) {
	fmt.Printf("%s login\n", accountName)
	return []byte("creadData"), nil
}
