package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type TextDataService struct {
	client pb.PrivateDataServiceClient
}

func NewTextDataService(socket *grpc.ClientConn) *TextDataService {
	return &TextDataService{
		client: pb.NewPrivateDataServiceClient(socket),
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
