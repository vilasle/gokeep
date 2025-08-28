package model

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Usepass_Delete(t *testing.T) {
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, id int, err error) {
		m.EXPECT().Delete(ctx, id).Return(err)
	}

	testCases := []struct {
		name          string
		id            int
		login         string
		password      string
		srcData       []byte
		owner         *User
		encryptedData *EncryptedData
		isExists      bool
		successPD     bool
	}{
		{
			name:     "new usepass, got error",
			id:       0,
			login:    "test1",
			password: "password1",
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
			},
			isExists:  false,
			successPD: false,
		},
		{
			name:     "new usepass, got error",
			id:       1234,
			login:    "test1",
			password: "password1",
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
			},
			isExists:  true,
			successPD: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvR := NewMockPrivateDataRepository(ctrl)

			model := &model{
				id:             tt.id,
				owner:          tt.owner,
				dataRepository: pvR,
				encryptedData:  tt.encryptedData,
			}
			ctx := context.Background()

			var errPD error

			if tt.isExists {
				if !tt.successPD {
					errPD = errors.New("error")
				}

				behaviorPD(pvR, ctx, tt.id, errPD)
			}

			err := model.Delete(ctx)

			assert.Equal(t, tt.successPD && tt.isExists, err == nil)
		})
	}
}
