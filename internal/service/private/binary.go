package private

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.BinaryDataService = (*BinaryDataService)(nil)

type BinaryDataService struct {
	//key encryption key
	kek encryption.Encoder
	//manager provide assess to work with models
	manager model.ModelManager
}

func NewBinaryDataService(manager model.ModelManager, masterKey encryption.Encoder) *BinaryDataService {
	return &BinaryDataService{
		manager: manager,
		kek:     masterKey,
	}
}

func (s *BinaryDataService) List(ctx context.Context,
	userID int, clientKey encryption.Encoder) (service.ListPrivateDataResponse, error) {

	log := logger.With("operation", "BinaryDataService.List")

	log.Info("getting list of data by user", "userId", userID)

	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.ListPrivateDataResponse{}, err
	}

	result, err := s.manager.BinaryData.List(ctx, user)

	if err != nil {
		log.Error("failed to get list of data", err)
		return service.ListPrivateDataResponse{}, err
	}

	if len(result) == 0 {
		log.Debug("result is empty")
		return service.ListPrivateDataResponse{}, nil
	}

	ls := make([]model.PrivateData, len(result))
	for i, v := range result {
		ls[i] = v
	}
	return prepareListOfPrivateData(ls, replacementKeys{kek: s.kek, newKek: clientKey})
}

func (s *BinaryDataService) Get(ctx context.Context,
	req service.GetPrivateData, clientKey encryption.Encoder) (response service.PrivateDataResponse, err error) {

	log := logger.With("operation", "BinaryDataService.Get")
	log.Info("getting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	result, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		log.Error("failed to get data", err)
		return service.PrivateDataResponse{}, err
	}

	keys := replacementKeys{
		kek:    s.kek,
		newKek: clientKey,
	}
	return getResponseFromModelWithEncryptedDEK(result, keys)
}

func (s *BinaryDataService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	log := logger.With("operation", "BinaryDataService.Delete")
	log.Info("deleting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return err
	}

	entity, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return err
	}
	return entity.Delete(ctx)
}

func (s *BinaryDataService) Add(ctx context.Context,
	req service.AddBinaryData, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	log := logger.With("operation", "BinaryDataService.Add")

	log.Info("add new user's entity", "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity := s.manager.BinaryData.New(user, req.Data, req.Name)

	return s.save(ctx, entity, clientKey)
}

func (s *BinaryDataService) Update(ctx context.Context,
	req service.UpdateBinaryData, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	log := logger.With("operation", "BinaryDataService.Update")

	log.Info("update existed user's entity", "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return service.PrivateDataResponse{}, err
	}

	entity.SetData(req.Data)
	entity.SetName(req.Name)

	return s.save(ctx, entity, clientKey)
}

func (s *BinaryDataService) save(ctx context.Context,
	entity *model.BinaryData, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	dek, err := encryption.GenerateNewAESKey()
	if err != nil {
		return service.PrivateDataResponse{},
			errors.Join(ErrGenerateDEK, err)
	}

	dekSrc := dek.JSON()
	encryptor := encryption.NewModelEncoding([]byte(dekSrc), s.kek, dek)

	if err := entity.Save(ctx, encryptor); err != nil {
		return service.PrivateDataResponse{},
			errors.Join(ErrSaveData, err)
	}

	keys := replacementKeys{
		kek:    s.kek,
		dek:    dek,
		newKek: clientKey,
	}

	return getResponseFromModel(entity, keys)
}
