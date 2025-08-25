package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type LoginPasswordService struct {
	client pb.PrivateDataServiceClient
}

func NewLoginPasswordDataService(socket *grpc.ClientConn) *LoginPasswordService {
	return &LoginPasswordService{
		client: pb.NewPrivateDataServiceClient(socket),
	}
}

func (s *LoginPasswordService) Save(ctx context.Context, req client.LoginPasswordSaveRequest) client.LoginPasswordSaveResponse {
	return client.LoginPasswordSaveResponse{
		Error: "",
		ID:    1,
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek1"),
			Data: []byte(`{"login": "test_account", "password": "test_password"}`),
		},
	}
}

func (s *LoginPasswordService) Get(ctx context.Context, req client.LoginPasswordGetRequest) client.LoginPasswordGetResponse {
	return client.LoginPasswordGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek"),
			Data: []byte("data"),
		},
		Error: "",
	}
}

func (s *LoginPasswordService) Delete(ctx context.Context, req client.LoginPasswordDeleteRequest) client.LoginPasswordDeleteResponse {
	return client.LoginPasswordDeleteResponse{
		Error: "",
	}
}
