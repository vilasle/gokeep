package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

// TODO implement work with local repository
func (c *Client) SaveLoginPassword(ctx context.Context, login, password string, id int) error {
	data := svc.LoginPasswordSaveRequest{
		ID:       id,
		Login:    login,
		Password: password,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeUsepass, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response := c.externalServices.credentials.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, model.TypeUsepass, repository.SaveRequest{
		ID:         id,
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
		View:       login,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}

	return nil

}

func (c *Client) SaveBankCard(ctx context.Context, number, expires string, cvv, id int) error {
	date, err := time.Parse("01/06", expires)
	if err != nil {
		return fmt.Errorf("invalid date format, expected $month/$year, e.g 01/2020")
	}

	data := svc.BankCardSaveRequest{
		Number:  number,
		Expires: date,
		CVV:     cvv,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeBankCard, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response := c.externalServices.bankCard.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, model.TypeBankCard, repository.SaveRequest{
		ID:         id,
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveTextDataAsIs(ctx context.Context, text string, name string, id int) error {
	data := svc.TextDataSaveRequest{
		Text: []byte(text),
		Name: name,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypePlainText, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response := c.externalServices.text.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, model.TypePlainText, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveTextDataFromFile(ctx context.Context, path string, name string, id int) error {
	if name == "" {
		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		name = stat.Name()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	data := svc.TextDataSaveRequest{
		Text: content,
		Name: name,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypePlainText, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response := c.externalServices.text.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, model.TypeUsepass, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
		View:       response.Data.View,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveBinaryData(ctx context.Context, path string, name string, id int) error {
	if name == "" {
		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		name = stat.Name()
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	data := svc.BinaryDataSaveRequest{
		Data: content,
		Name: name,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeUsepass, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response := c.externalServices.binary.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, model.TypeBinaryData, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
		View:       response.Data.View,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) GetLoginPassword(ctx context.Context, id int) error {
	t := model.TypeUsepass
	response, err := c.localStorage.Get(ctx, t, repository.GetRequest{
		ID: id,
	})

	if err != nil {
		return err
	}

	if len(response) == 0 {
		return fmt.Errorf("no data found")
	}

	return c.showFullEntity(response[0], t)
}

// TODO implement work with local repository
func (c *Client) GetBankCard(ctx context.Context, id int) error {
	t := model.TypeBankCard
	response, err := c.localStorage.Get(ctx, t, repository.GetRequest{
		ID: id,
	})

	if err != nil {
		return err
	}

	if len(response) == 0 {
		return fmt.Errorf("no data found")
	}

	return c.showFullEntity(response[0], t)
}

// TODO implement work with local repository
func (c *Client) GetTextData(ctx context.Context, id int) error {
	t := model.TypePlainText
	response, err := c.localStorage.Get(ctx, t, repository.GetRequest{
		ID: id,
	})

	if err != nil {
		return err
	}

	if len(response) == 0 {
		return fmt.Errorf("no data found")
	}

	return c.showFullEntity(response[0], t)
}

// TODO implement work with local repository
func (c *Client) GetBinaryData(ctx context.Context, id int) error {
	t := model.TypeBinaryData
	response, err := c.localStorage.Get(ctx, t, repository.GetRequest{
		ID: id,
	})

	if err != nil {
		return err
	}

	if len(response) == 0 {
		return fmt.Errorf("no data found")
	}

	return c.showFullEntity(response[0], t)
}

// TODO implement work with local repository
func (c *Client) DeleteLoginPassword(ctx context.Context, id int) error {
	r, err := c.localStorage.Get(ctx, model.TypeUsepass, repository.GetRequest{
		ID: id,
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

		response := c.externalServices.credentials.Delete(ctx, svc.LoginPasswordDeleteRequest{
			ID: r.ExternalID,
		})
		if response.Error != "" {
			errs = append(errs, fmt.Errorf(response.Error))
			continue
		}

		errs = append(errs, c.localStorage.Delete(ctx, model.TypeUsepass, repository.DeleteRequest{
			ID: r.ID,
		}))
	}

	return errors.Join(errs...)
}

// TODO implement work with local repository
func (c *Client) DeleteBankCard(ctx context.Context, id int) error {
	response := c.externalServices.bankCard.Delete(ctx, svc.BankCardDeleteRequest{
		ID: id,
	})
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return c.localStorage.Delete(ctx, model.TypeBankCard, repository.DeleteRequest{
		ID: id,
	})
}

// TODO implement work with local repository
func (c *Client) DeleteTextData(ctx context.Context, id int) error {
	response := c.externalServices.text.Delete(ctx, svc.TextDataDeleteRequest{
		ID: id,
	})
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return c.localStorage.Delete(ctx, model.TypePlainText, repository.DeleteRequest{
		ID: id,
	})
}

// TODO implement work with local repository
func (c *Client) DeleteBinaryData(ctx context.Context, id int) error {
	response := c.externalServices.binary.Delete(ctx, svc.BinaryDataDeleteRequest{
		ID: id,
	})
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	return c.localStorage.Delete(ctx, model.TypeBinaryData, repository.DeleteRequest{
		ID: id,
	})
}
