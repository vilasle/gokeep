package model

import (
	"errors"
	"testing"
	"time"

	"context"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BankCard_Save(t *testing.T) {
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
		data           string
		dataName       string
		srcData        []byte
		owner          *User
		encryptedData  *EncryptedData
		successPD      bool
		successEncrypt bool
	}{
		{
			name:     "new text data",
			id:       0,
			data:     "some text",
			dataName: "name",
			srcData:  []byte("some text"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("some text"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successEncrypt: true,
		},
		{
			name:     "text data is existed, need to update",
			id:       1234,
			data:     "some text",
			dataName: "name",
			srcData:  []byte("some text"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("some text"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successEncrypt: true,
		},
		{
			name:     "new text data, need to add, got encryption error",
			data:     "some text",
			dataName: "name",
			srcData:  []byte("some text"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("some text"),
				Key:  []byte("key"),
			},
			successPD:      false,
			successEncrypt: false,
		},
		{
			name:     "new text data, need to add, got saving private data error",
			data:     "some text",
			dataName: "name",
			srcData:  []byte("some text"),
			owner:    &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("some text"),
				Key:  []byte("key"),
			},
			successEncrypt: true,
			successPD:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockPrivateDataRepository(ctrl)
			encrypter := NewMockEncoder(ctrl)

			entity := newPlainText(tt.owner, []byte{}, "")
			entity.dataRepository = r
			entity.id = tt.id

			entity.SetText([]byte(tt.data))
			entity.SetName(tt.dataName)

			entity.AddMetadata("key", "value")
			entity.AddMetadata("key2", "value2")

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
				Type:     entity.modelType,
				UserID:   entity.owner.ID(),
				Data:     tt.encryptedData.Data,
				DEK:      tt.encryptedData.Key,
				View:     entity.view,
				Metadata: entity.metadata,
			}
			if tt.successEncrypt {
				behaviorRepository(r, ctx, !entity.isExists(), dto, errPD)
			}

			err := entity.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successEncrypt, err == nil)
		})
	}
}

func Test_BankCard_decryptData(t *testing.T) {
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
		number         string
		cvv            int
		expiration     string
		srcData        []byte
		owner          *User
		encryptedData  *EncryptedData
		successPD      bool
		successEncrypt bool
	}{
		{
			name:       "new text data",
			id:         0,
			number:     "123456780",
			cvv:        123,
			expiration: "12/20",
			srcData:    []byte("123456780\n123\n12/20"),
			owner:      &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("123456780\n123\n12/20"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successEncrypt: true,
		},
		{
			name:       "text data is existed, need to update",
			id:         1234,
			number:     "123456780",
			cvv:        123,
			expiration: "12/20",
			srcData:    []byte("123456780\n123\n12/20"),
			owner:      &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("123456780\n123\n12/20"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successEncrypt: true,
		},
		{
			name:       "new text data, need to add, got encryption error",
			number:     "123456780",
			cvv:        123,
			expiration: "12/20",
			srcData:    []byte("123456780\n123\n12/20"),
			owner:      &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("123456780\n123\n12/20"),
				Key:  []byte("key"),
			},
			successPD:      false,
			successEncrypt: false,
		},
		{
			name:       "new text data, need to add, got saving private data error",
			number:     "123456780",
			cvv:        123,
			expiration: "12/20",
			srcData:    []byte("123456780\n123\n12/20"),
			owner:      &User{login: "test", password: "password"},
			encryptedData: &EncryptedData{
				Data: []byte("123456780\n123\n12/20"),
				Key:  []byte("key"),
			},
			successEncrypt: true,
			successPD:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockPrivateDataRepository(ctrl)
			encrypter := NewMockEncoder(ctrl)

			exp, err := time.Parse("01/06", tt.expiration)
			require.NoError(t, err)

			entity := newBankCard(tt.owner, "", 0, time.Now())
			entity.dataRepository = r
			entity.id = tt.id

			entity.SetNumber(tt.number)
			entity.SetCVV(tt.cvv)
			entity.SetExpiration(exp)

			entity.AddMetadata("key", "value")
			entity.AddMetadata("key2", "value2")

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
				Type:     entity.modelType,
				UserID:   entity.owner.ID(),
				Data:     tt.encryptedData.Data,
				DEK:      tt.encryptedData.Key,
				View:     entity.view,
				Metadata: entity.metadata,
			}
			if tt.successEncrypt {
				behaviorRepository(r, ctx, !entity.isExists(), dto, errPD)
			}

			err = entity.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successEncrypt, err == nil)
		})
	}

}

func Test_findBankCardByID(t *testing.T) {
	type repositoryMockArgs struct {
		input  int
		output PrivateDataInfo
		err    error
	}
	testCases := []struct {
		name    string
		input   int
		user    *User
		output  *PlainText
		wantErr bool
		repositoryMockArgs
	}{
		{
			name:  "success",
			input: 1,
			user: &User{
				id: 1,
			},
			output: &PlainText{
				model: model{
					id: 1,
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
					},
					metadata: map[string]string{
						"name": "test",
						"view": "test",
					},
				},
			},

			repositoryMockArgs: repositoryMockArgs{
				input: 1,
				output: PrivateDataInfo{
					ID:     1,
					UserID: 1,
					Type:   2,
					View:   "test",
					Data:   []byte("data"),
					DEK:    []byte("key"),
					Metadata: map[string]string{
						"name": "test",
						"view": "test",
					},
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
			output: &PlainText{
				model: model{
					id: 1,
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
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
					Data:     []byte("data"),
					DEK:      []byte("key"),
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

			actual, err := findBankCardByID(ctx, tt.input, tt.user, r)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			ed := actual.EncryptedData()

			assert.NoError(t, err)
			assert.Equal(t, tt.output.ID(), actual.id)
			assert.Equal(t, tt.user, actual.owner)
			assert.Equal(t, TypeBankCard, actual.Type())
			assert.Equal(t, ed.Data, actual.encryptedData.Data)
			assert.Equal(t, ed.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.output.Metadata(), actual.Metadata())

		})
	}
}
