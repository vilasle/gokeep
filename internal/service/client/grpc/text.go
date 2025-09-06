package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TextDataService is wrapper over grpc.PrivateDataServiceClient for work with text data
type TextDataService struct {
	client proto.PrivateDataServiceClient
}

// NewTextDataService returns new instance of TextDataService
func NewTextDataService(socket *grpc.ClientConn) *TextDataService {
	return &TextDataService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

// Save - prepare and send request for grpc server
func (s *TextDataService) Save(ctx context.Context,
	req client.TextDataSaveRequest) (client.SaveResponse, error) {

	dto := proto.SaveTextDataRequest{
		Id:         int64(req.ID),
		Name:       req.Name,
		Data:       req.Text,
		Credential: &proto.ConfirmAssess{Token: req.JWT},
		Metadata:   make([]*proto.Metadata, 0, len(req.Metadata)),
	}

	for _, v := range req.Metadata {
		dto.Metadata = append(dto.Metadata, &proto.Metadata{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return handleSaveResponse(s.client.SaveTextData(ctx, &dto))
}

// Get - send Get request to grpc server and cast response to expected view
func (s *TextDataService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

// Delete - send Delete request to grpc server
func (s *TextDataService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, model.TypePlainText, s.client, req)
}
