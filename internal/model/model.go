package model

import (
	"context"
	"errors"
)

var _ PrivateData = (*model)(nil)

type model struct {
	id             int
	view           string
	modelType      Type
	owner          *User
	encryptedData  *EncryptedData
	dataRepository PrivateDataRepository
}

func (m *model) Save(ctx context.Context) (err error) {
	var saveFn func(context.Context, PrivateDataSave) (int, error)

	isExists := m.isExists()

	if isExists {
		saveFn = m.dataRepository.Update
	} else {
		saveFn = m.dataRepository.Add
	}

	dto := PrivateDataSave{
		ID:     m.id,
		Type:   m.modelType,
		UserID: m.owner.id,
		Data:   m.encryptedData.Data,
		DEK:    m.encryptedData.Key,
		View:   m.view,
	}
	if id, err := saveFn(ctx, dto); err == nil {
		m.id = id
	} else {
		return err
	}
	return nil
}

func (m *model) String() string {
	return m.view
}

func (m *model) Delete(ctx context.Context) (err error) {
	isExists := m.isExists()

	if !isExists {
		//TODO use package error
		return errors.New("entity is not exists")
	}
	return m.dataRepository.Delete(ctx, m.id)
}

// isExists - return false if user does not exists in storage
func (m model) isExists() bool {
	return m.id > 0
}

func (m model) ID() int {
	return m.id
}

func (m model) Owner() *User {
	return m.owner
}

func (m model) EncryptedData() EncryptedData {
	return *m.encryptedData
}

func (m model) Type() int {
	return int(m.modelType)
}
