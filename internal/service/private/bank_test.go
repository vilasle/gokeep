package private

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

func TestBankCardService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)
		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)
		clientKey.EXPECT().Encrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{
			{
				ID:     1,
				UserID: 1,
				Type:   model.TypeBankCard,
				View:   "card",
				Data:   []byte("data"),
				DEK:    []byte(dekHex),
				Metadata: map[string]string{
					"key": "value",
				},
			},
		}, nil)

		expectedKey := "7b226b6579223a2238393861313637306534326266316562363339616331303538396134386334346638313166613461643561656564306534333630376430643735623362633034222c226e6f6e6365223a22393366376266653335623934343266633136373937383565227d"
		expected := service.ListPrivateDataResponse{
			Data: []service.PrivateDataResponse{
				{
					ID:   1,
					View: "card",
					Data: "data",
					Key:  expectedKey,
					Metadata: []service.MetadataValue{
						{
							Key:   "key",
							Value: "value",
						},
					},
				},
			},
		}

		resp, err := svc.List(ctx, userID, clientKey)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("getting user failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, errors.New("error"))

		resp, err := svc.List(ctx, userID, clientKey)
		assert.Error(t, err)
		assert.Empty(t, resp)
	})

	t.Run("getting list of data failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{
			{
				ID:     1,
				UserID: 1,
				Type:   model.TypeBankCard,
				View:   "card",
				Data:   []byte("data"),
				DEK:    []byte(dekHex),
				Metadata: map[string]string{
					"key": "value",
				},
			},
		}, nil).Return(nil, errors.New("error"))

		resp, err := svc.List(ctx, userID, clientKey)
		assert.Error(t, err)
		assert.Empty(t, resp)
	})

	t.Run("empty list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{}, nil)

		resp, err := svc.List(ctx, userID, clientKey)
		assert.NoError(t, err)
		assert.Empty(t, resp)
	})

	t.Run("error decryption of DEK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return(nil, errors.New("error"))

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{
			{
				ID:     1,
				UserID: 1,
				Type:   model.TypeBankCard,
				View:   "card",
				Data:   []byte("data"),
				DEK:    []byte(dekHex),
				Metadata: map[string]string{
					"key": "value",
				},
			},
		}, nil)

		_, err := svc.List(ctx, userID, clientKey)
		assert.Error(t, err)
	})

	t.Run("replacing keys failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return(jsonKey, nil)
		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return(jsonKey, nil)
		clientKey.EXPECT().Encrypt([]byte(jsonKey)).Return(nil, errors.New("errors"))

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{
			{
				ID:     1,
				UserID: 1,
				Type:   model.TypeBankCard,
				View:   "card",
				Data:   []byte("data"),
				DEK:    []byte(dekHex),
				Metadata: map[string]string{
					"key": "value",
				},
			},
		}, nil)

		_, err := svc.List(ctx, userID, clientKey)
		assert.Error(t, err)
	})

	t.Run("wrong DEK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"93f7bfe35b9442fc1679785e"`)
		dekHex := hex.EncodeToString(jsonKey)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return(jsonKey, nil)

		private.EXPECT().List(ctx, model.TypeBankCard, 1).Return([]model.PrivateDataInfo{
			{
				ID:     1,
				UserID: 1,
				Type:   model.TypeBankCard,
				View:   "card",
				Data:   []byte("data"),
				DEK:    []byte(dekHex),
				Metadata: map[string]string{
					"key": "value",
				},
			},
		}, nil)

		_, err := svc.List(ctx, userID, clientKey)
		assert.Error(t, err)
	})
}

func TestBankCardService_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		dekHex := hex.EncodeToString(jsonKey)
		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)
		masterKey.EXPECT().Decrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)
		clientKey.EXPECT().Encrypt([]byte(jsonKey)).Return([]byte(jsonKey), nil)

		private.EXPECT().Get(ctx, id).Return(model.PrivateDataInfo{
			ID:     1,
			UserID: 1,
			Type:   model.TypeBankCard,
			View:   "card",
			Data:   []byte("data"),
			DEK:    []byte(dekHex),
			Metadata: map[string]string{
				"key": "value",
			},
		}, nil)

		expectedKey := "7b226b6579223a2238393861313637306534326266316562363339616331303538396134386334346638313166613461643561656564306534333630376430643735623362633034222c226e6f6e6365223a22393366376266653335623934343266633136373937383565227d"
		expected := service.PrivateDataResponse{
			ID:   1,
			View: "card",
			Data: "data",
			Key:  expectedKey,
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}

		resp, err := svc.Get(ctx, service.GetPrivateData{
			ID:     id,
			UserID: userID,
		}, clientKey)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})
	t.Run("getting user is failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, errors.New("error"))

		_, err := svc.Get(ctx, service.GetPrivateData{
			ID:     id,
			UserID: userID,
		}, clientKey)
		assert.Error(t, err)
	})

	t.Run("getting entity failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().Get(ctx, id).Return(model.PrivateDataInfo{}, errors.New("error"))

		_, err := svc.Get(ctx, service.GetPrivateData{
			ID:     id,
			UserID: userID,
		}, clientKey)
		assert.Error(t, err)
	})
}

