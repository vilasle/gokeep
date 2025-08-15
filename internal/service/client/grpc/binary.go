package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
)

// TODO implement it
type BinaryDataService struct {
	socket string
}

func NewBinaryDataService(socket string) *BinaryDataService {
	return &BinaryDataService{
		socket: socket,
	}
}

func (s *BinaryDataService) Save(ctx context.Context, req client.BinaryDataSaveRequest) client.BinaryDataSaveResponse {
	return client.BinaryDataSaveResponse{}
}

func (s *BinaryDataService) List(ctx context.Context) client.BinaryDataListResponse {
	return client.BinaryDataListResponse{
		Result: []client.BinaryDataView{
			{
				Name: "test",
			},
		},
	}
}

func (s *BinaryDataService) Get(ctx context.Context, req client.BinaryDataGetRequest) client.BinaryDataGetResponse {
	return client.BinaryDataGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK: []byte("dek"),
			Data: []byte("data"),
		},
	}
}

func (s *BinaryDataService) Delete(ctx context.Context, req client.BinaryDataDeleteRequest) client.BinaryDataDeleteResponse {
	return client.BinaryDataDeleteResponse{}
}
