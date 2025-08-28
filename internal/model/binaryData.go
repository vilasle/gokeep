package model

import (
	"context"
)

var _ PrivateData = (*BinaryData)(nil)

type BinaryData struct {
	model
	data []byte
}

func newBinaryData(owner *User, data []byte, name string) *BinaryData {
	return &BinaryData{
		model: model{
			view:      name,
			owner:     owner,
			modelType: TypeBinaryData,
		},
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

func (bd *BinaryData) SetName(name string) {
	bd.model.view = name
}

func (bd *BinaryData) SetData(data []byte) {
	bd.data = data
}

func (bd *BinaryData) prepareEncryptedData(encoder Encoder) (err error) {
	bd.encryptedData, err = encoder.Encrypt(bd.data)
	return err
}

func findBinaryDataByID(ctx context.Context, id int, owner *User, r PrivateDataRepository) (*BinaryData, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.id {
		return nil, ErrUserNotFound
	}

	model := model{
		id:             id,
		owner:          owner,
		modelType:      TypeBinaryData,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
		metadata: data.Metadata,
	}

	return &BinaryData{model: model}, nil
}
