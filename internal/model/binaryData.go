package model

import (
	"context"
	"fmt"
)

var _ PrivateData = (*BinaryData)(nil)

type BinaryData struct {
	Model
	data []byte
}

func newBinaryData(owner *User, data []byte) *BinaryData {
	return &BinaryData{
		Model: Model{
			owner: owner,
		},
		data: data,
	}
}

func (tp *BinaryData) Save(ctx context.Context, encrypter Encrypter) (err error) {
	if err := tp.prepareEncryptedData(encrypter); err != nil {
		//TODO wrap error with package error
		return err
	}
	return tp.Model.Save(ctx)
}

func (tp *BinaryData) prepareEncryptedData(encrypter Encrypter) error {
	encryptedData, err := encrypter.Encrypt(tp.data)
	if err != nil {
		return err
	}
	encryptedData.owner = tp

	tp.encryptedData = encryptedData

	return nil
}

func (tp *BinaryData) decryptData(encrypter Encrypter) (err error) {
	tp.data, err = encrypter.Decrypt(*tp.encryptedData)
	return err
}

func findBinaryDataByID(ctx context.Context, id int64, r PrivateDataRepository) (*BinaryData, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	u, ok := data.(*BinaryData)
	if !ok {
		return nil, fmt.Errorf("invalid data. private data(id=%d) does not have type '*BinaryData'", id)
	}

	return u, nil
}
