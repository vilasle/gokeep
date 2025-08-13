package grpc

type GRPCAuthService struct {
	socket string
}

func NewGRPCAuthService(socket string) *GRPCAuthService {
	return &GRPCAuthService{
		socket: socket,
	}
}

func (s *GRPCAuthService) CreateAccount(accountName, password string) error {
	panic("not implemented")
}

func (s *GRPCAuthService) Login(accountName, password string) (cread string, err error) {
	panic("not implemented")
}
