package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

func (c *Client) SaveLoginPassword(ctx context.Context,
	login, password string, id int, meta map[string]string) error {

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.LoginPasswordSaveRequest{
		ID:       id,
		Login:    login,
		Password: password,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeUsepass, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response, err := c.externalServices.credentials.Save(ctx, data)
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, model.TypeUsepass, repository.SaveRequest{
		ID:         id,
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       login,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
	})
}

func (c *Client) SaveBankCard(ctx context.Context,
	number, expires string, cvv, id int, meta map[string]string) error {

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.BankCardSaveRequest{
		Number:   number,
		Expires:  expires,
		CVV:      cvv,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeBankCard, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response, err := c.externalServices.bankCard.Save(ctx, data)
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, model.TypeBankCard, repository.SaveRequest{
		ID:         id,
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
	})
}

func (c *Client) SaveTextDataAsIs(ctx context.Context,
	text string, name string, id int, meta map[string]string) error {

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.TextDataSaveRequest{
		Text:     []byte(text),
		Name:     name,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypePlainText, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response, err := c.externalServices.text.Save(ctx, data)
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, model.TypePlainText, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
	})
}

func (c *Client) SaveTextDataFromFile(ctx context.Context,
	path string, name string, id int, meta map[string]string) error {

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

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.TextDataSaveRequest{
		Text:     content,
		Name:     name,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypePlainText, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response, err := c.externalServices.text.Save(ctx, data)
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, model.TypePlainText, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
	})
}

func (c *Client) SaveBinaryData(ctx context.Context,
	path string, name string, id int, meta map[string]string) error {

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

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.BinaryDataSaveRequest{
		Data:     content,
		Name:     name,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if id > 0 {
		result, err := c.localStorage.Get(ctx, model.TypeUsepass, repository.GetRequest{ID: id})
		if err == nil && len(result) > 0 {
			data.ID = result[0].ExternalID
		} else if err != nil {
			return err
		}
	}

	response, err := c.externalServices.binary.Save(ctx, data)
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, model.TypeBinaryData, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
	})
}

func (c *Client) GetLoginPassword(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeUsepass, id)
}

// TODO implement work with local repository
func (c *Client) GetBankCard(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeBankCard, id)
}

// TODO implement work with local repository
func (c *Client) GetTextData(ctx context.Context, id int) error {
	return c.get(ctx, model.TypePlainText, id)
}

// TODO implement work with local repository
func (c *Client) GetBinaryData(ctx context.Context, id int) error {
	return c.get(ctx, model.TypeBinaryData, id)
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

		err := c.externalServices.credentials.Delete(ctx, svc.DeleteRequest{
			ID:  r.ExternalID,
			JWT: string(c.credential),
		})
		if err != nil {
			errs = append(errs, err)
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
	if err := c.externalServices.bankCard.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}

	return c.localStorage.Delete(ctx, model.TypeBankCard, repository.DeleteRequest{ID: id})
}

// TODO implement work with local repository
func (c *Client) DeleteTextData(ctx context.Context, id int) error {
	if err := c.externalServices.text.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, model.TypePlainText, repository.DeleteRequest{ID: id})
}

// TODO implement work with local repository
func (c *Client) DeleteBinaryData(ctx context.Context, id int) error {
	if err := c.externalServices.binary.Delete(ctx, svc.DeleteRequest{ID: id, JWT: string(c.credential)}); err != nil {
		return err
	}
	return c.localStorage.Delete(ctx, model.TypeBinaryData, repository.DeleteRequest{ID: id})
}

func (c *Client) get(ctx context.Context, t model.Type, id int) error {
	var (
		response []repository.GetResponse
		err      error
	)

	//show specific data and all entity in short format
	if id > 0 {
		if response, err = c.localStorage.Get(ctx, t, repository.GetRequest{ID: id}); err != nil {
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
		return c.showListOfEntities(t, response...)
	}
}

func prepareMetadataForExternalStorage(meta map[string]string) []svc.MetadataValue {
	if meta == nil {
		return []svc.MetadataValue{}
	}
	metadata := make([]svc.MetadataValue, 0, len(meta))
	for k, v := range meta {
		metadata = append(metadata, svc.MetadataValue{Key: k, Value: v})
	}
	return metadata
}

func prepareMetadataForLocalStorage(meta []svc.MetadataValue) []repository.MetadataValue {
	if meta == nil {
		return []repository.MetadataValue{}
	}
	metadata := make([]repository.MetadataValue, 0, len(meta))
	for _, v := range meta {
		metadata = append(metadata, repository.MetadataValue{Key: v.Key, Value: v.Value})
	}
	return metadata
}
