package private

import (
	"context"
	"errors"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/logger"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.BankCardService = (*BankCardService)(nil)

type BankCardService struct {
	//key encryption key
	kek     encryption.Encoder
	manager model.ModelManager
}

func NewBankCardService(manager model.ModelManager, masterKey encryption.Encoder) *BankCardService {
	return &BankCardService{
		manager: manager,
		kek:     masterKey,
	}
}

func (s *BankCardService) List(ctx context.Context, 
	userID int, clientKey encryption.Encoder) (service.ListPrivateDataResponse, error) {

	log := logger.With("operation", "BankCardService.List")

	log.Info("getting list of data by user", "userId", userID)

	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.ListPrivateDataResponse{}, err
	}

	result, err := s.manager.BankCards.List(ctx, user)

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

func (s *BankCardService) Get(ctx context.Context,
	req service.GetPrivateData, clientKey encryption.Encoder) (response service.PrivateDataResponse, err error) {

	log := logger.With("operation", "BankCardService.Get")
	log.Info("getting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	result, err := s.manager.BankCards.Get(ctx, user, req.ID)
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

func (s *BankCardService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	log := logger.With("operation", "BankCardService.Delete")
	log.Info("deleting data by id", "id", req.ID, "userId", req.UserID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return err
	}

	entity, err := s.manager.BankCards.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return err
	}
	return entity.Delete(ctx)
}

func (s *BankCardService) Add(ctx context.Context, 
	req service.AddBankCard, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

	log := logger.With("operation", "BankCardService.Add")

	log.Info("add new user's entity", "userId", req.UserID)
	log.Debug("private data", "number", req.Number, "cvv", req.CVV, "expiration", req.Expiration)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity := s.manager.BankCards.New(user, req.Number, req.CVV, req.Expiration)

	return s.save(ctx, entity, clientKey)
}

func (s *BankCardService) Update(ctx context.Context, 
	req service.UpdateBankCard, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {
		
	log := logger.With("operation", "BankCardService.Update")

	log.Info("update existed user's entity", "userId", req.UserID)
	log.Debug("private data", "number", req.Number, "cvv", req.CVV, "expiration", req.Expiration, "id", req.ID)

	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		logger.Error("failed to get user", err)
		return service.PrivateDataResponse{}, err
	}

	entity, err := s.manager.BankCards.Get(ctx, user, req.ID)
	if err != nil {
		logger.Error("failed to get entity", err)
		return service.PrivateDataResponse{}, err
	}

	entity.SetNumber(req.Number)
	entity.SetCVV(req.CVV)
	entity.SetExpiration(req.Expiration)

	return s.save(ctx, entity, clientKey)
}

func (s *BankCardService) save(ctx context.Context,
	entity *model.BankCard, clientKey encryption.Encoder) (service.PrivateDataResponse, error) {

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
