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

func TestClient_DeleteLoginPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockLoginPasswordDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(nil)
		local.EXPECT().Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeUsepass}).Return(nil)
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)
		err := c.DeleteLoginPassword(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("failed external service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockLoginPasswordDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				credentials: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(errors.New("error"))
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeUsepass}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)

		err := c.DeleteLoginPassword(ctx, id)
		assert.Error(t, err)
	})
}

func TestClient_DeleteBankCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBankCardDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				bankCard: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(nil)
		local.EXPECT().Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBankCard}).Return(nil)
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBankCard}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)

		err := c.DeleteBankCard(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("failed external service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBankCardDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				bankCard: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(errors.New("error"))
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBankCard}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)

		err := c.DeleteBankCard(ctx, id)
		assert.Error(t, err)
	})
}

func TestClient_DeleteTextData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockTextDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				text: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(nil)
		local.EXPECT().Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypePlainText}).Return(nil)
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypePlainText}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)

		err := c.DeleteTextData(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("failed external service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockTextDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				text: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypePlainText}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(errors.New("error"))

		err := c.DeleteTextData(ctx, id)
		assert.Error(t, err)
	})
}

func TestClient_DeleteBinaryData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBinaryDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				binary: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}

		ctx := context.Background()
		id := 1

		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(nil)
		local.EXPECT().Delete(ctx, repository.DeleteRequest{ID: id, Type: model.TypeBinaryData}).Return(nil)
		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBinaryData}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)
		err := c.DeleteBinaryData(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("failed external service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := NewMockBinaryDataDataService(ctrl)
		local := NewMockClientRepository(ctrl)

		c := CommandLineClient{
			externalServices: &externalServices{
				binary: svc,
			},
			localStorage: local,
			credential:   []byte("jwt"),
		}
		ctx := context.Background()
		id := 1

		local.EXPECT().Get(ctx, repository.GetRequest{ID: id, Type: model.TypeBinaryData}).Return([]repository.GetResponse{
			{
				ID:         1,
				ExternalID: 1,
			},
		}, nil)
		svc.EXPECT().Delete(ctx, client.DeleteRequest{ID: 1, JWT: string(c.credential)}).Return(errors.New("error"))

		err := c.DeleteBinaryData(ctx, id)
		assert.Error(t, err)
	})
}
