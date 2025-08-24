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
		name           string
		id             int64
		numberCard     string
		cvv            int
		expiration     time.Time
		srcData        []byte
		encryptedData  *EncryptedData
		successPD      bool
		successED      bool
		successEncrypt bool
	}{
		{
			name:       "new bank card, need to add",
			id:         0,
			numberCard: "1234567890123",
			cvv:        123,
			expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			srcData:    []byte("1234567890123\n123\n2020-01-01"),
			encryptedData: &EncryptedData{
				Data: []byte("1234567890123\n123\n2020-01-01"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:       "bank card is existed, need to update",
			id:         1234,
			numberCard: "1234567890123",
			cvv:        123,
			expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			srcData:    []byte("1234567890123\n123\n2020-01-01"),
			encryptedData: &EncryptedData{
				ID:   1234,
				Data: []byte("1234567890123\n123\n2020-01-01"),
				Key:  []byte("key"),
			},
			successPD:      true,
			successED:      true,
			successEncrypt: true,
		},
		{
			name:       "new user, need to add, got encryption error",
			numberCard: "1234567890123",
			cvv:        123,
			expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			srcData:    []byte("1234567890123\n123\n2020-01-01"),
			encryptedData: &EncryptedData{
				Data: []byte("1234567890123\n123\n2020-01-01"),
				Key:  []byte("key"),
			},
			successPD:      false,
			successED:      false,
			successEncrypt: false,
		},
		{
			name:       "new user, need to add, got saving private data error",
			numberCard: "1234567890123",
			cvv:        123,
			expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			srcData:    []byte("1234567890123\n123\n2020-01-01"),
			encryptedData: &EncryptedData{
				Data: []byte("1234567890123\n123\n2020-01-01"),
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
			encryptedR := NewMockEncryptedDataRepository(ctrl)
			encrypter := NewMockEncrypter(ctrl)

			bankCard := newBankCard(nil, tt.numberCard, tt.cvv, tt.expiration)

			bankCard.dataRepository = usepassR
			bankCard.encryptionRepository = encryptedR

			bankCard.id = tt.id

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
				behaviorPD(usepassR, ctx, !bankCard.isExists(), &bankCard.model, errPD)
			}

			//if can not save private data, does not try to save encrypted data
			if tt.successPD {
				if !tt.successED {
					errED = errors.New("error")
				}
				behaviorED(encryptedR, ctx, tt.encryptedData, errED)
			}

			err := bankCard.Save(ctx, encrypter)

			assert.Equal(t, tt.successPD && tt.successED && tt.successEncrypt, err == nil)
		})
	}
}

func Test_BankCard_decryptData(t *testing.T) {
	behavior := func(m *MockEncrypter, ed EncryptedData, data []byte, err error) {
		m.EXPECT().Decrypt(ed).Return(data, err)
	}

	testCases := []struct {
		name               string
		ed                 *EncryptedData
		data               []byte
		encErr             error
		err                error
		expectedNumber     string
		expectedCVV        int
		expectedExpiration time.Time
	}{
		{
			name: "success decryption",
			ed: &EncryptedData{
				Key:  []byte("key"),
				Data: []byte("test\npassword"),
			},
			data:               []byte("1234567890123\n123\n2020-01-01"),
			expectedNumber:     "1234567890123",
			expectedCVV:        123,
			expectedExpiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "encryption error",
			ed: &EncryptedData{
				Key:  []byte("key"),
				Data: []byte("1234567890123\n123\n2020-01-01"),
			},
			data:               []byte{},
			encErr:             errors.New("error"),
			expectedNumber:     "",
			expectedCVV:        0,
			expectedExpiration: time.Time{},
		},
		{
			name: "unexpected encryption data, wrong format",
			ed: &EncryptedData{
				Key:  []byte("key"),
				Data: []byte("1234567890123\n123\n2020-01-01"),
			},
			data:               []byte("1234567890123_123_2020-01-01"),
			err:                errors.New("error"),
			expectedNumber:     "",
			expectedCVV:        0,
			expectedExpiration: time.Time{},
		},
		{
			name: "unexpected encryption data, cvv is not a number",
			ed: &EncryptedData{
				Key:  []byte("key"),
				Data: []byte("1234567890123\n123\n2020-01-01"),
			},
			data:               []byte("1234567890123\ndsfdsaf\n2020-01-01"),
			err:                errors.New("error"),
			expectedNumber:     "",
			expectedCVV:        0,
			expectedExpiration: time.Time{},
		},
		{
			name: "unexpected encryption data, expiration is not a time",
			ed: &EncryptedData{
				Key:  []byte("key"),
				Data: []byte("1234567890123\n123\n2020-01-01"),
			},
			data:               []byte("1234567890123\n123\nv2213-3v-vfd1"),
			err:                errors.New("error"),
			expectedNumber:     "",
			expectedCVV:        0,
			expectedExpiration: time.Time{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			enc := NewMockEncrypter(ctrl)

			behavior(enc, *tt.ed, tt.data, tt.encErr)

			bankCard := newBankCard(nil, "", 0, time.Time{})
			bankCard.encryptedData = tt.ed

			err := bankCard.decryptData(enc)

			if tt.err != nil || tt.encErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedNumber, bankCard.number)
			assert.Equal(t, tt.expectedCVV, bankCard.cvv)
			assert.Equal(t, tt.expectedExpiration, bankCard.expiration)
		})
	}

}

func Test_findBankCardByID(t *testing.T) {
	testCases := []struct {
		name        string
		id          int64
		privateData PrivateData
		repErr      error
		err         error
		expected    *BankCard
	}{
		{
			name: "success",
			id:   1,
			privateData: &BankCard{
				model: model{
					id: 1,
				},
				number:     "123456789012",
				cvv:        321,
				expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			err: nil,
			expected: &BankCard{
				model: model{
					id: 1,
				},
				number:     "123456789012",
				cvv:        321,
				expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
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

			actual, err := findBankCardByID(ctx, tt.id, nil, r)

			if tt.err != nil || tt.repErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
