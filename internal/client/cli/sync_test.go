package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service/client"
)

func TestClient_Sync(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		bank.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBankCard}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		text.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypePlainText}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		binary.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBinaryData}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		localStorage.EXPECT().Rewrite(ctx, gomock.Any()).Return(nil)

		err := c.Sync(ctx)
		assert.NoError(t, err)
	})

	t.Run("cred service failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return(nil, errors.New("error"))

		err := c.Sync(ctx)
		assert.Error(t, err)
	})

	t.Run("bank service is failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		bank.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBankCard}).Return(nil, errors.New("error"))

		err := c.Sync(ctx)
		assert.Error(t, err)
	})

	t.Run("text service failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		bank.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBankCard}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		text.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypePlainText}).Return(nil, errors.New("error"))

		err := c.Sync(ctx)
		assert.Error(t, err)
	})

	t.Run("binary service failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		bank.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBankCard}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		text.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypePlainText}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		binary.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBinaryData}).Return(nil, errors.New("error"))

		err := c.Sync(ctx)
		assert.Error(t, err)
	})

	t.Run("local storage failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cred := NewMockLoginPasswordDataService(ctrl)
		bank := NewMockBankCardDataService(ctrl)
		binary := NewMockBinaryDataDataService(ctrl)
		text := NewMockTextDataDataService(ctrl)

		localStorage := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: cred,
				bankCard:    bank,
				binary:      binary,
				text:        text,
			},
			credential:   []byte("jwt"),
			localStorage: localStorage,
		}

		ctx := context.Background()

		cred.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeUsepass}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: nil,
			},
		}, nil)

		bank.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBankCard}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		text.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypePlainText}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		binary.EXPECT().Get(ctx, client.GetRequest{JWT: "jwt", Type: model.TypeBinaryData}).Return([]client.EncryptedEntity{
			{
				ID: 1,
				Data: client.EncryptedData{
					View: "view",
					DEK:  []byte("dek"),
					Data: []byte("data"),
				},
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}, nil)

		localStorage.EXPECT().Rewrite(ctx, gomock.Any()).Return(errors.New("error"))

		err := c.Sync(ctx)
		assert.Error(t, err)
	})
}