func TestBankCardService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1

		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().Get(ctx, id).Return(model.PrivateDataInfo{
			ID:     1,
			UserID: 1,
			Type:   model.TypeBankCard,
			View:   "card",
			Data:   []byte("data"),
			DEK:    []byte("dek"),
			Metadata: map[string]string{
				"key": "value",
			},
		}, nil)

		private.EXPECT().Delete(ctx, id).Return(nil)

		err := svc.Delete(ctx, service.DeletePrivateData{
			ID:     id,
			UserID: userID,
		})
		assert.NoError(t, err)
	})

	t.Run("getting user failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1

		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, errors.New("error"))

		err := svc.Delete(ctx, service.DeletePrivateData{
			ID:     id,
			UserID: userID,
		})
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1

		id := 1
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().Get(ctx, id).Return(model.PrivateDataInfo{}, errors.New("error"))

		err := svc.Delete(ctx, service.DeletePrivateData{
			ID:     id,
			UserID: userID,
		})
		assert.Error(t, err)
	})
}

func TestBankCardService_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Encrypt(gomock.Any()).Return([]byte(jsonKey), nil)
		masterKey.EXPECT().Decrypt(gomock.Any()).Return([]byte(jsonKey), nil)
		clientKey.EXPECT().Encrypt(gomock.Any()).Return([]byte(jsonKey), nil)

		private.EXPECT().Add(ctx, gomock.Any()).Return(1, nil)

		_, err := svc.Add(ctx, service.AddBankCard{
			Number: "123456890",
			CVV:    123,
			UserID: userID,
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.NoError(t, err)
	})

	t.Run("failed getting user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, errors.New("error"))

		_, err := svc.Add(ctx, service.AddBankCard{
			Number: "123456890",
			CVV:    123,
			UserID: userID,
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.Error(t, err)
	})

	t.Run("saving failed entity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Encrypt(gomock.Any()).Return([]byte(jsonKey), nil)

		private.EXPECT().Add(ctx, gomock.Any()).Return(0, errors.New("error"))

		_, err := svc.Add(ctx, service.AddBankCard{
			Number: "123456890",
			CVV:    123,
			UserID: userID,
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.Error(t, err)
	})
}

func TestBankCardService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		jsonKey := []byte(`{"key":"898a1670e42bf1eb639ac10589a48c44f811fa4ad5aeed0e43607d0d75b3bc04","nonce":"93f7bfe35b9442fc1679785e"}`)
		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		masterKey.EXPECT().Encrypt(gomock.Any()).Return([]byte(jsonKey), nil)
		masterKey.EXPECT().Decrypt(gomock.Any()).Return([]byte(jsonKey), nil)
		clientKey.EXPECT().Encrypt(gomock.Any()).Return([]byte(jsonKey), nil)

		private.EXPECT().Get(ctx, 1).Return(model.PrivateDataInfo{
			ID:     1,
			UserID: 1,
			Type:   model.TypeBankCard,
			View:   "card",
			Data:   []byte("data"),
			DEK:    []byte("dek"),
			Metadata: map[string]string{
				"key": "value",
			},
		}, nil)
		private.EXPECT().Update(ctx, gomock.Any()).Return(1, nil)

		_, err := svc.Update(ctx, service.UpdateBankCard{
			ID: 1,
			AddBankCard: service.AddBankCard{
				Number: "123456890",
				CVV:    123,
				UserID: userID,
			},
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.NoError(t, err)
	})

	t.Run("getting user failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{}, errors.New("error"))

		_, err := svc.Update(ctx, service.UpdateBankCard{
			ID: 1,
			AddBankCard: service.AddBankCard{
				Number: "123456890",
				CVV:    123,
				UserID: userID,
			},
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		private := NewMockPrivateDataRepository(ctrl)
		user := NewMockUserRepository(ctrl)

		masterKey := NewMockEncoder(ctrl)

		collector := NewMockRepositoryCollector(ctrl)
		collector.EXPECT().User().Return(user)
		collector.EXPECT().Private().Return(private)

		manager := model.NewModelManager(collector)

		svc := NewBankCardService(manager, masterKey)

		ctx := context.Background()
		userID := 1
		clientKey := NewMockEncoder(ctrl)

		user.EXPECT().Get(ctx, userID).Return(model.UserInfo{
			ID:       1,
			Login:    "login",
			Password: "password",
		}, nil)

		private.EXPECT().Get(ctx, 1).Return(model.PrivateDataInfo{}, errors.New("error"))
		_, err := svc.Update(ctx, service.UpdateBankCard{
			ID: 1,
			AddBankCard: service.AddBankCard{
				Number: "123456890",
				CVV:    123,
				UserID: userID,
			},
			Metadata: []service.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}, clientKey)
		assert.Error(t, err)
	})
}
