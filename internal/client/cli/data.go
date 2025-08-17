package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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
		result, err := c.localStorage.Get(ctx, repository.TypeLoginPassword, repository.GetRequest{ID: id})
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

	if result := c.localStorage.Save(ctx, repository.TypeLoginPassword, repository.SaveRequest{
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

// TODO implement work with local repository
func (c *Client) SaveBankCard(ctx context.Context, number, expires string, cvv int) error {
	date, err := time.Parse("01/2006", expires)
	if err != nil {
		return fmt.Errorf("invalid date format, expected $month/$year, e.g 01/2020")
	}

	data := svc.BankCardSaveRequest{
		Number:  number,
		Expires: date,
		CVV:     cvv,
	}

	response := c.externalServices.bankCard.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, repository.TypeBankCard, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveTextDataAsIs(ctx context.Context, text string, name string) error {
	data := svc.TextDataSaveRequest{
		Text: []byte(text),
		Name: name,
	}

	response := c.externalServices.text.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, repository.TypeTextData, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveTextDataFromFile(ctx context.Context, path string, name string) error {
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

	response := c.externalServices.text.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, repository.TypeTextData, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) SaveBinaryData(ctx context.Context, path string, name string) error {
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

	response := c.externalServices.binary.Save(ctx, data)
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}

	if result := c.localStorage.Save(ctx, repository.TypeBinaryData, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        response.Data.DEK,
		Data:       response.Data.Data,
	}); result.Error != "" {
		return fmt.Errorf(result.Error)
	}
	return nil
}

func (c *Client) ListLoginPassword(ctx context.Context) error {
	if response, err := c.localStorage.Get(ctx, repository.TypeLoginPassword); err == nil {
		c.showList(response)
	} else {
		return err
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) ListBankCard(ctx context.Context) error {
	if response, err := c.localStorage.Get(ctx, repository.TypeBankCard); err == nil {
		c.showList(response)
	} else {
		return err
	}
	return nil

}

// TODO implement work with local repository
func (c *Client) ListTextData(ctx context.Context) error {
	if response, err := c.localStorage.Get(ctx, repository.TypeTextData); err == nil {
		c.showList(response)
	} else {
		return err
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) ListBinaryData(ctx context.Context) error {
	if response, err := c.localStorage.Get(ctx, repository.TypeBinaryData); err == nil {
		c.showList(response)
	} else {
		return err
	}
	return nil
}

// TODO implement work with local repository
func (c *Client) GetLoginPassword(ctx context.Context, id int) error {
	t := repository.TypeLoginPassword
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
	t := repository.TypeBankCard
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
	t := repository.TypeTextData
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
	t := repository.TypeBinaryData
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
	r, err := c.localStorage.Get(ctx, repository.TypeLoginPassword, repository.GetRequest{
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

		errs = append(errs, c.localStorage.Delete(ctx, repository.TypeLoginPassword, repository.DeleteRequest{
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

	return c.localStorage.Delete(ctx, repository.TypeBankCard, repository.DeleteRequest{
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

	return c.localStorage.Delete(ctx, repository.TypeTextData, repository.DeleteRequest{
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

	return c.localStorage.Delete(ctx, repository.TypeBinaryData, repository.DeleteRequest{
		ID: id,
	})
}
