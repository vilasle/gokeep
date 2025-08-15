package grpc

import (
	"context"
	"time"

	"github.com/vilasle/gokeep/internal/service/client"
)

// TODO implement it
type BankCardService struct {
	socket string
}

func NewBankCardService(socket string) *BankCardService {
	return &BankCardService{
		socket: socket,
	}
}

func (s *BankCardService) Save(ctx context.Context, req client.BankCardSaveRequest) client.BankCardSaveResponse {
	return client.BankCardSaveResponse{}
}

func (s *BankCardService) List(ctx context.Context) client.BankCardListResponse {
	return client.BankCardListResponse{
		Result: []client.BankCardView{
			{
				Number:  "1234567890123456",
				Expires: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func (s *BankCardService) Get(ctx context.Context, req client.BankCardGetRequest) client.BankCardGetResponse {
	return client.BankCardGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek"),
			Data: []byte("data"),
		},
	}
}

func (s *BankCardService) Delete(ctx context.Context, req client.BankCardDeleteRequest) client.BankCardDeleteResponse {
	return client.BankCardDeleteResponse{}
}
