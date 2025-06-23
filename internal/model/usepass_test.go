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
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, up *Usepass, err error) {
		if up.id == 0 {
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
			name:     "new user, need to add, got encryption error",
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
			name:     "new user, need to add, got saving private data error",
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

			usepass := newUsepass(
				tt.owner,
				tt.login,
				tt.password,
				encrypter,
				usepassR,
				encryptedR)

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
				behaviorPD(usepassR, ctx, usepass, errPD)
			}

			//if can not save private data, does not try to save encrypted data
			if tt.successPD {
				if !tt.successED {
					errED = errors.New("error")
				}
				behaviorED(encryptedR, ctx, tt.encryptedData, errED)
			}

			err := usepass.Save(ctx)

			assert.Equal(t, tt.successPD && tt.successED && tt.successEncrypt, err == nil)
		})
	}

}

func Test_Usepass_Delete(t *testing.T) {
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, up *Usepass, err error) {
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
				id:   0,
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
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
				id:   1234,
				data: []byte("test1\npassword1"),
				key:  []byte("key"),
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

			usepassR := NewMockPrivateDataRepository(ctrl)
			encryptedR := NewMockEncryptedDataRepository(ctrl)
			encrypter := NewMockEncrypter(ctrl)

			usepass := newUsepass(
				tt.owner,
				tt.login,
				tt.password,
				encrypter,
				usepassR,
				encryptedR)

			usepass.id = tt.id
			usepass.encryptedData = tt.encryptedData

			ctx := context.Background()

			var errPD, errED error

			if tt.isExists {
				if !tt.successPD {
					errPD = errors.New("error")
				}

				behaviorPD(usepassR, ctx, usepass, errPD)

				//if can not save private data, does not try to save encrypted data
				if tt.successPD {
					if !tt.successED {
						errED = errors.New("error")
					}
					behaviorED(encryptedR, ctx, tt.encryptedData, errED)
				}
			}

			err := usepass.Delete(ctx)

			assert.Equal(t, tt.successPD && tt.successED && tt.isExists, err == nil)
		})
	}
}

func Test_Usepass_ID(t *testing.T) {
	expected := int64(1)
	usepass := newUsepass(nil, "", "", nil, nil, nil)
	usepass.id = expected

	assert.Equal(t, expected, usepass.ID())
}

func Test_Usepass_Owner(t *testing.T) {
	expected := newUser("test", "password", nil)
	usepass := newUsepass(expected, "", "", nil, nil, nil)

	assert.Equal(t, expected, usepass.Owner())
}

func Test_Usepass_EncryptedData(t *testing.T) {
	usepass := newUsepass(nil, "", "", nil, nil, nil)
	expected := EncryptedData{
		key:  []byte("key"),
		data: []byte("data"),
	}
	usepass.encryptedData = &expected
	assert.Equal(t, expected, usepass.EncryptedData())
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

			behavior(enc, *tt.ed, tt.data, tt.err)

			usepass := newUsepass(nil, "", "", enc, nil, nil)
			usepass.encryptedData = tt.ed

			err := usepass.decryptData()

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
				id:       1,
				login:    "test",
				password: "password",
			},
			err: nil,
			expected: &Usepass{
				id:       1,
				login:    "test",
				password: "password",
			},
		},
		{
			name: "repository return wrong type",
			id:   1,
			privateData: &Usepass{
				id:       1,
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
				id:       1,
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
