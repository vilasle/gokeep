package encryption

import "github.com/vilasle/gokeep/internal/model"

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

	//encrypt data with DEK
	if err := ed.Encrypt(data); err != nil {
		return nil, err
	}

	//encrypt DEK key with master key
	if err := ed.EncryptKey(em.masterKey, em.dekSrc); err != nil {
		return nil, err
	}

	return &model.EncryptedData{
		Data: []byte(ed.Data),
		Key:  []byte(ed.Key),
	}, nil

}

func (em EncryptionModel) Decrypt(data model.EncryptedData) ([]byte, error) {
	ed := &EncryptedData{
		dek:  em.dek,
		Data: string(data.Data),
	}

	//decrypt data with DEK
	return ed.Decrypt()

}

/*

 */

// model interface
// type Encrypter interface {
// 	Encrypt([]byte) (*EncryptedData, error)
// 	Decrypt(EncryptedData) ([]byte, error)
// }

// func (tp *BinaryData) Save(ctx context.Context, encrypter Encrypter) (err error) {
// 	if err := tp.prepareEncryptedData(encrypter); err != nil {
// 		//TODO wrap error with package error
// 		return err
// 	}
// 	return tp.Model.Save(ctx)
// }

// func (tp *BinaryData) decryptData(encrypter Encrypter) (err error) {
// 	tp.data, err = encrypter.Decrypt(*tp.encryptedData)
// 	return err
// }
