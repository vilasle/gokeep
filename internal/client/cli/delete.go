package cli

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

//DeleteLoginPassword - delete login password from external storage and local storage
func (c *CommandLineClient) DeleteLoginPassword(ctx context.Context, id int) error {
	externalId := c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass})
	if externalId == 0 {
		return nil
	}

	if err := c.externalServices.credentials.Delete(ctx, svc.DeleteRequest{ID: externalId, JWT: string(c.credential)}); err != nil {
		return err
	}

	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeUsepass})
}

//DeleteLoginPassword - delete bank card from external storage and local storage
func (c *CommandLineClient) DeleteBankCard(ctx context.Context, id int) error {
	externalId := c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass})
	if externalId == 0 {
		return nil
	}

	if err := c.externalServices.bankCard.Delete(ctx, svc.DeleteRequest{ID: externalId, JWT: string(c.credential)}); err != nil {
		return err
	}

	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBankCard})
}

//DeleteLoginPassword - delete text data from external storage and local storage
func (c *CommandLineClient) DeleteTextData(ctx context.Context, id int) error {
	externalId := c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass})
	if externalId == 0 {
		return nil
	}

	if err := c.externalServices.text.Delete(ctx, svc.DeleteRequest{ID: externalId, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypePlainText})
}

//DeleteLoginPassword - delete binary data from external storage and local storage
func (c *CommandLineClient) DeleteBinaryData(ctx context.Context, id int) error {
	externalId := c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass})
	if externalId == 0 {
		return nil
	}

	if err := c.externalServices.binary.Delete(ctx, svc.DeleteRequest{ID: externalId, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBinaryData})
}
