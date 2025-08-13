package client

import (
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/internal/service/client/grpc"
)

func createAuthService(grpcSocket string) client.AuthService {
	return grpc.NewGRPCAuthService(grpcSocket)
}
