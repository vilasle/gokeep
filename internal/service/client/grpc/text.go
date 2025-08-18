package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
)

// TODO implement it
type TextDataService struct {
	socket string
}

func NewTextDataService(socket string) *TextDataService {
	return &TextDataService{
		socket: socket,
	}
}

func (s *TextDataService) Save(ctx context.Context, req client.TextDataSaveRequest) client.TextDataSaveResponse {
	return client.TextDataSaveResponse{
		ID: 1,
		Data: client.EncryptedData{
			View: req.Name,
			DEK:  []byte("dek"),
			Data: req.Text,
		},
	}
}

func (s *TextDataService) List(ctx context.Context) client.TextDataListResponse {
	return client.TextDataListResponse{
		Result: []client.TextDataView{
			{
				Name: "test",
			},
		},
	}
}

func (s *TextDataService) Get(ctx context.Context, req client.TextDataGetRequest) client.TextDataGetResponse {
	return client.TextDataGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek"),
			Data: []byte("data"),
		},
	}
}

func (s *TextDataService) Delete(ctx context.Context, req client.TextDataDeleteRequest) client.TextDataDeleteResponse {
	return client.TextDataDeleteResponse{}
}
