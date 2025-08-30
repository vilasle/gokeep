package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type BinaryDataService struct {
	client proto.PrivateDataServiceClient
}

func NewBinaryDataService(socket *grpc.ClientConn) *BinaryDataService {
	return &BinaryDataService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

func (s *BinaryDataService) Save(ctx context.Context, req client.BinaryDataSaveRequest) (client.SaveResponse, error) {
	dto := &proto.SaveBinaryDataRequest{
		Name:       req.Name,
		Data:       req.Data,
		Credential: &proto.ConfirmAssess{Token: req.JWT},
		Metadata:   make([]*proto.Metadata, 0, len(req.Metadata)),
	}

	for _, v := range req.Metadata {
		dto.Metadata = append(dto.Metadata, &proto.Metadata{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return handleSaveResponse(s.client.SaveBinaryData(ctx, dto))
}

func (s *BinaryDataService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

func (s *BinaryDataService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, s.client, req)
}
