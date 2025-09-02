package model

import (
	"context"
)

var _ PrivateData = (*BinaryData)(nil)

//BinaryData is model for presentation binary data
type BinaryData struct {
	model
	data []byte
}

func newBinaryData(owner UserAccess, data []byte, name string) *BinaryData {
	return &BinaryData{
		model: model{
			view:      name,
			owner:     owner,
			modelType: TypeBinaryData,
		},
		data: data,
	}
}

//Save prepare and encrypt binary data and save it in repository
func (bd *BinaryData) Save(ctx context.Context, encoder Encoder) (err error) {
	return bd.model.save(ctx, bd.data, encoder)
}

//SetName set new name for binary data
func (bd *BinaryData) SetName(name string) {
	bd.model.view = name
}

//SetData set new data for binary data
func (bd *BinaryData) SetData(data []byte) {
	bd.data = data
}

func findBinaryDataByID(ctx context.Context, id int, owner UserAccess, r PrivateDataRepository) (*BinaryData, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.ID() {
		return nil, ErrUserNotFound
	}

	model := model{
		id:             id,
		view:           data.View,
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
