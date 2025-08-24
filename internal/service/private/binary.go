package private

import (
	"context"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.LoginPasswordService = (*UsepassService)(nil)

type BinaryDataService struct {
	//key encryption key
	kek     encryption.Encoder
	manager model.ModelManager
}

func NewBinaryDataService(manager model.ModelManager, masterKey encryption.Encoder) *BinaryDataService {
	return &BinaryDataService{
		manager: manager,
		kek:     masterKey,
	}
}

func (s *BinaryDataService) List(ctx context.Context, userID int) (service.ListPrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{}, err
	}

	result, err := s.manager.BinaryData.List(ctx, user)

	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{}, err
	}

	if len(result) == 0 {
		//TODO improve message
		return service.ListPrivateDataResponse{}, nil
	}

	response := service.ListPrivateDataResponse{
		Data: make([]map[string]any, len(result)),
	}

	for i, pv := range result {
		response.Data[i] = map[string]any{
			"id":   pv.ID(),
			"name": pv.String(),
		}
	}

	return response, nil
}

func (s *BinaryDataService) Get(ctx context.Context, req service.GetPrivateData) (response service.PrivateDataResponse, err error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.PrivateDataResponse{}, err
	}

	result, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		//TODO change error, and if not found user, message about it
		return service.PrivateDataResponse{}, err
	}

	response.Fields = map[string]any{
		"id":   result.ID(),
		"name": result.String(),
	}

	return
}

func (s *BinaryDataService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		return err
	}

	entity, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		return err
	}
	//TODO wrap error
	return entity.Delete(ctx)
}

func (s *BinaryDataService) Add(ctx context.Context, req service.AddBinaryData, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	entity := s.manager.BinaryData.New(user, req.Name, req.Data)

	return s.save(ctx, entity, clientKey)
}

func (s *BinaryDataService) Update(ctx context.Context, req service.UpdateBinaryData, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	entity, err := s.manager.BinaryData.Get(ctx, user, req.ID)
	if err != nil {
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	entity.SetName(req.Name)
	entity.SetData(req.Data)

	return s.save(ctx, entity, clientKey)
}

func (s *BinaryDataService) save(ctx context.Context, entity *model.BinaryData, clientKey encryption.Encoder) (service.AddingUpdatePrivateDataResponse, error) {
	//generate new key for data
	dek, err := encryption.GenerateNewAESKey()
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	dekSrc := dek.JSON()
	encryptor := encryption.NewModelEncoding([]byte(dekSrc), s.kek, dek)

	if err := entity.Save(ctx, encryptor); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	//replace key to client key
	savedData := entity.EncryptedData()

	encData := encryption.NewEncryptedDataFromReadyData(dek, savedData.Data, savedData.Key)

	if err := encData.ReplaceKey(s.kek, clientKey); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{}, err
	}

	response := service.AddingUpdatePrivateDataResponse{
		ID:   entity.ID(),
		Data: encData.Data,
		Key:  encData.Key,
	}

	return response, nil
}
