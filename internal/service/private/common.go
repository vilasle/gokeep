package private

import (
	"errors"

	"github.com/vilasle/gokeep/internal/encryption"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

type replacementKeys struct {
	dek, kek, newKek encryption.Encoder
}

type encryptedData struct {
	data, key []byte
}

func prepareListOfPrivateData(ls []model.PrivateData, keys replacementKeys) (service.ListPrivateDataResponse, error) {
	response := service.ListPrivateDataResponse{
		Data: make([]service.PrivateDataResponse, 0, len(ls)),
	}

	errs := make([]error, 0, len(ls))
	for _, pv := range ls {
		keys := replacementKeys{
			kek:    keys.kek,
			newKek: keys.newKek,
		}

		r, err := getResponseFromModelWithEncryptedDEK(pv, keys)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		response.Data = append(response.Data, r)
	}

	return response, errors.Join(errs...)
}

func getResponseFromModelWithEncryptedDEK(m model.PrivateData, keys replacementKeys) (service.PrivateDataResponse, error) {
	savedData := m.EncryptedData()

	encDEK := savedData.Key

	cleanDEK, err := keys.kek.Decrypt(encDEK)
	if err != nil {
		return service.PrivateDataResponse{},
			errors.Join(ErrDecryptionDEK, err)
	}

	dek, err := encryption.NewAESKeyFromJSON(cleanDEK)
	if err != nil {
		return service.PrivateDataResponse{},
			errors.Join(ErrCreateDEK, err)
	}

	keys.dek = dek

	return getResponseFromModel(m, keys)
}

// getResponseFromModel - replace KEK to client key
func getResponseFromModel(m model.PrivateData, keys replacementKeys) (service.PrivateDataResponse, error) {
	meta := make([]service.MetadataValue, 0, len(m.Metadata()))
	for k, v := range m.Metadata() {
		meta = append(meta, service.MetadataValue{Key: k, Value: v})
	}

	saveData := m.EncryptedData()

	ed := encryptedData{data: saveData.Data, key: saveData.Key}

	newData, err := replaceKey(ed, keys)

	return service.PrivateDataResponse{
		ID:       m.ID(),
		View:     m.String(),
		Data:     string(newData.data),
		Key:      string(newData.key),
		Metadata: meta,
	}, err
}

// replaceKey - replace key in encrypted data
func replaceKey(ed encryptedData, keys replacementKeys) (encryptedData, error) {
	encData := encryption.NewEncryptedDataFromReadyData(keys.dek, ed.data, ed.key)

	if err := encData.ReplaceKey(keys.kek, keys.newKek); err != nil {
		//TODO improve message
		return encryptedData{},
			errors.Join(ErrReplaceKey, err)
	}

	return encryptedData{
		data: []byte(encData.Data),
		key:  []byte(encData.Key),
	}, nil
}
