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

// String is a string representation of the model
func (m *model) String() string {
	return m.view
}

// Delete - delete entity from repository
func (m *model) Delete(ctx context.Context) (err error) {
	isExists := m.isExists()

	if !isExists {
		return errors.New("entity is not exists")
	}
	return m.dataRepository.Delete(ctx, m.id)
}

// AddMetadata - add metadata to model, if metadata already exists, it will be overwritten
func (m *model) AddMetadata(key string, value string) {
	if m.metadata == nil {
		m.metadata = make(map[string]string)
	}
	m.metadata[key] = value
}

// SetMetadata - set metadata to model
func (m *model) SetMetadata(metadata map[string]string) {
	m.metadata = metadata
}

// ID return id of model
func (m model) ID() int {
	return m.id
}

// Owner return owner of model
func (m model) Owner() UserAccess {
	return m.owner
}

// EncryptedData return encrypted data of model
func (m model) EncryptedData() EncryptedData {
	return *m.encryptedData
}

// Type return type of model
func (m model) Type() int {
	return int(m.modelType)
}

// Metadata return metadata of model
func (m model) Metadata() map[string]string {
	r := make(map[string]string)
	for k, v := range m.metadata {
		r[k] = v
	}
	return r
}

// Save prepare dto for writing data on repository
func (m *model) save(ctx context.Context, content []byte, encoder Encoder) (err error) {
	if m.encryptedData, err = encoder.Encrypt(content); err != nil {
		return err
	}

	dto := PrivateDataSave{
		ID:       m.id,
		Type:     m.modelType,
		UserID:   m.owner.ID(),
		Data:     m.encryptedData.Data,
		DEK:      m.encryptedData.Key,
		View:     m.view,
		Metadata: m.metadata,
	}

	if m.isExists() {
		m.id, err = m.dataRepository.Update(ctx, dto)
	} else {
		m.id, err = m.dataRepository.Add(ctx, dto)
	}
	return err
}

// isExists - return false if user does not exists in storage
func (m model) isExists() bool {
	return m.id > 0
}
