package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type BinaryDataService struct {
	client pb.PrivateDataServiceClient
}

func NewBinaryDataService(socket *grpc.ClientConn) *BinaryDataService {
	return &BinaryDataService{
		client: pb.NewPrivateDataServiceClient(socket),
	}
}

func (s *BinaryDataService) Save(ctx context.Context, req client.BinaryDataSaveRequest) client.BinaryDataSaveResponse {
	return client.BinaryDataSaveResponse{}
}

func (s *BinaryDataService) Get(ctx context.Context, req client.BinaryDataGetRequest) client.BinaryDataGetResponse {
	return client.BinaryDataGetResponse{
		Data: client.EncryptedData{
			View: "test",
			DEK:  []byte("dek"),
			Data: []byte("data"),
		},
	}
}

func (s *BinaryDataService) Delete(ctx context.Context, req client.BinaryDataDeleteRequest) client.BinaryDataDeleteResponse {
	return client.BinaryDataDeleteResponse{}
}
