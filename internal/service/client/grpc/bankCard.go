package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type BankCardService struct {
	client pb.PrivateDataServiceClient
}

func NewBankCardService(socket *grpc.ClientConn) *BankCardService {
	return &BankCardService{
		client: pb.NewPrivateDataServiceClient(socket),
	}
}

func (s *BankCardService) Save(ctx context.Context, req client.BankCardSaveRequest) client.BankCardSaveResponse {
	return client.BankCardSaveResponse{
		ID: 1,
		Data: client.EncryptedData{
			View: "some card",
			DEK:  []byte("dek"),
			Data: []byte(`{"number": "1234567890123456","expires": "10/31","cvv":123}`),
		},
	}
}

func (s *BankCardService) Get(ctx context.Context, req client.BankCardGetRequest) client.BankCardGetResponse {
	return client.BankCardGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek"),
			Data: []byte(`"number": "1234567890123456","expires": "10/31","cvv":123`),
		},
	}
}

func (s *BankCardService) Delete(ctx context.Context, req client.BankCardDeleteRequest) client.BankCardDeleteResponse {
	return client.BankCardDeleteResponse{}
}
