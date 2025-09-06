package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// BinaryDataService is wrapper over grpc.PrivateDataServiceClient for work with binary data
type BinaryDataService struct {
	client proto.PrivateDataServiceClient
}

// NewBinaryDataService returns new instance of BinaryDataService
func NewBinaryDataService(socket *grpc.ClientConn) *BinaryDataService {
	return &BinaryDataService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

// Save - prepare and send request for grpc server
func (s *BinaryDataService) Save(ctx context.Context, req client.BinaryDataSaveRequest) (client.SaveResponse, error) {
	dto := &proto.SaveBinaryDataRequest{
		Id:         int64(req.ID),
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


// Get - send Get request to grpc server and cast response to expected view
func (s *BinaryDataService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

// Delete - send Delete request to grpc server
func (s *BinaryDataService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, model.TypeBinaryData, s.client, req)
}
