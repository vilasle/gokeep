package model

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewModelManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockRepositoryCollector(ctrl)

	ur := NewMockUserRepository(ctrl)
	pvR := NewMockPrivateDataRepository(ctrl)
	encR := NewMockEncryptedDataRepository(ctrl)

	m.EXPECT().User().Return(ur)
	m.EXPECT().Private().Return(pvR).Times(3)
	m.EXPECT().Encryption().Return(encR).Times(3)

	NewModelManager(m)
}

func Test_userManager_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := NewMockUserRepository(ctrl)

	manager := userManager{
		repository: repository,
	}
	expected := &User{
		login:    "test",
		password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		r:        repository,
	}

	u := manager.New("test", "test")
	assert.Equal(t, expected, u)
}

func Test_userManager_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := NewMockUserRepository(ctrl)

	expected := &User{
		login:    "test",
		password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		r:        repository,
	}
	ctx := context.Background()
	repository.EXPECT().Get(ctx, "test").Return(expected, nil)

	manager := userManager{
		repository: repository,
	}

	u, _ := manager.FindByLogin(ctx, "test")
	assert.Equal(t, expected, u)
}

func Test_usepassManager_New(t *testing.T) {
	manager := usepassManager{
		pvRepository:  nil,
		encRepository: nil,
	}

	expected := &Usepass{
		Model:    Model{},
		login:    "test",
		password: "password",
	}

	usepass := manager.New(nil, "test", "password")

	assert.Equal(t, expected, usepass)
}

func Test_usepassManager_FindByID(t *testing.T) {
	testCases := []struct {
		name     string
		id       int64
		expected *Usepass
		err      error
	}{
		{
			name: "usepass is found",
			id:   1234,
			expected: &Usepass{
				Model: Model{
					id:    1234,
					owner: newUser("test", "test", nil),
				},
				login:    "test",
				password: "password",
			},
			err: nil,
		},
		{
			name:     "usepass is not found",
			id:       1234,
			expected: nil,
			err:      errors.New("usepass not found"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)
			encRepository := NewMockEncryptedDataRepository(ctrl)

			manager := usepassManager{
				pvRepository:  pvRepository,
				encRepository: encRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.id).Return(tt.expected, tt.err)

			usepass, err := manager.FindByID(ctx, tt.id)

			if tt.err != nil {
				assert.Error(t, tt.err, err)
				return
			}
			expected := tt.expected
			expected.dataRepository = pvRepository
			expected.encryptionRepository = encRepository

			assert.Equal(t, expected, usepass)
		})
	}
}

func Test_bankCardManager_New(t *testing.T) {
	manager := bankCardManager{
		pvRepository:  nil,
		encRepository: nil,
	}

	expected := &BankCard{
		Model: Model{
			owner: nil,
		},
		number:     "123456789012",
		cvv:        123,
		expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	usepass := manager.New(nil, "123456789012", 123, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))

	assert.Equal(t, expected, usepass)
}

func Test_bankCardManager_FindByID(t *testing.T) {
	testCases := []struct {
		name     string
		id       int64
		expected *BankCard
		err      error
	}{
		{
			name: "bank card is found",
			id:   1234,
			expected: &BankCard{
				Model: Model{
					id: 1234,
				},
				number:     "123456789012",
				cvv:        123,
				expiration: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			err: nil,
		},
		{
			name:     "bank card is not found",
			id:       1234,
			expected: nil,
			err:      errors.New("usepass not found"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)
			encRepository := NewMockEncryptedDataRepository(ctrl)

			manager := bankCardManager{
				pvRepository:  pvRepository,
				encRepository: encRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.id).Return(tt.expected, tt.err)

			bankCard, err := manager.FindByID(ctx, tt.id)

			if tt.err != nil {
				assert.Error(t, tt.err, err)
				return
			}
			expected := tt.expected
			expected.dataRepository = pvRepository
			expected.encryptionRepository = encRepository

			assert.Equal(t, expected, bankCard)
		})
	}
}
