package encryption

import (
	"errors"

	"github.com/vilasle/gokeep/internal/model"
)

type EncryptionModel struct {
	kek    Encoder
	dek    Encoder
	dekSrc []byte
}

func NewEncryptionModel(dekSrc []byte, kek, dek Encoder) *EncryptionModel {
	return &EncryptionModel{
		kek:    kek,
		dek:    dek,
		dekSrc: dekSrc,
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
	errs = append(errs, ed.EncryptKey(em.kek, em.dekSrc))
	return &model.EncryptedData{
		Data: []byte(ed.Data),
		Key:  []byte(ed.Key),
	}, errors.Join(errs...)
}

// Decrypt - decrypts data with DEK
func (em EncryptionModel) Decrypt(data model.EncryptedData) ([]byte, error) {
	return (&EncryptedData{dek: em.dek, Data: string(data.Data)}).Decrypt()
}
