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
	id                   int64
	owner                *User
	login                string
	password             string
	encryptedData        *EncryptedData
	encrypter            Encrypter
	dataRepository       PrivateDataRepository
	encryptionRepository EncryptedDataRepository
}

func newUsepass(owner *User, login, password string,
	encrypter Encrypter,
	dR PrivateDataRepository,
	eR EncryptedDataRepository) *Usepass {

	return &Usepass{
		owner:                owner,
		login:                login,
		password:             password,
		encrypter:            encrypter,
		dataRepository:       dR,
		encryptionRepository: eR,
	}
}

func (u *Usepass) Save(ctx context.Context) (err error) {
	if err := u.prepareEncryptedData(); err != nil {
		//TODO wrap error with package error
		return err
	}

	var saveUsepassFn func(context.Context, PrivateData) error
	var saveEncryptedDataFn func(context.Context, *EncryptedData) error

	isExists := u.isExists()

	if isExists {
		saveUsepassFn = u.dataRepository.Update
		saveEncryptedDataFn = u.encryptionRepository.Update
	} else {
		saveUsepassFn = u.dataRepository.Add
		saveEncryptedDataFn = u.encryptionRepository.Add
	}

	if err := saveUsepassFn(ctx, u); err == nil {
		return saveEncryptedDataFn(ctx, u.encryptedData)
	} else {
		return err
	}
}

func (u *Usepass) Delete(ctx context.Context) (err error) {
	isExists := u.isExists()

	if !isExists {
		//TODO use package error
		return errors.New("entity is not exists")
	}

	errs := make([]error, 0, 2)

	errs = append(errs, u.dataRepository.Delete(ctx, u))
	errs = append(errs, u.encryptionRepository.Delete(ctx, u.encryptedData))

	return errors.Join(errs...)
}

func (u *Usepass) prepareEncryptedData() error {
	data := u.dataForEncryption()

	encryptedData, err := u.encrypter.Encrypt(data)
	if err != nil {
		return err
	}

	u.encryptedData = encryptedData

	u.encryptedData.owner = u

	return nil
}

// isExists - return false if user does not exists in storage
func (u Usepass) isExists() bool {
	return u.id > 0
}

func (u Usepass) ID() int64 {
	return u.id
}

func (u Usepass) Owner() *User {
	return u.owner
}

func (u Usepass) EncryptedData() EncryptedData {
	return *u.encryptedData
}

func (u *Usepass) decryptData() error {
	data, err := u.encrypter.Decrypt(*u.encryptedData)
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
