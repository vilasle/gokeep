package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

func (c *Client) DeleteLoginPassword(ctx context.Context, id int) error {
	r, err := c.localStorage.Get(ctx, repository.GetRequest{
		ID:   id,
		Type: model.TypeUsepass,
	})
	if err != nil {
		return err
	}

	if len(r) == 0 {
		return fmt.Errorf("no data found")
	}

	errs := make([]error, 0)
	for _, r := range r {
		if r.ExternalID == 0 {
			continue
		}

		err := c.externalServices.credentials.Delete(ctx, svc.DeleteRequest{
			ID:  r.ExternalID,
			JWT: string(c.credential),
		})
		if err != nil {
			errs = append(errs, err)
			continue
		}

		errs = append(errs, c.localStorage.Delete(ctx, repository.DeleteRequest{
			ID:   r.ID,
			Type: model.TypeUsepass,
		}))
	}

	return errors.Join(errs...)
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
