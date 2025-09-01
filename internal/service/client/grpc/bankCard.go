package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// BankCardService is wrapper over grpc.PrivateDataServiceClient for work with bank card data
type BankCardService struct {
	client proto.PrivateDataServiceClient
}

// NewBankCardService returns new instance of BankCardService
func NewBankCardService(socket *grpc.ClientConn) *BankCardService {
	return &BankCardService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

// Save - prepare and send request for grpc server
func (s *BankCardService) Save(ctx context.Context, req client.BankCardSaveRequest) (client.SaveResponse, error) {
	dto := proto.SaveBankCardRequest{
		Id:         int64(req.ID),
		Number:     req.Number,
		Cvv:        int64(req.CVV),
		Expires:    req.Expires,
		Credential: &proto.ConfirmAssess{Token: req.JWT},
		Metadata:   make([]*proto.Metadata, 0, len(req.Metadata)),
	}

	for _, v := range req.Metadata {
		dto.Metadata = append(dto.Metadata, &proto.Metadata{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return handleSaveResponse(s.client.SaveBankCard(ctx, &dto))
}

// Get - send Get request to grpc server and cast response to expected view
func (s *BankCardService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

// Delete - send Delete request to grpc server
func (s *BankCardService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, model.TypeBankCard, s.client, req)
}
