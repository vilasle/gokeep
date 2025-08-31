package cli

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

func (c *Client) DeleteLoginPassword(ctx context.Context, id int) error {
	if err := c.externalServices.credentials.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}

	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeUsepass})
}

func (c *Client) DeleteBankCard(ctx context.Context, id int) error {
	if err := c.externalServices.bankCard.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}

	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBankCard})
}

func (c *Client) DeleteTextData(ctx context.Context, id int) error {
	if err := c.externalServices.text.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypePlainText})
}

func (c *Client) DeleteBinaryData(ctx context.Context, id int) error {
	if err := c.externalServices.binary.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBinaryData})
}
