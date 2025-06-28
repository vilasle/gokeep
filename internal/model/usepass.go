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
	Model
	login    string
	password string
}

func newUsepass(owner *User, login, password string) *Usepass {
	return &Usepass{
		Model: Model{
			owner: owner,
		},
		login:    login,
		password: password,
	}
}

func (u *Usepass) Save(ctx context.Context, encrypter Encrypter) (err error) {
	if err := u.prepareEncryptedData(encrypter); err != nil {
		//TODO wrap error with package error
		return err
	}
	return u.Model.Save(ctx)
}
func (u *Usepass) prepareEncryptedData(encrypter Encrypter) error {
	data := u.dataForEncryption()

	encryptedData, err := encrypter.Encrypt(data)
	if err != nil {
		return err
	}

	u.encryptedData = encryptedData

	u.encryptedData.Owner = u

	return nil
}

func (u *Usepass) decryptData(encrypter Encrypter) error {
	data, err := encrypter.Decrypt(*u.encryptedData)
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

func findUsepassByID(ctx context.Context, id int64, r PrivateDataRepository) (*Usepass, error) {
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
