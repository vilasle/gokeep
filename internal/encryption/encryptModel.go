package encryption

import (
	"errors"

	"github.com/vilasle/gokeep/internal/model"
)

type ModelEncoding struct {
	kek    Encoder
	dek    Encoder
	dekSrc []byte
}
// NewModelEncoding - creates new ModelEncoding
func NewModelEncoding(dekSrc []byte, kek, dek Encoder) *ModelEncoding {
	return &ModelEncoding{
		kek:    kek,
		dek:    dek,
		dekSrc: dekSrc,
	}
}

//Encrypt - encrypts data with DEK and return EncryptedData for model
func (em ModelEncoding) Encrypt(data []byte) (*model.EncryptedData, error) {
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

// Decrypt - decrypts data with DEK from model data
func (em ModelEncoding) Decrypt(data model.EncryptedData) ([]byte, error) {
	return (&EncryptedData{dek: em.dek, Data: string(data.Data)}).Decrypt()
}
