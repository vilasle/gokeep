package grpc

import (
	"context"
	"fmt"

	"github.com/vilasle/gokeep/internal/service/client"
)

// TODO implement it
type LoginPasswordService struct {
	socket string
}

func NewLoginPasswordDataService(socket string) *LoginPasswordService {
	return &LoginPasswordService{
		socket: socket,
	}
}

func (s *LoginPasswordService) Save(ctx context.Context, req client.LoginPasswordSaveRequest) client.LoginPasswordSaveResponse {
	fmt.Printf("create login %s and password %s in server\n", req.Login, req.Password)
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

func (s *LoginPasswordService) List(ctx context.Context) client.LoginPasswordListResponse {
	return client.LoginPasswordListResponse{
		Result: []client.LoginPasswordView{
			{
				Login: "test",
			},
		},
		Error: "",
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
