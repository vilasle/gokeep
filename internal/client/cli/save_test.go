package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
	"github.com/vilasle/gokeep/internal/service/client"
)

func TestClient_SaveLoginPassword(t *testing.T) {
	t.Run("success saving login password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := Client{
			externalServices: &externalServices{
				credentials: cred,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()

		id := 0
		meta := map[string]string{
			"key": "value",
		}
		login, password := "login", "password"

		cred.EXPECT().Save(ctx, client.LoginPasswordSaveRequest{
			ID:       id,
			Login:    login,
			Password: password,
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
		}).Return(client.SaveResponse{
			ID: 10,
			Data: client.EncryptedData{
				DEK:  []byte("dek"),
				Data: []byte("data"),
				View: "view",
			},
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, nil)

		local.EXPECT().Save(ctx, repository.SaveRequest{
			ExternalID: 10,
			DEK:        "dek",
			Data:       "data",
			View:       "view",
			Metadata: []repository.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			Type: model.TypeUsepass,
		}).Return(nil)

		err := c.SaveLoginPassword(ctx, login, password, id, meta)
		assert.NoError(t, err)
	})

	t.Run("saving login password failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := Client{
			externalServices: &externalServices{
				credentials: cred,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()

		id := 0
		meta := map[string]string{
			"key": "value",
		}
		login, password := "login", "password"

		cred.EXPECT().Save(ctx, client.LoginPasswordSaveRequest{
			ID:       id,
			Login:    login,
			Password: password,
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
		}).Return(client.SaveResponse{}, errors.New("error"))

		err := c.SaveLoginPassword(ctx, login, password, id, meta)
		assert.Error(t, err)
	})
}

func TestClient_SaveBankCard(t *testing.T) {
	t.Run("success saving bank card", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBankCardDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := Client{
			externalServices: &externalServices{
				bankCard: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()

		id := 0
		number, expires, cvv := "1234567890", "03/31", 123

		svc.EXPECT().Save(ctx, client.BankCardSaveRequest{
			ID:       id,
			Number:   number,
			CVV:      cvv,
			Expires:  expires,
			Metadata: []client.MetadataValue{},
			JWT:      "jwt",
		}).Return(client.SaveResponse{
			ID: 10,
			Data: client.EncryptedData{
				DEK:  []byte("dek"),
				Data: []byte("data"),
				View: "view",
			},
			Metadata: []client.MetadataValue{},
		}, nil)

		local.EXPECT().Save(ctx, repository.SaveRequest{
			ExternalID: 10,
			DEK:        "dek",
			Data:       "data",
			View:       "view",
			Metadata:   []repository.MetadataValue{},
			Type:       model.TypeBankCard,
		}).Return(nil)

		err := c.SaveBankCard(ctx, number, expires, cvv, id, nil)
		assert.NoError(t, err)
	})
}

func TestClient_SaveTextData(t *testing.T) {
	t.Run("success saving text data as is", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockTextDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := Client{
			externalServices: &externalServices{
				text: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()

		id := 0
		meta := map[string]string{
			"key": "value",
		}
		name, text := "text_name", "some text"

		svc.EXPECT().Save(ctx, client.TextDataSaveRequest{
			ID:   id,
			Name: name,
			Text: []byte(text),
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
		}).Return(client.SaveResponse{
			ID: 10,
			Data: client.EncryptedData{
				DEK:  []byte("dek"),
				Data: []byte("data"),
				View: "view",
			},
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, nil)

		local.EXPECT().Save(ctx, repository.SaveRequest{
			ExternalID: 10,
			DEK:        "dek",
			Data:       "data",
			View:       "view",
			Metadata: []repository.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			Type: model.TypePlainText,
		}).Return(nil)

		err := c.SaveTextData(ctx, text, name, id, meta)
		assert.NoError(t, err)
	})
}

func TestClient_SaveBinaryData(t *testing.T) {
	t.Run("success saving binary data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBinaryDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := Client{
			externalServices: &externalServices{
				binary: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()

		id := 0
		meta := map[string]string{
			"key": "value",
		}
		name := "text_name"
		text := []byte("some text")

		svc.EXPECT().Save(ctx, client.BinaryDataSaveRequest{
			ID:   id,
			Name: name,
			Data: []byte(text),
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			JWT: "jwt",
		}).Return(client.SaveResponse{
			ID: 10,
			Data: client.EncryptedData{
				DEK:  []byte("dek"),
				Data: []byte("data"),
				View: "view",
			},
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, nil)

		local.EXPECT().Save(ctx, repository.SaveRequest{
			ExternalID: 10,
			DEK:        "dek",
			Data:       "data",
			View:       "view",
			Metadata: []repository.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
			Type: model.TypeBinaryData,
		}).Return(nil)

		err := c.SaveBinaryData(ctx, text, name, id, meta)
		assert.NoError(t, err)
	})
}
