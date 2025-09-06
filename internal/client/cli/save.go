package cli

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	svc "github.com/vilasle/gokeep/internal/service/client"
)

//SaveLoginPassword - save login password to external storage and local storage
func (c *CommandLineClient) SaveLoginPassword(ctx context.Context,
	login, password string, id int, meta map[string]string) (err error) {

	data := svc.LoginPasswordSaveRequest{
		ID:       c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}),
		Login:    login,
		Password: password,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	response, err := c.externalServices.credentials.Save(ctx, data)

	return c.handleSaveResponse(ctx, id, model.TypeUsepass, response, err)
}

//SaveBankCard - save band card to external storage and local storage
func (c *CommandLineClient) SaveBankCard(ctx context.Context,
	number, expires string, cvv, id int, meta map[string]string) (err error) {

	data := svc.BankCardSaveRequest{
		ID:       c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeBankCard}),
		Number:   number,
		Expires:  expires,
		CVV:      cvv,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	response, err := c.externalServices.bankCard.Save(ctx, data)

	return c.handleSaveResponse(ctx, id, model.TypeBankCard, response, err)
}

//SaveTextData - save text data to external storage and local storage
func (c *CommandLineClient) SaveTextData(ctx context.Context,
	text string, name string, id int, meta map[string]string) (err error) {

	data := svc.TextDataSaveRequest{
		ID:       c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypePlainText}),
		Text:     []byte(text),
		Name:     name,
		JWT:      string(c.credential),
		Metadata: prepareMetadataForExternalStorage(meta),
	}

	response, err := c.externalServices.text.Save(ctx, data)

	return c.handleSaveResponse(ctx, id, model.TypePlainText, response, err)
}

//SaveBinaryData - save binary data to external storage and local storage
func (c *CommandLineClient) SaveBinaryData(ctx context.Context,
	content []byte, name string, id int, meta map[string]string) error {

	metadata := prepareMetadataForExternalStorage(meta)
	data := svc.BinaryDataSaveRequest{
		ID:       c.getExternalID(ctx, repository.GetRequest{ID: id, Type: model.TypeBinaryData}),
		Data:     content,
		Name:     name,
		JWT:      string(c.credential),
		Metadata: metadata,
	}

	response, err := c.externalServices.binary.Save(ctx, data)
	return c.handleSaveResponse(ctx, id, model.TypeBinaryData, response, err)
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

func (c *CommandLineClient) handleSaveResponse(ctx context.Context, id int, t model.Type, response svc.SaveResponse, err error) error {
	if err != nil {
		return err
	}

	return c.localStorage.Save(ctx, repository.SaveRequest{
		ID:         id,
		ExternalID: response.ID,
		DEK:        string(response.Data.DEK),
		Data:       string(response.Data.Data),
		View:       response.Data.View,
		Metadata:   prepareMetadataForLocalStorage(response.Metadata),
		Type:       t,
	})
}

