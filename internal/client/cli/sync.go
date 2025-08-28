package cli

import (
	"context"
	"fmt"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	"github.com/vilasle/gokeep/internal/service/client"
)

func (c *Client) Sync(ctx context.Context) error {
	//get all information from external storage
	all := make([]repository.SaveRequest, 0)

	fmt.Print("getting logins and passwords...")
	if result, err := c.externalServices.credentials.Get(ctx,
		client.GetRequest{JWT: string(c.credential), Type: model.TypeUsepass},
	); err == nil {
		fmt.Printf("got %d records\n", len(result))
		all = append(all, prepareCollectionForSyncStorages(result, model.TypeUsepass)...)
	} else {
		fmt.Print("failed\n")
		return err
	}

	fmt.Print("getting bank card...")
	if result, err := c.externalServices.bankCard.Get(ctx,
		client.GetRequest{JWT: string(c.credential), Type: model.TypeBankCard},
	); err == nil {
		fmt.Printf("got %d records\n", len(result))
		all = append(all, prepareCollectionForSyncStorages(result, model.TypeBankCard)...)
	} else {
		fmt.Print("failed\n")
		return err
	}

	fmt.Print("getting text data...")
	if result, err := c.externalServices.text.Get(ctx,
		client.GetRequest{JWT: string(c.credential), Type: model.TypePlainText},
	); err == nil {
		fmt.Printf("got %d records\n", len(result))
		all = append(all, prepareCollectionForSyncStorages(result, model.TypePlainText)...)
	} else {
		fmt.Print("failed\n")
		return err
	}

	fmt.Print("getting binary data...")
	if result, err := c.externalServices.binary.Get(ctx,
		client.GetRequest{JWT: string(c.credential), Type: model.TypeBinaryData},
	); err == nil {
		fmt.Printf("got %d records\n", len(result))
		all = append(all, prepareCollectionForSyncStorages(result, model.TypeBinaryData)...)
	} else {
		fmt.Print("failed\n")
		return err
	}

	err := c.localStorage.Rewrite(ctx, all)
	if err != nil {
		return err
	}

	return nil
}

func prepareCollectionForSyncStorages(entities []client.EncryptedEntity, tData int) []repository.SaveRequest {
	var result []repository.SaveRequest
	for _, entity := range entities {
		result = append(result, repository.SaveRequest{
			ID:       entity.ID,
			View:     entity.Data.View,
			DEK:      string(entity.Data.DEK),
			Data:     string(entity.Data.Data),
			Type:     tData,
			Metadata: prepareMetadataForLocalStorage(entity.Metadata),
		})
	}
	return result
}
