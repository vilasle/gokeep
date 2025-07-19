package private

import (
	"context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/service"
)

func Test_BankCardService_List(t *testing.T) {

	type argsUserGet struct {
		userID int64
		err    error
	}

	type argsPrivateRepositoryList struct {
		number   string
		cvv      int
		expireAt time.Time
		empty    bool
		err      error
	}

	behaviorUserGet := func(m *MockUserRepository, ctx context.Context, args argsUserGet, user *model.User) {
		m.EXPECT().Get(ctx, args.userID).Return(user, args.err).Times(1)
	}

	behaviorPrivateGet := func(m *MockPrivateDataRepository, ctx context.Context, args argsPrivateRepositoryList, user *model.User, dto []model.PrivateData) {
		d := dto
		if args.empty {
			d = []model.PrivateData{}
		}
		m.EXPECT().List(ctx, model.ModelTypeBankCard, user).Return(d, args.err).Times(1)

	}

	testCases := []struct {
		name           string
		userID         int64
		userRepArgs    argsUserGet
		privateRepArgs argsPrivateRepositoryList
		expected       service.ListPrivateDataResponse
		wantErr        bool
	}{
		{
			name:   "success",
			userID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgs: argsPrivateRepositoryList{
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      nil,
			},
			expected: service.ListPrivateDataResponse{
				Data: []map[string]any{
					{
						"id":   int64(0),
						"card": "******7890",
					},
				},
			},
		},
		{
			name:   "not found user",
			userID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    errors.New("error`"),
			},
			privateRepArgs: argsPrivateRepositoryList{},
			expected:       service.ListPrivateDataResponse{},
			wantErr:        true,
		},
		{
			name:   "error from private repository",
			userID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgs: argsPrivateRepositoryList{
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      errors.New("error"),
			},
			expected: service.ListPrivateDataResponse{
				Data: []map[string]any{
					{
						"id":   int64(0),
						"card": "******7890",
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "repository does not content any data",
			userID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgs: argsPrivateRepositoryList{
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				empty:    true,
				err:      nil,
			},
			expected: service.ListPrivateDataResponse{
				Error: "there are not any bank cards",
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			collector := NewMockRepositoryCollector(ctrl)

			userRep := NewMockUserRepository(ctrl)

			pvRep := NewMockPrivateDataRepository(ctrl)

			collector.EXPECT().User().Return(userRep).Times(1)
			collector.EXPECT().Private().Return(pvRep).Times(4)
			collector.EXPECT().Encryption().Return(nil).Times(4)

			manager := model.NewModelManager(collector)

			user := manager.Users.New("test", "test")

			behaviorUserGet(userRep, ctx, tt.userRepArgs, user)

			if tt.userRepArgs.err == nil {
				bk := &model.BankCard{}
				bk.SetCVV(tt.privateRepArgs.cvv)
				bk.SetNumber(tt.privateRepArgs.number)
				bk.SetExpiration(tt.privateRepArgs.expireAt)

				behaviorPrivateGet(pvRep, ctx, tt.privateRepArgs, user, []model.PrivateData{bk})
			}

			svc := NewBankCardService(manager)

			dto, err := svc.List(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, dto)

		})

	}
}

func Test_BankCardService_Get(t *testing.T) {
	type argsUserGet struct {
		userID int64
		err    error
	}

	type argsPrivateRepositoryGet struct {
		id       int64
		number   string
		cvv      int
		expireAt time.Time
		err      error
	}

	behaviorUserGet := func(m *MockUserRepository, ctx context.Context, args argsUserGet, user *model.User) {
		m.EXPECT().Get(ctx, args.userID).Return(user, args.err).Times(1)
	}

	behaviorPrivateList := func(m *MockPrivateDataRepository, ctx context.Context, args argsPrivateRepositoryGet, user *model.User, dto model.PrivateData) {
		d := dto
		m.EXPECT().Get(ctx, args.id).Return(d, args.err).Times(1)
	}

	testCases := []struct {
		name           string
		userID         int64
		entityID       int64
		userRepArgs    argsUserGet
		privateRepArgs argsPrivateRepositoryGet
		expected       service.PrivateDataResponse
		wantErr        bool
	}{
		{
			name:     "success",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgs: argsPrivateRepositoryGet{
				id:       1,
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      nil,
			},
			expected: service.PrivateDataResponse{
				Fields: map[string]any{
					"id":   int64(0),
					"name": "******7890",
				},
			},
		},
		{
			name:     "user not found",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    errors.New("error`"),
			},
			privateRepArgs: argsPrivateRepositoryGet{},
			expected: service.PrivateDataResponse{
				Error: "user not found",
			},
			wantErr: true,
		},
		{
			name:     "repository error",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgs: argsPrivateRepositoryGet{
				id:       1,
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      errors.New("error"),
			},
			expected: service.PrivateDataResponse{
				Error: "error getting bank card",
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			collector := NewMockRepositoryCollector(ctrl)

			userRep := NewMockUserRepository(ctrl)

			pvRep := NewMockPrivateDataRepository(ctrl)

			collector.EXPECT().User().Return(userRep).Times(1)
			collector.EXPECT().Private().Return(pvRep).Times(4)
			collector.EXPECT().Encryption().Return(nil).Times(4)

			manager := model.NewModelManager(collector)

			user := manager.Users.New("test", "test")

			behaviorUserGet(userRep, ctx, tt.userRepArgs, user)

			if tt.userRepArgs.err == nil {
				bk := &model.BankCard{}
				bk.SetCVV(tt.privateRepArgs.cvv)
				bk.SetNumber(tt.privateRepArgs.number)
				bk.SetExpiration(tt.privateRepArgs.expireAt)

				behaviorPrivateList(pvRep, ctx, tt.privateRepArgs, user, bk)
			}

			svc := NewBankCardService(manager)

			dto, err := svc.Get(ctx, service.GetPrivateData{
				UserID: tt.userID,
				ID:     tt.privateRepArgs.id,
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, dto)

		})

	}
}

func Test_BankCardService_Delete(t *testing.T) {
	type argsUserGet struct {
		userID int64
		err    error
	}

	type argsPrivateRepositoryGet struct {
		id       int64
		number   string
		cvv      int
		expireAt time.Time
		err      error
	}

	behaviorUserGet := func(m *MockUserRepository, ctx context.Context, args argsUserGet, user *model.User) {
		m.EXPECT().Get(ctx, args.userID).Return(user, args.err).Times(1)
	}

	behaviorPrivateGet := func(m *MockPrivateDataRepository, ctx context.Context, args argsPrivateRepositoryGet, dto model.PrivateData) {
		d := dto
		m.EXPECT().Get(ctx, args.id).Return(d, args.err).Times(1)
	}

	testCases := []struct {
		name              string
		userID            int64
		entityID          int64
		userRepArgs       argsUserGet
		privateRepArgsGet argsPrivateRepositoryGet
	}{
		{
			//check only error because without id in entity and will return error
			name:     "got error on delete",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgsGet: argsPrivateRepositoryGet{
				id:       1,
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      nil,
			},
		},
		{
			name:     "got error on getting user",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    errors.New("error"),
			},
			privateRepArgsGet: argsPrivateRepositoryGet{
				id:       1,
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      nil,
			},
		},
		{
			name:     "got error on getting bank card",
			userID:   1,
			entityID: 1,
			userRepArgs: argsUserGet{
				userID: 1,
				err:    nil,
			},
			privateRepArgsGet: argsPrivateRepositoryGet{
				id:       1,
				number:   "1234567890",
				cvv:      123,
				expireAt: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
				err:      errors.New("error"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			collector := NewMockRepositoryCollector(ctrl)

			userRep := NewMockUserRepository(ctrl)

			pvRep := NewMockPrivateDataRepository(ctrl)
			encRep := NewMockEncryptedDataRepository(ctrl)

			collector.EXPECT().User().Return(userRep).Times(1)
			collector.EXPECT().Private().Return(pvRep).Times(4)
			collector.EXPECT().Encryption().Return(encRep).Times(4)

			manager := model.NewModelManager(collector)

			user := manager.Users.New("test", "test")

			behaviorUserGet(userRep, ctx, tt.userRepArgs, user)

			if tt.userRepArgs.err == nil {
				bk := &model.BankCard{}
				bk.SetCVV(tt.privateRepArgsGet.cvv)
				bk.SetNumber(tt.privateRepArgsGet.number)
				bk.SetExpiration(tt.privateRepArgsGet.expireAt)

				behaviorPrivateGet(pvRep, ctx, tt.privateRepArgsGet, bk)
			}

			svc := NewBankCardService(manager)

			err := svc.Delete(ctx, service.DeletePrivateData{
				UserID: tt.userID,
				ID:     tt.privateRepArgsGet.id,
			})
			assert.Error(t, err)
		})

	}
}

func Test_BankCardService_Add(t *testing.T) {
	// type argsUserGet struct {
	// 	userID int64
	// 	err    error
	// }

	// behaviorUserGet := func(m *MockUserRepository, ctx context.Context, args argsUserGet, user *model.User) {
	// 	m.EXPECT().Get(ctx, args.userID).Return(user, args.err).Times(1)
	// }

	// ctx := context.Background()

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// collector := NewMockRepositoryCollector(ctrl)

	// userRep := NewMockUserRepository(ctrl)

	// pvRep := NewMockPrivateDataRepository(ctrl)
	// encRep := NewMockEncryptedDataRepository(ctrl)
	// encryptor := NewMockEncryptor(ctrl)

	// collector.EXPECT().User().Return(userRep).Times(1)
	// collector.EXPECT().Private().Return(pvRep).Times(4)
	// collector.EXPECT().Encryption().Return(encRep).Times(4)

	// manager := model.NewModelManager(collector)

	// user := manager.Users.New("test", "test")

	// pvRep.EXPECT().Add(ctx, gomock.Any()).Return(nil)
	// encRep.EXPECT().Add(ctx, gomock.Any()).Return(nil)

	// behaviorUserGet(userRep, ctx, argsUserGet{1, nil}, user)

	// svc := NewBankCardService(manager)

	// dto := service.AddBankCard{
	// 	UserID:     1,
	// 	Number:     "1234567890",
	// 	CVV:        123,
	// 	Expiration: time.Date(2030, 5, 5, 0, 0, 0, 0, time.UTC),
	// }

	// r, err := svc.Add(ctx, dto, encryptor)

	// assert.NoError(t, err)

	// _ = r
}
