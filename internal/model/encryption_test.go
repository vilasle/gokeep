package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncryptedData_ID(t *testing.T) {
	expected := int64(1)
	ed := &EncryptedData{
		id: 1,
	}
	assert.Equal(t, expected, ed.ID())
}

func Test_EncryptedData_Owner(t *testing.T) {
	expected := &Usepass{
		Model: Model{
			id: 1,
		},
	}

	ed := &EncryptedData{
		owner: expected,
	}
	assert.Equal(t, expected, ed.Owner())
}

func Test_EncryptedData_Data(t *testing.T) {
	expected := []byte("data")
	ed := &EncryptedData{
		data: expected,
	}
	assert.Equal(t, expected, ed.Data())

}

func Test_EncryptedData_Key(t *testing.T) {
	expected := []byte("key")
	ed := &EncryptedData{
		key: expected,
	}
	assert.Equal(t, expected, ed.Key())

}
