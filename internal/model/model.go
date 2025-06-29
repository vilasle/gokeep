package model

import (
	"context"
	"errors"
)

var _ PrivateData = (*model)(nil)

type model struct {
	id                   int64
	owner                *User
	encryptedData        *EncryptedData
	dataRepository       PrivateDataRepository
	encryptionRepository EncryptedDataRepository
}

func (m *model) Save(ctx context.Context) (err error) {
	var saveFn func(context.Context, PrivateData) error
	var saveEncryptedDataFn func(context.Context, *EncryptedData) error

	isExists := m.isExists()

	if isExists {
		saveFn = m.dataRepository.Update
		saveEncryptedDataFn = m.encryptionRepository.Update
	} else {
		saveFn = m.dataRepository.Add
		saveEncryptedDataFn = m.encryptionRepository.Add
	}

	if err := saveFn(ctx, m); err == nil {
		return saveEncryptedDataFn(ctx, m.encryptedData)
	} else {
		return err
	}
}

func (m *model) String() string {
	return ""
}

func (m *model) Delete(ctx context.Context) (err error) {
	isExists := m.isExists()

	if !isExists {
		//TODO use package error
		return errors.New("entity is not exists")
	}

	errs := make([]error, 0, 2)

	errs = append(errs, m.dataRepository.Delete(ctx, m))
	errs = append(errs, m.encryptionRepository.Delete(ctx, m.encryptedData))

	return errors.Join(errs...)
}

// isExists - return false if user does not exists in storage
func (m model) isExists() bool {
	return m.id > 0
}

func (m model) ID() int64 {
	return m.id
}

func (m model) Owner() *User {
	return m.owner
}

func (m model) EncryptedData() EncryptedData {
	return *m.encryptedData
}
