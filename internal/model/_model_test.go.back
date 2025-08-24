package model

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Usepass_Delete(t *testing.T) {
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, up PrivateData, err error) {
		m.EXPECT().Delete(ctx, up).Return(err)
	}

	behaviorED := func(m *MockEncryptedDataRepository, ctx context.Context, ed *EncryptedData, err error) {
		m.EXPECT().Delete(ctx, ed).Return(err)
	}

	testCases := []struct {
		name          string
		id            int64
		login         string
		password      string
		srcData       []byte
		owner         *User
		encryptedData *EncryptedData
		isExists      bool
		successPD     bool
		successED     bool
	}{
		{
			name:     "new usepass, got error",
			id:       0,
			login:    "test1",
			password: "password1",
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				ID:   0,
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
			},
			isExists:  false,
			successPD: false,
			successED: false,
		},
		{
			name:     "new usepass, got error",
			id:       1234,
			login:    "test1",
			password: "password1",
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				ID:   1234,
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
			},
			isExists:  true,
			successPD: true,
			successED: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvR := NewMockPrivateDataRepository(ctrl)
			encR := NewMockEncryptedDataRepository(ctrl)

			model := &model{
				id:                   tt.id,
				owner:                tt.owner,
				dataRepository:       pvR,
				encryptionRepository: encR,
				encryptedData:        tt.encryptedData,
			}
			ctx := context.Background()

			var errPD, errED error

			if tt.isExists {
				if !tt.successPD {
					errPD = errors.New("error")
				}

				behaviorPD(pvR, ctx, model, errPD)

				//if can not save private data, does not try to save encrypted data
				if tt.successPD {
					if !tt.successED {
						errED = errors.New("error")
					}
					behaviorED(encR, ctx, tt.encryptedData, errED)
				}
			}

			err := model.Delete(ctx)

			assert.Equal(t, tt.successPD && tt.successED && tt.isExists, err == nil)
		})
	}
}

func Test_Model_ID(t *testing.T) {
	expected := int64(1)
	usepass := model{
		id: expected,
	}
	usepass.id = expected

	assert.Equal(t, expected, usepass.ID())
}

func Test_Model_Owner(t *testing.T) {
	expected := newUser("test", "password", nil)
	usepass := model{
		owner: expected,
	}
	assert.Equal(t, expected, usepass.Owner())
}

func Test_Model_EncryptedData(t *testing.T) {
	usepass := model{}
	expected := EncryptedData{
		Key:  []byte("key"),
		Data: []byte("data"),
	}
	usepass.encryptedData = &expected
	assert.Equal(t, expected, usepass.EncryptedData())
}
