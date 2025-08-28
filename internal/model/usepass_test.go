package model

import (
	"errors"
	"testing"

	"context"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Usepass_Save(t *testing.T) {
	behaviorRepository := func(m *MockPrivateDataRepository, ctx context.Context, isExists bool, up PrivateDataSave, err error) {
		if isExists {
			m.EXPECT().Add(ctx, up).Return(1, err)
		} else {
			m.EXPECT().Update(ctx, up).Return(1, err)
		}
	}

	behaviorEncrypter := func(m *MockEncoder, src []byte, ed *EncryptedData, err error) {
		m.EXPECT().Encrypt(src).Return(ed, err)
	}

	testCases := []struct {
		name           string
		id             int
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
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
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
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
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
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
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
				Data: []byte("test1\npassword1"),
				Key:  []byte("key"),
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
			encrypter := NewMockEncoder(ctrl)

			usepass := newUsepass(tt.owner, "", "")
			usepass.SetUsername(tt.login)
			usepass.SetPassword(tt.password)
			usepass.dataRepository = usepassR

			usepass.id = tt.id

			ctx := context.Background()

			var errPD, errEnc error

			if !tt.successEncrypt {
				errEnc = errors.New("error")
			}

			behaviorEncrypter(encrypter, tt.srcData, tt.encryptedData, errEnc)

			if !tt.successPD {
				errPD = errors.New("error")
			}

			dto := PrivateDataSave{
				ID:       tt.id,
				Type:     usepass.modelType,
				UserID:   usepass.owner.id,
				Data:     tt.encryptedData.Data,
				DEK:      tt.encryptedData.Key,
				View:     usepass.view,
				Metadata: usepass.metadata,
			}
			if tt.successEncrypt {
				behaviorRepository(usepassR, ctx, !usepass.isExists(), dto, errPD)
			}

			err := usepass.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successED && tt.successEncrypt, err == nil)
		})
	}
}

func Test_findUsepassByID(t *testing.T) {
	type repositoryMockArgs struct {
		input  int
		output PrivateDataInfo
		err    error
	}
	testCases := []struct {
		name    string
		input   int
		user    *User
		output  *Usepass
		wantErr bool
		repositoryMockArgs
	}{
		{
			name:  "success",
			input: 1,
			user: &User{
				id: 1,
			},
			output: &Usepass{
				model: model{
					id: 1,
					encryptedData: &EncryptedData{
						Data: []byte("password"),
						Key:  []byte("password"),
					},
					metadata: make(map[string]string),
				},
			},

			repositoryMockArgs: repositoryMockArgs{
				input: 1,
				output: PrivateDataInfo{
					ID:       1,
					UserID:   1,
					Type:     2,
					View:     "test",
					Data:     []byte("password"),
					DEK:      []byte("password"),
					Metadata: make(map[string]string),
				},
			},
		},
		{
			name:  "repository error",
			input: 1,
			user: &User{
				id: 1,
			},
			wantErr: true,
			output:  nil,
			repositoryMockArgs: repositoryMockArgs{
				input:  1,
				output: PrivateDataInfo{},
				err:    errors.New("error"),
			},
		},
		{
			name:  "wrong user",
			input: 1,
			user: &User{
				id: 2,
			},
			wantErr: true,
			output: &Usepass{
				model: model{
					id: 1,
					encryptedData: &EncryptedData{
						Data: []byte("password"),
						Key:  []byte("password"),
					},
					metadata: make(map[string]string),
				},
			},

			repositoryMockArgs: repositoryMockArgs{
				input: 1,
				output: PrivateDataInfo{
					ID:       1,
					UserID:   1,
					Type:     2,
					View:     "test",
					Data:     []byte("password"),
					DEK:      []byte("password"),
					Metadata: make(map[string]string),
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockPrivateDataRepository(ctrl)

			ctx := context.Background()

			r.EXPECT().Get(ctx, tt.repositoryMockArgs.input).Return(tt.repositoryMockArgs.output, tt.repositoryMockArgs.err)

			actual, err := findUsepassByID(ctx, tt.input, tt.user, r)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.output.id, actual.id)
			assert.Equal(t, tt.user, actual.owner)
			assert.Equal(t, TypeUsepass, actual.modelType)
			assert.Equal(t, tt.output.encryptedData.Data, actual.encryptedData.Data)
			assert.Equal(t, tt.output.encryptedData.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.output.metadata, actual.metadata)

		})
	}
}
