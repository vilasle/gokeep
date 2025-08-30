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
	owner          UserAccess
	encryptedData  *EncryptedData
	dataRepository PrivateDataRepository
	metadata       map[string]string
}

func (m *model) Save(ctx context.Context) (err error) {
	isExists := m.isExists()

	dto := PrivateDataSave{
		ID:       m.id,
		Type:     m.modelType,
		UserID:   m.owner.ID(),
		Data:     m.encryptedData.Data,
		DEK:      m.encryptedData.Key,
		View:     m.view,
		Metadata: m.metadata,
	}

	if isExists {
		m.id, err = m.dataRepository.Update(ctx, dto)
	} else {
		m.id, err = m.dataRepository.Add(ctx, dto)
	}
	return err
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

func (m *model) AddMetadata(key string, value string) {
	if m.metadata == nil {
		m.metadata = make(map[string]string)
	}
	m.metadata[key] = value
}

func (m *model) SetMetadata(metadata map[string]string) {
	m.metadata = metadata
}

// isExists - return false if user does not exists in storage
func (m model) isExists() bool {
	return m.id > 0
}

func (m model) ID() int {
	return m.id
}

func (m model) Owner() UserAccess {
	return m.owner
}

func (m model) EncryptedData() EncryptedData {
	return *m.encryptedData
}

func (m model) Type() int {
	return int(m.modelType)
}

func (m model) Metadata() map[string]string {
	r := make(map[string]string)
	for k, v := range m.metadata {
		r[k] = v
	}
	return r
}
