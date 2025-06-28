package encryption

import (
	"errors"

	"github.com/vilasle/gokeep/internal/model"
)

type EncryptionModel struct {
	masterKey Encryptor
	dek       Encryptor
	dekSrc    []byte
}

func NewEncryptionModel(dekSrc []byte, masterKey, dek Encryptor) *EncryptionModel {
	return &EncryptionModel{
		masterKey: masterKey,
		dek:       dek,
		dekSrc:    dekSrc,
	}
}

func (em EncryptionModel) Encrypt(data []byte) (*model.EncryptedData, error) {
	ed := &EncryptedData{
		dek: em.dek,
	}
	errs := make([]error, 0, 2)
	//encrypt data with DEK
	errs = append(errs, ed.Encrypt(data))
	//encrypt DEK key with master key
	errs = append(errs, ed.EncryptKey(em.masterKey, em.dekSrc))
	return &model.EncryptedData{
		Data: []byte(ed.Data),
		Key:  []byte(ed.Key),
	}, errors.Join(errs...)
}

// Decrypt - decrypts data with DEK
func (em EncryptionModel) Decrypt(data model.EncryptedData) ([]byte, error) {
	return (&EncryptedData{dek: em.dek, Data: string(data.Data)}).Decrypt()
}
