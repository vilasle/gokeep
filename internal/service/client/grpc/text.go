package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type TextDataService struct {
	client proto.PrivateDataServiceClient
}

func NewTextDataService(socket *grpc.ClientConn) *TextDataService {
	return &TextDataService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

func (s *TextDataService) Save(ctx context.Context,
	req client.TextDataSaveRequest) (client.SaveResponse, error) {

	dto := proto.SaveTextDataRequest{
		Name:       req.Name,
		Data:       req.Text,
		Credential: &proto.ConfirmAssess{Token: req.JWT},
	}

	resp, err := s.client.SaveTextData(ctx, &dto)
	return handleSaveResponse(resp, err)
}

func (s *TextDataService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

func (s *TextDataService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, s.client, req)
}
