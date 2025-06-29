package private

import (
	"context"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

var _ service.TextDataService = (*TextService)(nil)

type TextService struct {
	//key encryption key
	kek     encryption.Encryptor
	manager model.ModelManager
}

func NewTextService(manager model.ModelManager) *TextService {
	return &TextService{
		manager: manager,
	}
}

func (s *TextService) List(ctx context.Context, userID int64) (service.ListPrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, userID)
	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "user not found"}, err
	}

	result, err := s.manager.PlainTexts.List(ctx, user)

	if err != nil {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "error listing text data"}, err
	}

	if len(result) == 0 {
		//TODO improve message
		return service.ListPrivateDataResponse{Error: "there are not any text data"}, nil
	}

	response := service.ListPrivateDataResponse{
		Data: make([]map[string]any, len(result)),
	}

	for i, pv := range result {
		response.Data[i] = map[string]any{
			"id":   pv.ID(),
			"text": pv.String(),
		}
	}

	return response, nil
}

func (s *TextService) Get(ctx context.Context, req service.GetPrivateData) (response service.PrivateDataResponse, err error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.PrivateDataResponse{Error: "user not found"}, err
	}

	result, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
	if err != nil {
		//TODO change error, and if not found user, message about it
		return service.PrivateDataResponse{Error: "error getting text data"}, err
	}

	response.Fields = map[string]any{
		"id":   result.ID(),
		"text": result.String(),
	}

	return
}

func (s *TextService) Delete(ctx context.Context, req service.DeletePrivateData) error {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		return err
	}

	entity, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
	if err != nil {
		return err
	}
	//TODO wrap error
	return entity.Delete(ctx)
}

func (s *TextService) Add(ctx context.Context, req service.AddTextData, clientKey encryption.Encryptor) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "user not found"}, err
	}

	entity := s.manager.PlainTexts.New(user, req.Text)

	return s.save(ctx, entity, clientKey)
}

func (s *TextService) Update(ctx context.Context, req service.UpdateTextData, clientKey encryption.Encryptor) (service.AddingUpdatePrivateDataResponse, error) {
	user, err := s.manager.Users.Get(ctx, req.UserID)
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "user not found"}, err
	}

	entity, err := s.manager.PlainTexts.Get(ctx, user, req.ID)
	if err != nil {
		return service.AddingUpdatePrivateDataResponse{Error: "text data not found"}, err
	}

	entity.SetText(req.Text)

	return s.save(ctx, entity, clientKey)
}

func (s *TextService) save(ctx context.Context, usepass *model.PlainText, clientKey encryption.Encryptor) (service.AddingUpdatePrivateDataResponse, error) {
	//generate new key for data
	dek, err := encryption.GenerateNewAESKey()
	if err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error generating new AES key"}, err
	}

	dekSrc := dek.JSON()
	encryptor := encryption.NewEncryptionModel([]byte(dekSrc), s.kek, dek)

	if err := usepass.Save(ctx, encryptor); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error saving text data"}, err
	}

	//replace key to client key
	savedData := usepass.EncryptedData()

	encData := encryption.NewEncryptedDataFromReadyData(dek, savedData.Data, savedData.Key)

	if err := encData.ReplaceKey(s.kek, clientKey); err != nil {
		//TODO improve message
		return service.AddingUpdatePrivateDataResponse{Error: "error replacing key"}, err
	}

	response := service.AddingUpdatePrivateDataResponse{
		ID:   usepass.ID(),
		Data: encData.Data,
		Key:  encData.Key,
	}

	return response, nil
}
