package model

import (
	"context"
)

var _ PrivateData = (*BinaryData)(nil)

type BinaryData struct {
	model
	name string
	data []byte
}

func newBinaryData(owner *User, data []byte, name string) *BinaryData {
	return &BinaryData{
		model: model{
			owner:     owner,
			modelType: TypeBinaryData,
		},
		name: name,
		data: data,
	}
}

func (bd *BinaryData) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := bd.prepareEncryptedData(encoder); err != nil {
		//TODO wrap error with package error
		return err
	}
	return bd.model.Save(ctx)
}

func (bd *BinaryData) String() string {
	return bd.name
}

func (bd *BinaryData) SetName(name string) {
	bd.name = name
}

func (bd *BinaryData) SetData(data []byte) {
	bd.data = data
}

func (bd *BinaryData) prepareEncryptedData(encoder Encoder) (err error) {
	bd.encryptedData, err = encoder.Encrypt(bd.data)
	return err
}

func (bd *BinaryData) decryptData(encoder Encoder) (err error) {
	bd.data, err = encoder.Decrypt(*bd.encryptedData)
	return err
}

// FIXME add getting binary data by id and check that owner was right id
func findBinaryDataByID(ctx context.Context, id int, owner *User, r PrivateDataRepository) (*BinaryData, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	model := model{
		id:             id,
		owner:          owner,
		modelType:      TypeUsepass,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
	}

	return &BinaryData{model: model}, nil
}
