package cli

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	config "github.com/vilasle/gokeep/internal/client"
	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/model"
	repository "github.com/vilasle/gokeep/internal/repository/client"
)

func TestClient_GetLoginPassword(t *testing.T) {
	t.Run("getting specific entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("login\npassword"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "test login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetLoginPassword(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("getting specific entity - invalid data length", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("login"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "test login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting specific entity failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			localStorage: local,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return(nil, errors.New("error"))

		err := c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting specific entity return empty result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			localStorage: local,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{}, nil)

		err := c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting all entities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			localStorage: local,
			// encoder:      encoder,
		}
		id := 0
		ctx := context.Background()

		local.EXPECT().
			All(ctx, model.TypeUsepass).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					// DEK:        []byte(dekHex),
					// Data:       []byte(dataHex),
					View: "Login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err := c.GetLoginPassword(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("getting all entities failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			localStorage: local,
			// encoder:      encoder,
		}
		id := 0
		ctx := context.Background()

		local.EXPECT().
			All(ctx, model.TypeUsepass).
			Return(nil, errors.New("error"))

		err := c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting specific entity - decrypt dek failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("login\npassword"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(nil, errors.New("error"))

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "test login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting specific entity - invalid dek", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`key898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte("data"),
					View:       "test login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err := c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

	t.Run("getting specific entity - invalid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte("data"),
					View:       "test login",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err := c.GetLoginPassword(ctx, id)
		assert.Error(t, err)
	})

}

func TestClient_GetBankCard(t *testing.T) {
	t.Run("getting specific entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("1234567890\n123\n03/31"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBankCard}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "card number",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetBankCard(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("getting specific entity - invalid length", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("1234567890\n123"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBankCard}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "card number",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetBankCard(ctx, id)
		assert.Error(t, err)
	})
}

func TestClient_GetTextData(t *testing.T) {
	t.Run("getting specific entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pwd, err := os.Getwd()
		require.NoError(t, err)

		defer os.Remove(filepath.Join(pwd, "data.txt"))

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("some test"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
			workspace: config.WorkplaceConfig{
				UploadDirectory: config.PathInfo{
					Path: pwd,
				},
			},
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypePlainText}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "data.txt",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetTextData(ctx, id)
		assert.NoError(t, err)
	})
}

func TestClient_GetBinaryData(t *testing.T) {
	t.Run("getting specific entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pwd, err := os.Getwd()
		require.NoError(t, err)

		defer os.Remove(filepath.Join(pwd, "data.txt"))

		local := NewMockClientRepository(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		enc, err := encryption.NewAESKeyFromJSON(jsonKey)
		require.NoError(t, err)

		data, err := enc.Encrypt([]byte("some test"))
		dataHex := hex.EncodeToString(data)
		require.NoError(t, err)

		encoder := NewMockEncoder(ctrl)
		encoder.EXPECT().Decrypt(jsonKey).Return(jsonKey, nil)

		c := CommandLineClient{
			localStorage: local,
			encoder:      encoder,
			workspace: config.WorkplaceConfig{
				UploadDirectory: config.PathInfo{
					Path: pwd,
				},
			},
		}
		id := 1
		ctx := context.Background()

		local.EXPECT().
			Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBinaryData}).
			Return([]repository.GetResponse{
				{
					ID:         id,
					ExternalID: 1,
					DEK:        []byte(dekHex),
					Data:       []byte(dataHex),
					View:       "data.txt",
					Metadata: []repository.MetadataValue{
						{
							Key:   "Key",
							Value: "Value",
						},
					},
				},
			}, nil)

		err = c.GetBinaryData(ctx, id)
		assert.NoError(t, err)
	})
}
