package grpc

import (
	"context"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	pb "github.com/vilasle/gokeep/proto"
	"google.golang.org/grpc"
)

// TODO implement it
type BankCardService struct {
	client pb.PrivateDataServiceClient
}

func NewBankCardService(socket *grpc.ClientConn) *BankCardService {
	return &BankCardService{
		client: pb.NewPrivateDataServiceClient(socket),
	}
}

func (s *BankCardService) Save(ctx context.Context, req client.BankCardSaveRequest) (client.SaveResponse, error) {
	dto := proto.SaveBankCardRequest{
		Number:     req.Number,
		Cvv:        int64(req.CVV),
		Expires:    req.Expires,
		Credential: &pb.ConfirmAssess{Token: req.JWT},
	}

	resp, err := s.client.SaveBankCard(ctx, &dto)
	return handleSaveResponse(resp, err)
}

func (s *BankCardService) Get(ctx context.Context, req client.GetRequest) ([]client.EncryptedData, error) {
	return get(ctx, s.client, req)
}

func (s *BankCardService) Delete(ctx context.Context, req client.DeleteRequest) error {
	return delete(ctx, s.client, req)
}
