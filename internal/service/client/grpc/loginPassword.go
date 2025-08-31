package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type LoginPasswordService struct {
	client proto.PrivateDataServiceClient
}

func NewLoginPasswordDataService(socket *grpc.ClientConn) *LoginPasswordService {
	return &LoginPasswordService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

func (s *LoginPasswordService) Save(ctx context.Context,
	req client.LoginPasswordSaveRequest) (client.SaveResponse, error) {

	dto := proto.SaveLoginPasswordRequest{
		Id:         int64(req.ID),
		Login:      req.Login,
		Password:   req.Password,
		Credential: &proto.ConfirmAssess{Token: req.JWT},
		Metadata:   make([]*proto.Metadata, 0, len(req.Metadata)),
	}

	for _, v := range req.Metadata {
		dto.Metadata = append(dto.Metadata, &proto.Metadata{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return handleSaveResponse(s.client.SaveLoginPassword(ctx, &dto))
}

func (s *LoginPasswordService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

func (s *LoginPasswordService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, s.client, req)
}
