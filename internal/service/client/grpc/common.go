package grpc

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/service/client"
	"github.com/vilasle/gokeep/proto"
)

func handleSaveResponse(resp *proto.EncryptedDataResponse, err error) (client.SaveResponse, error) {
	if err != nil {
		return client.SaveResponse{}, err
	}

	if resp.Error != "" {
		return client.SaveResponse{}, errors.New(resp.Error)
	}

	return client.SaveResponse{
		ID: int(resp.Entity.Id),
		Data: client.EncryptedData{
			View: resp.Entity.Data.View,
			DEK:  resp.Entity.Data.Dek,
			Data: resp.Entity.Data.Data,
		},
	}, nil
}
func get(ctx context.Context, svc proto.PrivateDataServiceClient, req client.GetRequest) ([]client.EncryptedEntity, error) {
	resp, err := svc.Get(ctx, &proto.GetDataRequest{
		Id:         int64(req.ID),
		Credential: &proto.ConfirmAssess{Token: req.JWT},
	})
	return handleGetResponse(resp, err)
}

func handleGetResponse(resp *proto.GetDataResponse, err error) ([]client.EncryptedEntity, error) {
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	result := make([]client.EncryptedEntity, len(resp.Data))
	for i, data := range resp.Data {
		result[i] = client.EncryptedEntity{
			ID: int(data.Id),
			Data: client.EncryptedData{
				View: data.Data.View,
				DEK:  data.Data.Dek,
				Data: data.Data.Data,
			},
		}
	}
	return result, nil
}

func delete(ctx context.Context, svc proto.PrivateDataServiceClient, req client.DeleteRequest) error {
	resp, err := svc.Delete(ctx, &proto.DeleteDataRequest{
		Id:         int64(req.ID),
		Credential: &proto.ConfirmAssess{Token: req.JWT},
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
