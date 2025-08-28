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

	m.EXPECT().User().Return(ur)
	m.EXPECT().Private().Return(pvR).Times(4)

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

	//success case
	expectedResponse := UserInfo{
		ID:       1,
		Login:    "test",
		Password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	user := &User{
		id:       1,
		login:    "test",
		password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}

	ctx := context.Background()
	repository.EXPECT().Find(ctx, "test").Return(expectedResponse, nil)

	manager := userManager{
		repository: repository,
	}

	u, _ := manager.FindByLogin(ctx, "test")
	assert.Equal(t, user.id, u.id)
	assert.Equal(t, user.login, u.login)
	assert.Equal(t, user.password, u.password)

	//fail case
	repository.EXPECT().Find(ctx, "test").Return(UserInfo{}, errors.New("user not found"))
	_, err := manager.FindByLogin(ctx, "test")
	assert.Error(t, err)

}

func Test_userManager_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := NewMockUserRepository(ctrl)

	//success case
	expectedResponse := UserInfo{
		ID:       1,
		Login:    "test",
		Password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}
	user := &User{
		id:       1,
		login:    "test",
		password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	}

	ctx := context.Background()
	repository.EXPECT().Get(ctx, 1).Return(expectedResponse, nil)

	manager := userManager{
		repository: repository,
	}

	u, _ := manager.Get(ctx, 1)
	assert.Equal(t, user.id, u.id)
	assert.Equal(t, user.login, u.login)
	assert.Equal(t, user.password, u.password)

	//fail case
	repository.EXPECT().Get(ctx, 1).Return(UserInfo{}, errors.New("user not found"))
	_, err := manager.Get(ctx, 1)
	assert.Error(t, err)

}

func Test_usepassManager_New(t *testing.T) {
	manager := usepassManager{
		pvRepository: nil,
	}

	expected := &Usepass{
		model:    model{},
		login:    "test",
		password: "password",
	}

	usepass := manager.New(nil, "test", "password")

	assert.Equal(t, expected.login, usepass.login)
	assert.Equal(t, expected.password, usepass.password)
}

func Test_usepassManager_Get(t *testing.T) {
	type argsMock struct {
		id  int
		pd  PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name     string
		id       int
		user     *User
		usepass  *Usepass
		mockArgs argsMock
		wantErr  bool
	}{
		{
			name: "usepass is found",
			id:   1234,
			user: user,
			usepass: &Usepass{
				model: model{
					id:   1234,
					view: "test",
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
					},
					metadata: map[string]string{},
				},
			},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeUsepass,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: nil,
			},

			wantErr: false,
		},
		{
			name:    "usepass is not found",
			id:      1234,
			user:    user,
			usepass: &Usepass{},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeUsepass,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: errors.New("usepass not found"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := usepassManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.mockArgs.id).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.Get(ctx, tt.user, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.usepass.id, actual.id)
			assert.Equal(t, tt.usepass.view, actual.view)
			assert.Equal(t, tt.usepass.encryptedData.Data, actual.encryptedData.Data)
			assert.Equal(t, tt.usepass.encryptedData.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.usepass.metadata, actual.metadata)

		})
	}
}

