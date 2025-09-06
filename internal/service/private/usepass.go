package private

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.LoginPasswordService = (*UsepassService)(nil)

// UsepassService - is object for work LoginPassword entity
type UsepassService struct {
	//key encryption key
	kek encryption.Encoder
	//manager provide assess to work with models
	manager model.ModelManager
}

// NewUsepassService - create new UsepassService object
func NewUsepassService(manager model.ModelManager, masterKey encryption.Encoder) *UsepassService {
	return &UsepassService{
		manager: manager,
		kek:     masterKey,
	}
}

// List - get list of data by user, and replace KEK to clientKey
func (s *UsepassService) List(ctx context.Context, userID int, clientKey encryption.Encoder) (service.ListPrivateDataResponse, error) {
	log := logger.With("operation", "UsepassService.List")

	log.Info("getting list of data by user", "userId", userID)

	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.ListPrivateDataResponse{}, err
	}

	result, err := s.manager.Usepass.List(ctx, user)

	if err != nil {
		log.Error("failed to get list of data", "error", err)
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
	keys := replacementKeys{kek: s.kek, newKek: clientKey}
	return prepareListOfPrivateData(ls, keys)
}

// Get - get data by id, and replace KEK to clientKey
func (s *UsepassService) Get(ctx context.Context, req service.GetPrivateData, clientKey encryption.Encoder) (response service.PrivateDataResponse, err error) {
	log := logger.With("operation", "UsepassService.Get")
	log.Info("getting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.PrivateDataResponse{}, err
	}

	result, err := s.manager.Usepass.Get(ctx, user, req.ID)
	if err != nil {
		log.Error("failed to get data", "error", err)
		return service.PrivateDataResponse{}, err
	}

	keys := replacementKeys{
		kek:    s.kek,
		newKek: clientKey,
	}
	return getResponseFromModelWithEncryptedDEK(result, keys)
}

// Delete - delete data by id
func (s *UsepassService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	log := logger.With("operation", "UsepassService.Delete")
	log.Info("deleting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return err
	}

	entity, err := s.manager.Usepass.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return err
	}
	return entity.Delete(ctx)
}

// Add - add new data to repository and replace KEK to clientKey
func (s *UsepassService) Add(ctx context.Context, req service.AddLoginPassword, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {
	log := logger.With("operation", "UsepassService.Add")

	log.Info("add new user's entity", "userId", req.UserID)
	log.Debug("private data", "login", req.Username, "password", req.Password)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.PrivateDataResponse{}, err
	}

	entity := s.manager.Usepass.New(user, req.Username, req.Password)
	for _, v := range req.Metadata {
		entity.AddMetadata(v.Key, v.Value)
	}

	return s.save(ctx, entity, clientKey)
}

// Update - update existed entity in repository, and replace KEK to clientKey
func (s *UsepassService) Update(ctx context.Context, req service.UpdateLoginPassword, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {
	log := logger.With("operation", "UsepassService.Update")

	log.Info("update existed user's entity", "userId", req.UserID)
	log.Debug("private data", "login", req.Username, "password", req.Password, "id", req.ID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity, err := s.manager.Usepass.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return service.PrivateDataResponse{}, err
	}

	entity.SetUsername(req.Username)
	entity.SetPassword(req.Password)
	for _, v := range req.Metadata {
		entity.AddMetadata(v.Key, v.Value)
	}

	return s.save(ctx, entity, clientKey)
}

// save - save entity to repository and replace KEK to clientKey
func (s *UsepassService) save(ctx context.Context, entity *model.Usepass, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {
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
