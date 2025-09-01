package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// BankCardService is wrapper over grpc.PrivateDataServiceClient for work with login password data
type LoginPasswordService struct {
	client proto.PrivateDataServiceClient
}

// NewLoginPasswordDataService returns new instance of LoginPasswordService
func NewLoginPasswordDataService(socket *grpc.ClientConn) *LoginPasswordService {
	return &LoginPasswordService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

// Save - prepare and send request for grpc server
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

// Get - send Get request to grpc server and cast response to expected view
func (s *LoginPasswordService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

// Delete - send Delete request to grpc server
func (s *LoginPasswordService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, model.TypeUsepass, s.client, req)
}
