package private

import (
	"context"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.BandCardService = (*BankCardService)(nil)

type BankCardService struct {
	//key encryption key
	kek     encryption.Encoder
	manager model.ModelManager
}

func NewBankCardService(manager model.ModelManager) *BankCardService {
	return &BankCardService{
		manager: manager,
	}
}

func (s *BankCardService) List(ctx context.Context, userID int64) (service.ListPrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "user not found"}, err
	}

	result, err := s.manager.BankCards.List(ctx, user)

	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "error listing bank card"}, err
	}

	if len(result) == 0 {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "there are not any bank cards"}, nil
	}

	response := service.ListPrivateDataResponse{
		Data: make([]map[string]any, len(result)),
	}

	for i, pv := range result {
		response.Data[i] = map[string]any{
			"id":   pv.ID(),
			"card": pv.String(),
		}
	}

	return response, nil
}

func (s *BankCardService) Get(ctx context.Context, req service.GetPrivateData) (response service.PrivateDataResponse, err error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.PrivateDataResponse{Error: "user not found"}, err
	}

	result, err := s.manager.BankCards.Get(ctx, user, req.ID)
	if err != nil {
		//TODO change error, and if not found user, message about it
		return service.PrivateDataResponse{Error: "error getting bank card"}, err
	}

	response.Fields = map[string]any{
		"id":   result.ID(),
		"name": result.String(),
	}

	return
}

func (s *BankCardService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		return err
	}

	entity, err := s.manager.BankCards.Get(ctx, user, req.ID)
	if err != nil {
		return err
	}
	//TODO wrap error
	return entity.Delete(ctx)
}

func (s *BankCardService) Add(ctx context.Context, req service.AddBankCard, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "user not found"}, err
	}
	entity := s.manager.BankCards.New(user, req.Number, req.CVV, req.Expiration)

	return s.save(ctx, entity, clientKey)
}

func (s *BankCardService) Update(ctx context.Context, req service.UpdateBankCard, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "user not found"}, err
	}

	entity, err := s.manager.BankCards.Get(ctx, user, req.ID)
	if err != nil {
		return service.AddingUpdatePrivateDataResponse{Error: "bank card not found"}, err
	}

	entity.SetNumber(req.Number)
	entity.SetCVV(req.CVV)
	entity.SetExpiration(req.Expiration)

	return s.save(ctx, entity, clientKey)
}

func (s *BankCardService) save(ctx context.Context, entity *model.BankCard, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	//generate new key for data
	dek, err := encryption.GenerateNewAESKey()
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error generating new AES key"}, err
	}

	dekSrc := dek.JSON()
	encryptor := encryption.NewEncryptionModel([]byte(dekSrc), s.kek, dek)

	if err := entity.Save(ctx, encryptor); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error saving bank card"}, err
	}

	//replace key to client key
	savedData := entity.EncryptedData()

	encData := encryption.NewEncryptedDataFromReadyData(dek, savedData.Data, savedData.Key)

	if err := encData.ReplaceKey(s.kek, clientKey); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error replacing key"}, err
	}

	response := service.AddingUpdatePrivateDataResponse{
		ID:   entity.ID(),
		Data: encData.Data,
		Key:  encData.Key,
	}

	return response, nil
}