func Test_usepassManager_List(t *testing.T) {
	type argsMock struct {
		id  int
		pd  []PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name      string
		id        int
		user      *User
		modelType int
		expected  []*Usepass
		mockArgs  argsMock
		wantErr   bool
	}{
		{
			name:      "success getting list of usepass",
			id:        1234,
			user:      user,
			modelType: TypeUsepass,
			mockArgs: argsMock{
				id: 1234,
				pd: []PrivateDataInfo{
					{
						ID:       1234,
						UserID:   4321,
						Type:     TypeUsepass,
						View:     "test",
						Data:     []byte("data"),
						DEK:      []byte("key"),
						Metadata: map[string]string{},
					},
				},
				err: nil,
			},
			expected: []*Usepass{
				{
					model: model{
						id:    1234,
						view:  "test",
						owner: user,
						encryptedData: &EncryptedData{
							Data: []byte("data"),
							Key:  []byte("key"),
						},
						modelType: TypeUsepass,
						metadata:  map[string]string{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := usepassManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().List(ctx, tt.modelType, tt.user).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.List(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, actual, len(tt.expected))
		})
	}
}

func Test_bankCardManager_New(t *testing.T) {
	manager := bankCardManager{
		pvRepository: nil,
	}

	expected := &BankCard{
		model:      model{},
		number:     "123456780",
		cvv:        123,
		expiration: time.Now(),
	}

	card := manager.New(nil, expected.number, expected.cvv, expected.expiration)

	assert.Equal(t, expected.number, card.number)
	assert.Equal(t, expected.cvv, card.cvv)
	assert.Equal(t, expected.expiration, card.expiration)

}

func Test_bankCardManager_Get(t *testing.T) {
	type argsMock struct {
		id  int
		pd  PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name     string
		id       int
		user     *User
		usepass  *Usepass
		mockArgs argsMock
		wantErr  bool
	}{
		{
			name: "bank card is found",
			id:   1234,
			user: user,
			usepass: &Usepass{
				model: model{
					id:   1234,
					view: "test",
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
					},
					metadata: map[string]string{},
				},
			},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeBankCard,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: nil,
			},

			wantErr: false,
		},
		{
			name:    "bank card is not found",
			id:      1234,
			user:    user,
			usepass: &Usepass{},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeBankCard,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: errors.New("usepass not found"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := bankCardManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.mockArgs.id).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.Get(ctx, tt.user, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.usepass.id, actual.id)
			assert.Equal(t, tt.usepass.view, actual.view)
			assert.Equal(t, tt.usepass.encryptedData.Data, actual.encryptedData.Data)
			assert.Equal(t, tt.usepass.encryptedData.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.usepass.metadata, actual.metadata)

		})
	}
}

func Test_bankCardManager_List(t *testing.T) {
	type argsMock struct {
		id  int
		pd  []PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name      string
		id        int
		user      *User
		modelType int
		expected  []*BankCard
		mockArgs  argsMock
		wantErr   bool
	}{
		{
			name:      "success getting list of bank cards",
			id:        1234,
			user:      user,
			modelType: TypeBankCard,
			mockArgs: argsMock{
				id: 1234,
				pd: []PrivateDataInfo{
					{
						ID:       1234,
						UserID:   4321,
						Type:     TypeBankCard,
						View:     "test",
						Data:     []byte("data"),
						DEK:      []byte("key"),
						Metadata: map[string]string{},
					},
				},
				err: nil,
			},
			expected: []*BankCard{
				{
					model: model{
						id:    1234,
						view:  "test",
						owner: user,
						encryptedData: &EncryptedData{
							Data: []byte("data"),
							Key:  []byte("key"),
						},
						modelType: TypeBankCard,
						metadata:  map[string]string{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := bankCardManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().List(ctx, tt.modelType, tt.user).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.List(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, actual, len(tt.expected))
		})
	}
}

func Test_textDataManager_New(t *testing.T) {
	manager := plainTextManager{
		pvRepository: nil,
	}

	expected := &PlainText{
		model: model{},
		text:  []byte("123456780"),
	}

	text := manager.New(nil, expected.text, "test")

	assert.Equal(t, expected.text, text.text)
	assert.Equal(t, "test", text.view)

}

func Test_textDataManager_Get(t *testing.T) {
	type argsMock struct {
		id  int
		pd  PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name     string
		id       int
		user     *User
		text     *PlainText
		mockArgs argsMock
		wantErr  bool
	}{
		{
			name: "binary data is found",
			id:   1234,
			user: user,
			text: &PlainText{
				model: model{
					id:   1234,
					view: "test",
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
					},
					metadata: map[string]string{},
				},
			},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypePlainText,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: nil,
			},

			wantErr: false,
		},
		{
			name: "binary data is not found",
			id:   1234,
			user: user,
			text: &PlainText{},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypePlainText,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: errors.New("binary data not found"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := plainTextManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.mockArgs.id).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.Get(ctx, tt.user, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.text.id, actual.id)
			assert.Equal(t, tt.text.view, actual.view)
			assert.Equal(t, tt.text.encryptedData.Data, actual.encryptedData.Data)
			assert.Equal(t, tt.text.encryptedData.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.text.metadata, actual.metadata)

		})
	}
}

