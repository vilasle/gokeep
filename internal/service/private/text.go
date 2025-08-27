package private

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.TextDataService = (*TextService)(nil)

type TextService struct {
	//key encryption key
	kek encryption.Encoder
	//manager provide assess to work with models
	manager model.ModelManager
}

func NewTextService(manager model.ModelManager, masterKey encryption.Encoder) *TextService {
	return &TextService{
		manager: manager,
		kek:     masterKey,
	}
}

func (s *TextService) List(ctx context.Context, userID int, clientKey encryption.Encoder) (service.ListPrivateDataResponse, error) {
	log := logger.With("operation", "TextService.List")

	log.Info("getting list of data by user", "userId", userID)

	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.ListPrivateDataResponse{}, err
	}

	result, err := s.manager.PlainTexts.List(ctx, user)

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
	return prepareListOfPrivateData(ls, replacementKeys{kek: s.kek, newKek: clientKey})
}

func (s *TextService) Get(ctx context.Context,
	req service.GetPrivateData, clientKey encryption.Encoder) (response service.PrivateDataResponse, err error) {

	log := logger.With("operation", "TextService.Get")
	log.Info("getting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.PrivateDataResponse{}, err
	}

	result, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
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

func (s *TextService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	log := logger.With("operation", "TextService.Delete")
	log.Info("deleting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return err
	}

	entity, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return err
	}
	return entity.Delete(ctx)
}

func (s *TextService) Add(ctx context.Context,
	req service.AddTextData, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	log := logger.With("operation", "TextService.Add")

	log.Info("add new user's entity", "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return service.PrivateDataResponse{}, err
	}

	entity := s.manager.PlainTexts.New(user, req.Text, req.Name)
	for _, v := range req.Metadata {
		entity.AddMetadata(v.Key, v.Value)
	}

	return s.save(ctx, entity, clientKey)
}

func (s *TextService) Update(ctx context.Context,
	req service.UpdateTextData, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	log := logger.With("operation", "TextService.Update")

	log.Info("update existed user's entity", "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return service.PrivateDataResponse{}, err
	}

	entity.SetText(req.Text)

	for _, v := range req.Metadata {
		entity.AddMetadata(v.Key, v.Value)
	}

	return s.save(ctx, entity, clientKey)
}

func (s *TextService) save(ctx context.Context,
	entity *model.PlainText, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

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
