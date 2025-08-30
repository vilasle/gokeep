package model

import (
	"context"
)

var _ PrivateData = (*PlainText)(nil)

type PlainText struct {
	model
	text []byte
}

func newPlainText(owner UserAccess, text []byte, name string) *PlainText {
	return &PlainText{
		model: model{
			view:      name,
			owner:     owner,
			modelType: TypePlainText,
		},
		text: text,
	}
}

func (tp *PlainText) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := tp.prepareEncryptedData(encoder); err != nil {
		//TODO wrap error with package error
		return err
	}
	return tp.model.Save(ctx)
}

func (tp *PlainText) prepareEncryptedData(encoder Encoder) (err error) {
	tp.encryptedData, err = encoder.Encrypt(tp.text)
	return err
}

func (tp *PlainText) SetText(text []byte) {
	tp.text = text
}

func (tp *PlainText) SetName(name string) {
	tp.model.view = name
}

// FIXME add getting plain text by id and check that owner was right id
func findPlainTextByID(ctx context.Context, id int, owner UserAccess, r PrivateDataRepository) (*PlainText, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.ID() {
		return nil, ErrNotFound
	}

	model := model{
		id:             id,
		view:           data.View,
		owner:          owner,
		modelType:      TypePlainText,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
		metadata: data.Metadata,
	}
	return &PlainText{model: model}, nil
}