func Test_textDataManager_List(t *testing.T) {
	type argsMock struct {
		id  int
		pd  []PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name      string
		id        int
		user      *User
		modelType int
		expected  []*PlainText
		mockArgs  argsMock
		wantErr   bool
	}{
		{
			name:      "success getting list of text data",
			id:        1234,
			user:      user,
			modelType: TypePlainText,
			mockArgs: argsMock{
				id: 1234,
				pd: []PrivateDataInfo{
					{
						ID:       1234,
						UserID:   4321,
						Type:     TypePlainText,
						View:     "test",
						Data:     []byte("data"),
						DEK:      []byte("key"),
						Metadata: map[string]string{},
					},
				},
				err: nil,
			},
			expected: []*PlainText{
				{
					model: model{
						id:    1234,
						view:  "test",
						owner: user,
						encryptedData: &EncryptedData{
							Data: []byte("data"),
							Key:  []byte("key"),
						},
						modelType: TypePlainText,
						metadata:  map[string]string{},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := plainTextManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().List(ctx, tt.modelType, tt.user).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.List(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, actual, len(tt.expected))
		})
	}
}

func Test_binaryDataManager_New(t *testing.T) {
	manager := binaryDataManager{
		pvRepository: nil,
	}

	expected := &BinaryData{
		model: model{},
		data:  []byte("data"),
	}

	binary := manager.New(nil, expected.data, "test")

	assert.Equal(t, expected.data, binary.data)
	assert.Equal(t, "test", binary.view)

}

func Test_binaryDataManager_Get(t *testing.T) {
	type argsMock struct {
		id  int
		pd  PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name     string
		id       int
		user     *User
		binary   *BinaryData
		mockArgs argsMock
		wantErr  bool
	}{
		{
			name: "binary data is found",
			id:   1234,
			user: user,
			binary: &BinaryData{
				model: model{
					id:   1234,
					view: "test",
					encryptedData: &EncryptedData{
						Data: []byte("data"),
						Key:  []byte("key"),
					},
					metadata: map[string]string{},
				},
			},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeBinaryData,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: nil,
			},

			wantErr: false,
		},
		{
			name:   "binary data is not found",
			id:     1234,
			user:   user,
			binary: &BinaryData{},
			mockArgs: argsMock{
				id: 1234,
				pd: PrivateDataInfo{
					ID:       1234,
					UserID:   4321,
					Type:     TypeBinaryData,
					View:     "test",
					Data:     []byte("data"),
					DEK:      []byte("key"),
					Metadata: map[string]string{},
				},
				err: errors.New("binary data not found"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := binaryDataManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().Get(ctx, tt.mockArgs.id).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.Get(ctx, tt.user, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.binary.id, actual.id)
			assert.Equal(t, tt.binary.view, actual.view)
			assert.Equal(t, tt.binary.encryptedData.Data, actual.encryptedData.Data)
			assert.Equal(t, tt.binary.encryptedData.Key, actual.encryptedData.Key)
			assert.Equal(t, tt.binary.metadata, actual.metadata)

		})
	}
}

func Test_binaryDataManager_List(t *testing.T) {
	type argsMock struct {
		id  int
		pd  []PrivateDataInfo
		err error
	}

	user := &User{
		id:       4321,
		login:    "test",
		password: "password",
	}

	testCases := []struct {
		name      string
		id        int
		user      *User
		modelType int
		expected  []*BinaryData
		mockArgs  argsMock
		wantErr   bool
	}{
		{
			name:      "success getting list of binary data",
			id:        1234,
			user:      user,
			modelType: TypeBinaryData,
			mockArgs: argsMock{
				id: 1234,
				pd: []PrivateDataInfo{
					{
						ID:       1234,
						UserID:   4321,
						Type:     TypeUsepass,
						View:     "test",
						Data:     []byte("data"),
						DEK:      []byte("key"),
						Metadata: map[string]string{
							"key": "value",
						},
					},
				},
				err: nil,
			},
			expected: []*BinaryData{
				{
					model: model{
						id:    1234,
						view:  "test",
						owner: user,
						encryptedData: &EncryptedData{
							Data: []byte("data"),
							Key:  []byte("key"),
						},
						modelType: TypeUsepass,
						metadata:  map[string]string{
							"key": "value",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pvRepository := NewMockPrivateDataRepository(ctrl)

			manager := binaryDataManager{
				pvRepository: pvRepository,
			}
			ctx := context.Background()
			pvRepository.EXPECT().List(ctx, tt.modelType, tt.user).Return(tt.mockArgs.pd, tt.mockArgs.err)

			actual, err := manager.List(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, actual, len(tt.expected))
		})
	}
}
