package grpc

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
	pb "github.com/vilasle/gokeep/proto"
)

func handleSaveResponse(resp *pb.EncryptedDataResponse, err error) (client.SaveResponse, error) {
	if err != nil {
		return client.SaveResponse{}, err
	}

	if resp.Error != "" {
		return client.SaveResponse{}, errors.New(resp.Error)
	}

	return client.SaveResponse{
		ID: int(resp.Id),
		Data: client.EncryptedData{
			View: resp.Data.View,
			DEK:  resp.Data.Dek,
			Data: resp.Data.Data,
		},
	}, nil
}
func get(ctx context.Context, svc proto.PrivateDataServiceClient, req client.GetRequest) ([]client.EncryptedData, error) {
	resp, err := svc.Get(ctx, &proto.GetDataRequest{
		Id:         int64(req.ID),
		Credential: &pb.ConfirmAssess{Token: req.JWT},
	})
	return handleGetResponse(resp, err)
}

func handleGetResponse(resp *proto.GetDataResponse, err error) ([]client.EncryptedData, error) {
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	result := make([]client.EncryptedData, len(resp.Data))
	for i, data := range resp.Data {
		result[i] = client.EncryptedData{
			View: data.View,
			DEK:  data.Dek,
			Data: data.Data,
		}
	}
	return result, nil
}

func delete(ctx context.Context, svc proto.PrivateDataServiceClient, req client.DeleteRequest) error {
	resp, err := svc.Delete(ctx, &proto.DeleteDataRequest{
		Id:         int64(req.ID),
		Credential: &pb.ConfirmAssess{Token: req.JWT},
	})
	return handleDeleteResponse(resp, err)
}

func handleDeleteResponse(resp *proto.DeleteDataResponse, err error) error {
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}
