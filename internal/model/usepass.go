package model

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
)

var _ PrivateData = (*Usepass)(nil)

// Usepass - work with logins and password and encrypt it
type Usepass struct {
	model
	login    string
	password string
}

func newUsepass(owner *User, login, password string) *Usepass {
	return &Usepass{
		model: model{
			owner: owner,
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

func (u *Usepass) String() (string) {
	return u.login
}

func (u *Usepass) SetUsername(username string) {
	u.login = username
}

func (u *Usepass) SetPassword(password string) {
	u.password = password
}

func (u *Usepass) prepareEncryptedData(encoder Encoder) error {
	data := u.dataForEncryption()

	encryptedData, err := encoder.Encrypt(data)
	if err != nil {
		return err
	}

	u.encryptedData = encryptedData

	u.encryptedData.Owner = u

	return nil
}

func (u *Usepass) decryptData(encoder Encoder) error {
	data, err := encoder.Decrypt(*u.encryptedData)
	if err != nil {
		return err
	}

	lp := strings.Split(string(data), "\n")
	if len(lp) != 2 {
		return errors.New("invalid data")
	}

	u.login, u.password = lp[0], lp[1]

	return nil
}

func (u Usepass) dataForEncryption() []byte {
	buf := bytes.Buffer{}
	buf.WriteString(u.login)
	buf.WriteString("\n")
	buf.WriteString(u.password)

	return buf.Bytes()
}

//FIXME add getting usepass by id and check that owner was right id
func findUsepassByID(ctx context.Context, id int64, owner *User, r PrivateDataRepository) (*Usepass, error) {
	data, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	u, ok := data.(*Usepass)
	if !ok {
		return nil, fmt.Errorf("invalid data. private data(id=%d) does not have type '*Usepass'", id)
	}

	return u, nil
}
