package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	"github.com/vilasle/gokeep/internal/service/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

func (c *Client) SaveLoginPassword(ctx context.Context,
	login, password string, id int, meta map[string]string) (err error) {

	data := svc.LoginPasswordSaveRequest{
		ID:       id,
		Login:    login,
		Password: password,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	if data.ID, err = c.getExistedId(ctx, model.TypeUsepass, id); err != nil {
		return err
	}

	response, err := c.externalServices.credentials.Save(ctx, data)

	return c.handleSaveResponse(ctx, model.TypeUsepass, response, err)
}

func (c *Client) SaveBankCard(ctx context.Context,
	number, expires string, cvv, id int, meta map[string]string) (err error) {

	data := svc.BankCardSaveRequest{
		Number:   number,
		Expires:  expires,
		CVV:      cvv,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	if data.ID, err = c.getExistedId(ctx, model.TypeBankCard, id); err != nil {
		return err
	}

	response, err := c.externalServices.bankCard.Save(ctx, data)

	return c.handleSaveResponse(ctx, model.TypeBankCard, response, err)
}

func (c *Client) SaveTextDataAsIs(ctx context.Context,
	text string, name string, id int, meta map[string]string) (err error) {

	data := svc.TextDataSaveRequest{
		Text:     []byte(text),
		Name:     name,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	if data.ID, err = c.getExistedId(ctx, model.TypePlainText, id); err != nil {
		return err
	}

	response, err := c.externalServices.text.Save(ctx, data)

	return c.handleSaveResponse(ctx, model.TypePlainText, response, err)
}

func (c *Client) SaveTextDataFromFile(ctx context.Context,
	path string, name string, id int, meta map[string]string) (err error) {

	content, err := getContentFromFile(name, path)
	if err != nil {
		return errors.Join(err, errors.New("failed to get content from file"))
	}

	data := svc.TextDataSaveRequest{
		Text:     content,
		Name:     name,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	if data.ID, err = c.getExistedId(ctx, model.TypePlainText, id); err != nil {
		return err
	}

	response, err := c.externalServices.text.Save(ctx, data)
	return c.handleSaveResponse(ctx, model.TypePlainText, response, err)
}

func (c *Client) SaveBinaryData(ctx context.Context,
	path string, name string, id int, meta map[string]string) error {

	content, err := getContentFromFile(name, path)
	if err != nil {
		return errors.Join(err, errors.New("failed to get content from file"))
	}

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.BinaryDataSaveRequest{
		Data:     content,
		Name:     name,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	if data.ID, err = c.getExistedId(ctx, model.TypeBinaryData, id); err != nil {
		return err
	}

	response, err := c.externalServices.binary.Save(ctx, data)
	return c.handleSaveResponse(ctx, model.TypeBinaryData, response, err)
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

func (c *Client) handleSaveResponse(ctx context.Context, t model.Type, response client.SaveResponse, err error) error {
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, repository.SaveRequest{
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
		Type:       t,
	})
}

func (c *Client) getExistedId(ctx context.Context, tData int, id int) (int, error) {
	if id > 0 {
		result, err := c.localStorage.Get(ctx, repository.GetRequest{ID: id, Type: tData})
		if err == nil && len(result) > 0 {
			return result[0].ExternalID, nil
		} else if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func getContentFromFile(name, path string) ([]byte, error) {
	if name == "" {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		name = stat.Name()
	}

	return os.ReadFile(path)
}
