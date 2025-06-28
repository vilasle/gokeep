package model

import (
	"errors"
	"testing"

	"context"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PlainText_Save(t *testing.T) {
	behaviorPD := func(m *MockPrivateDataRepository, ctx context.Context, isExists bool, up PrivateData, err error) {
		if isExists {
			m.EXPECT().Add(ctx, up).Return(err)
		} else {
			m.EXPECT().Update(ctx, up).Return(err)
		}
	}

	behaviorED := func(m *MockEncryptedDataRepository, ctx context.Context, ed *EncryptedData, err error) {
		if ed.ID == 0 {
			m.EXPECT().Add(ctx, ed).Return(err)
		} else {
			m.EXPECT().Update(ctx, ed).Return(err)
		}
	}

	behaviorEncrypter := func(m *MockEncrypter, src []byte, ed *EncryptedData, err error) {
		m.EXPECT().Encrypt(src).Return(ed, err)
	}

	testCases := []struct {
		name string
		id   int64
		// src text
		data []byte
		// compressed text
		srcData        []byte
		encryptedData  *EncryptedData
		successPD      bool
		successED      bool
		successEncrypt bool
	}{
		{
			name:    "new plain text, need to add",
			id:      0,
			data:    []byte("some text"),
			srcData: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0},
			encryptedData: &EncryptedData{
				Data: []byte{},
				Key:  []byte{},
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:    "plain text is existed, need to update",
			id:      1234,
			data:    []byte("some text"),
			srcData: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0},
			encryptedData: &EncryptedData{
				ID:   1234,
				Data: []byte{},
				Key:  []byte{},
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:    "new plain text, need to add, got encryption error",
			data:    []byte("some text"),
			srcData: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0},
			encryptedData: &EncryptedData{
				Data: []byte{},
				Key:  []byte{},
			},
			successPD:      false,
			successED:      false,
			successEncrypt: false,
		},
		{
			name:    "new plain text, need to add, got saving private data error",
			data:    []byte("some text"),
			srcData: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0},
			encryptedData: &EncryptedData{
				Data: []byte{},
				Key:  []byte{},
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

			plainText := newPlainText(nil, tt.data)

			plainText.dataRepository = usepassR
			plainText.encryptionRepository = encryptedR

			plainText.id = tt.id

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
				behaviorPD(usepassR, ctx, !plainText.isExists(), &plainText.Model, errPD)
			}

			//if can not save private data, does not try to save encrypted data
			if tt.successPD {
				if !tt.successED {
					errED = errors.New("error")
				}
				behaviorED(encryptedR, ctx, tt.encryptedData, errED)
			}

			err := plainText.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successED && tt.successEncrypt, err == nil)
		})
	}
}

func Test_PlainText_decryptData(t *testing.T) {
	behavior := func(m *MockEncrypter, ed EncryptedData, data []byte, err error) {
		m.EXPECT().Decrypt(ed).Return(data, err)
	}

	testCases := []struct {
		name        string
		ed          *EncryptedData
		data        []byte
		encErr      error
		err         error
		expectedTxt []byte
	}{
		{
			name: "success decryption",
			ed: &EncryptedData{
				Key:  []byte{},
				Data: []byte{},
			},
			data: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0}, //some text
			expectedTxt: []byte("some text"),
		},
		{
			name: "encryption error",
			ed: &EncryptedData{
				Key:  []byte{},
				Data: []byte{},
			},
			data: []byte{31,139,8,0,0,0,0,0,2,255,
				42,206,207,77,85,40,73,173,40,1,4,0,0,255,255,186,189,186,79,9,0,0,0}, //some text
			encErr: errors.New("error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			enc := NewMockEncrypter(ctrl)

			behavior(enc, *tt.ed, tt.data, tt.encErr)

			plainText := newPlainText(nil, []byte{})
			plainText.encryptedData = tt.ed

			err := plainText.decryptData(enc)

			if tt.err != nil || tt.encErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedTxt, plainText.text)
		})
	}

}

func Test_findPlainTextByID(t *testing.T) {
	testCases := []struct {
		name        string
		id          int64
		privateData PrivateData
		repErr      error
		err         error
		expected    *PlainText
	}{
		{
			name: "success",
			id:   1,
			privateData: &PlainText{
				Model: Model{
					id: 1,
				},
				text: []byte("some text"),
			},
			err: nil,
			expected: &PlainText{
				Model: Model{
					id: 1,
				},
				text: []byte("some text"),
			},
		},
		{
			name:        "repository return wrong type",
			id:          1,
			privateData: nil,
			repErr:      errors.New("repository error"),
			expected:    nil,
		},
		{
			name:        "repository error",
			id:          1,
			privateData: &MockPrivateData{},
			err:         errors.New("wrong type"),
			expected:    nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockPrivateDataRepository(ctrl)

			ctx := context.Background()

			r.EXPECT().Get(ctx, tt.id).Return(tt.privateData, tt.repErr)

			actual, err := findPlainTextByID(ctx, tt.id, r)

			if tt.err != nil || tt.repErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
