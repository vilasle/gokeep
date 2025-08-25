package model

import (
	"context"
)

var _ PrivateData = (*PlainText)(nil)

type PlainText struct {
	model
	name string
	text []byte
}

func newPlainText(owner *User, text []byte, name string) *PlainText {
	return &PlainText{
		model: model{
			view:      name,
			owner:     owner,
			modelType: TypePlainText,
		},
		name: name,
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

// FIXME add getting plain text by id and check that owner was right id
func findPlainTextByID(ctx context.Context, id int, owner *User, r PrivateDataRepository) (*PlainText, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.id {
		return nil, ErrNotFound
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
	return &PlainText{model: model}, nil
}
