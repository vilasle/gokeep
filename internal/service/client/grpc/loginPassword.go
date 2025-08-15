package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
)

//TODO implement it
type LoginPasswordService struct {
	socket string
}

func NewLoginPasswordDataService(socket string) *LoginPasswordService {
	return &LoginPasswordService{
		socket: socket,
	}
}

func (s *LoginPasswordService) Save(ctx context.Context, req client.LoginPasswordSaveRequest) client.LoginPasswordSaveResponse {
	return  client.LoginPasswordSaveResponse{
		Error: "",
	}
}

func (s *LoginPasswordService) List(ctx context.Context) client.LoginPasswordListResponse {
	return  client.LoginPasswordListResponse{
		Result: []client.LoginPasswordView{
			{
				Login:	"test",
			},
		},
		Error: "",
	}
}

func (s *LoginPasswordService) Get(ctx context.Context, req client.LoginPasswordGetRequest) client.LoginPasswordGetResponse {
	return client.LoginPasswordGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK: []byte("dek"),
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
