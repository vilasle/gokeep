package cli

import (
	"context"
	"fmt"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
)

func (c *CommandLineClient) GetLoginPassword(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeUsepass, id)
}
func (c *CommandLineClient) GetBankCard(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeBankCard, id)
}

func (c *CommandLineClient) GetTextData(ctx context.Context, id int) error {
	return c.get(ctx, model.TypePlainText, id)
}

func (c *CommandLineClient) GetBinaryData(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeBinaryData, id)
}

func (c *CommandLineClient) get(ctx context.Context, t model.Type, id int) error {
	var (
		response []repository.GetResponse
		err      error
	)

	//show specific data and all entity in short format
	if id > 0 {
		if response, err = c.localStorage.Get(ctx, repository.GetRequest{ID: id, Type: t}); err != nil {
			return err
		}
		if len(response) == 0 {
			return fmt.Errorf("no data found")
		}
		return c.showFullEntities(t, response...)

	} else {
		if response, err = c.localStorage.All(ctx, t); err != nil {
			return err
		}
		return c.showListOfEntities(response...)
	}
}
