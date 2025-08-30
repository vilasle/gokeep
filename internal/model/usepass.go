package model

import (
	"bytes"
	"context"
	"errors"
)

var _ PrivateData = (*Usepass)(nil)

// Usepass - work with logins and password and encrypt it
type Usepass struct {
	model
	login    string
	password string
}

func newUsepass(owner UserAccess, login, password string) *Usepass {
	return &Usepass{
		model: model{
			view:      login,
			owner:     owner,
			modelType: TypeUsepass,
		},
		login:    login,
		password: password,
	}
}

func (u *Usepass) Save(ctx context.Context, encoder Encoder) (err error) {
	if err := u.prepareEncryptedData(encoder); err != nil {
		//TODO wrap error with package error
		return err
	}
	return u.model.Save(ctx)
}

func (u *Usepass) SetUsername(username string) {
	u.login = username
}

func (u *Usepass) SetPassword(password string) {
	u.password = password
}

func (u *Usepass) prepareEncryptedData(encoder Encoder) (err error) {
	data := u.dataForEncryption()

	u.encryptedData, err = encoder.Encrypt(data)
	return err
}

func (u Usepass) dataForEncryption() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(u.login)
	buf.WriteString("\n")
	buf.WriteString(u.password)

	return buf.Bytes()
}

func findUsepassByID(ctx context.Context, id int, owner UserAccess, r PrivateDataRepository) (*Usepass, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.UserID != owner.ID() {
		return nil, errors.New("user not found")
	}

	model := model{
		id:             id,
		view:           data.View,
		owner:          owner,
		modelType:      TypeUsepass,
		dataRepository: r,
		encryptedData: &EncryptedData{
			Data: data.Data,
			Key:  data.DEK,
		},
		metadata: data.Metadata,
	}

	return &Usepass{model: model}, nil
}
