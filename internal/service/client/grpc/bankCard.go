package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type BankCardService struct {
	client proto.PrivateDataServiceClient
}

func NewBankCardService(socket *grpc.ClientConn) *BankCardService {
	return &BankCardService{
		client: proto.NewPrivateDataServiceClient(socket),
	}
}

func (s *BankCardService) Save(ctx context.Context, req client.BankCardSaveRequest) (client.SaveResponse, error) {
	dto := proto.SaveBankCardRequest{
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

func (s *BankCardService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedEntity, error) {
	return get(ctx, s.client, req)
}

func (s *BankCardService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, s.client, req)
}
