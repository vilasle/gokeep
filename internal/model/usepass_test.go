package model

import (
	"errors"
	"testing"

	"context"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Usepass_Save(t *testing.T) {
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, isExists bool, up PrivateData, err error) {
		if isExists {
			m.EXPECT().Add(ctx, up).Return(err)
		} else {
			m.EXPECT().Update(ctx, up).Return(err)
		}
	}

	behaviorED := func(m *MockEncryptedDataRepository, ctx context.Context, ed *EncryptedData, err error) {
		if ed.id == 0 {
			m.EXPECT().Add(ctx, ed).Return(err)
		} else {
			m.EXPECT().Update(ctx, ed).Return(err)
		}
	}

	behaviorEncrypter := func(m *MockEncrypter, src []byte, ed *EncryptedData, err error) {
		m.EXPECT().Encrypt(src).Return(ed, err)
	}

	testCases := []struct {
		name           string
		id             int64
		login          string
		password       string
		srcData        []byte
		owner          *User
		encryptedData  *EncryptedData
		successPD      bool
		successED      bool
		successEncrypt bool
	}{
		{
			name:     "new usepass, need to add",
			id:       0,
			login:    "test1",
			password: "password1",
			srcData:  []byte("test1\npassword1"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				id:   0,
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:     "usepass is existed, need to update",
			id:       1234,
			login:    "test1",
			password: "password1",
			srcData:  []byte("test1\npassword1"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				id:   1234,
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:     "new usepass, need to add, got encryption error",
			login:    "test1",
			password: "password1",
			srcData:  []byte("test1\npassword1"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
			},
			successPD:      false,
			successED:      false,
			successEncrypt: false,
		},
		{
			name:     "new usepass, need to add, got saving private data error",
			login:    "test1",
			password: "password1",
			srcData:  []byte("test1\npassword1"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
			},
			successEncrypt: true,
			successPD:      false,
			successED:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usepassR := NewMockPrivateDataRepository(ctrl)
			encryptedR := NewMockEncryptedDataRepository(ctrl)
			encrypter := NewMockEncrypter(ctrl)

			usepass := newUsepass(tt.owner, tt.login, tt.password)

			usepass.dataRepository = usepassR
			usepass.encryptionRepository = encryptedR

			usepass.id = tt.id

			ctx := context.Background()

			var errPD, errED, errEnc error

			if !tt.successEncrypt {
				errEnc = errors.New("error")
			}

			behaviorEncrypter(encrypter, tt.srcData, tt.encryptedData, errEnc)

			if !tt.successPD {
				errPD = errors.New("error")
			}

			if tt.successEncrypt {
				behaviorPD(usepassR, ctx, !usepass.isExists(), &usepass.Model, errPD)
			}

			//if can not save private data, does not try to save encrypted data
			if tt.successPD {
				if !tt.successED {
					errED = errors.New("error")
				}
				behaviorED(encryptedR, ctx, tt.encryptedData, errED)
			}

			err := usepass.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successED && tt.successEncrypt, err == nil)
		})
	}
}

func Test_Usepass_decryptData(t *testing.T) {
	behavior := func(m *MockEncrypter, ed EncryptedData, data []byte, err error) {
		m.EXPECT().Decrypt(ed).Return(data, err)
	}

	testCases := []struct {
		name             string
		ed               *EncryptedData
		data             []byte
		encErr           error
		err              error
		expectedLogin    string
		expectedPassword string
	}{
		{
			name: "success decryption",
			ed: &EncryptedData{
				key:  []byte("key"),
				data: []byte("test\npassword"),
			},
			data:             []byte("test\npassword"),
			expectedLogin:    "test",
			expectedPassword: "password",
		},
		{
			name: "encryption error",
			ed: &EncryptedData{
				key:  []byte("key"),
				data: []byte("test\npassword"),
			},
			data:             []byte{},
			encErr:           errors.New("error"),
			expectedLogin:    "",
			expectedPassword: "",
		},
		{
			name: "unexpected encryption data",
			ed: &EncryptedData{
				key:  []byte("key"),
				data: []byte("test\npassword"),
			},
			data:             []byte("test_password"),
			err:              errors.New("error"),
			expectedLogin:    "",
			expectedPassword: "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			enc := NewMockEncrypter(ctrl)

			behavior(enc, *tt.ed, tt.data, tt.encErr)

			usepass := newUsepass(nil, "", "")
			usepass.encryptedData = tt.ed

			err := usepass.decryptData(enc)

			if tt.err != nil || tt.encErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedLogin, usepass.login)
			assert.Equal(t, tt.expectedPassword, usepass.password)
		})
	}

}

func Test_findUsepassByID(t *testing.T) {
	testCases := []struct {
		name        string
		id          int64
		privateData PrivateData
		repErr      error
		err         error
		expected    *Usepass
	}{
		{
			name: "success",
			id:   1,
			privateData: &Usepass{
				Model: Model{
					id: 1,
				},
				login:    "test",
				password: "password",
			},
			err: nil,
			expected: &Usepass{
				Model: Model{
					id: 1,
				},
				login:    "test",
				password: "password",
			},
		},
		{
			name: "repository return wrong type",
			id:   1,
			privateData: &Usepass{
				Model: Model{
					id: 1,
				},
				login:    "test",
				password: "password",
			},
			repErr:   errors.New("repository error"),
			expected: nil,
		},
		{
			name:        "repository error",
			id:          1,
			privateData: &MockPrivateData{},
			err:         errors.New("wrong type"),
			expected: &Usepass{
				Model: Model{
					id: 1,
				},
				login:    "test",
				password: "password",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockPrivateDataRepository(ctrl)

			ctx := context.Background()

			r.EXPECT().Get(ctx, tt.id).Return(tt.privateData, tt.repErr)

			actual, err := findUsepassByID(ctx, tt.id, r)

			if tt.err != nil || tt.repErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
