package model

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUser_Save(t *testing.T) {
	behavior := func(m *MockUserRepository, user *User, err error) {
		if user.id == 0 {
			m.EXPECT().
				Add(gomock.Any(), UserAdd{Login: user.login, Password: user.password}).
				Return(1234, err)
		} else {
			m.EXPECT().
				Update(gomock.Any(), UserUpdate{ID: user.id, Login: user.login, Password: user.password}).
				Return(err)
		}
	}
	testCases := []struct {
		name     string
		id       int
		login    string
		password string
		success  bool
	}{
		{
			name:     "new user, need to create",
			id:       0,
			login:    "test",
			password: "password",
			success:  true,
		},
		{
			name:     "new user, get error for repository",
			id:       0,
			login:    "test",
			password: "password",
			success:  false,
		},
		{
			name:     "existed user, need to update",
			id:       1234,
			login:    "test",
			password: "password",
			success:  true,
		},
		{
			name:     "existed user, get error for repository",
			id:       1234,
			login:    "test",
			password: "password",
			success:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockUserRepository(ctrl)

			var testError error
			if !tt.success {
				testError = errors.New("error")
			}

			user := newUser(tt.login, tt.password, r)
			user.id = tt.id
			ctx := context.Background()
			behavior(r, user, testError)

			err := user.Save(ctx)
			assert.Equal(t, tt.success, err == nil)
		})
	}
}

func TestUser_Delete(t *testing.T) {
	behavior := func(m *MockUserRepository, user *User, err error) {
		if user.id > 0 {
			m.EXPECT().Delete(gomock.Any(), user.id).Return(err)
		}
	}
	testCases := []struct {
		name     string
		id       int
		login    string
		password string
		success  bool
	}{
		{
			name:     "new user, get error",
			id:       0,
			login:    "test",
			password: "password",
			success:  false,
		},
		{
			name:     "existed user, need to delete",
			id:       1234,
			login:    "test",
			password: "password",
			success:  true,
		},
		{
			name:     "existed user, get error for repository",
			id:       1234,
			login:    "test",
			password: "password",
			success:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockUserRepository(ctrl)

			var testError error
			if !tt.success {
				testError = errors.New("error")
			}

			user := newUser(tt.login, tt.password, r)
			user.id = tt.id
			ctx := context.Background()
			behavior(r, user, testError)

			err := user.Delete(ctx)
			assert.Equal(t, tt.success, err == nil)
		})
	}
}

func TestUser_PasswordIsValid(t *testing.T) {
	testCases := []struct {
		name          string
		login         string
		password      string
		checkPassword string
		match         bool
	}{
		{
			name:          "valid password",
			login:         "test",
			password:      "password",
			checkPassword: "password",
			match:         true,
		},
		{
			name:          "invalid password",
			login:         "test",
			password:      "password",
			checkPassword: "another_password",
			match:         false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := NewMockUserRepository(ctrl)

			user := newUser(tt.login, tt.password, r)
			assert.Equal(t, tt.match, user.PasswordIsValid(tt.checkPassword))
		})
	}
}

func Test_hash256(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "valid password",
			password: "password",
			expected: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		},
		{
			name:     "empty password",
			password: "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "valid password",
			password: "dsgk43;le4.rmecd.vkc43pelr.ffd",
			expected: "2a5cabe5aee156a4d82c6bbf421e657a71a84b3624907d4ef545a0323051c1a7",
		},
		{
			name:     "valid password",
			password: "ffd;vlb,43o43e0-3oefkdvm,hgui3ex",
			expected: "62380a6d5f11ffe38bc27254c39a5873a897f3bc5ece23c8d93678d500361e67",
		},
		{
			name:     "valid password",
			password: "432l,vf.2403edc_3421-_213w2	34ed",
			expected: "815bb53b06bef79608dae1eefdb09713c52febe71e4cf3e38ce734828d777bfd",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := hash256(tt.password)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_findUserByLogin(t *testing.T) {
	type mockArgs struct {
		input  string
		output UserInfo
		err error
	}

	testCases := []struct {
		name     string
		login    string
		password string
		mockArgs mockArgs
		expected *User
		wantErr bool
	}{
		{
			name:     "valid user",
			login:    "test",
			password: "password",
			expected: &User{
				id:       1,
				login:    "test",
				password: "5e884898da28047151d42d8",
			},
			wantErr: false,
			mockArgs: mockArgs{
				input: "test",
				output: UserInfo{
					ID:       1,
					Login:    "test",
					Password: "5e884898da28047151d42d8",
				},
				err: nil,
			},
		},
		{
			name:     "repository error",
			login:    "test",
			password: "password",
			expected: nil,
			wantErr: true,
			mockArgs: mockArgs{
				input: "test",
				output: UserInfo{},
				err: errors.New("repository error"),
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			r := NewMockUserRepository(ctrl)
			r.EXPECT().Find(ctx, tt.mockArgs.input).Return(tt.mockArgs.output, tt.mockArgs.err)

			actual, err := findUserByLogin(ctx, tt.login, r)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.id, actual.ID())
			assert.Equal(t, tt.expected.login, actual.login)
			assert.Equal(t, tt.expected.password, actual.password)
		})
	}
}

func Test_getUserByID(t *testing.T) {
	type mockArg struct {
		input  int
		output UserInfo
		err    error
	}

	testCases := []struct {
		name      string
		input     int
		output    *User
		wantError bool
		mockArg
	}{
		{
			name:  "valid user",
			input: 1,
			output: &User{
				id:       1,
				login:    "test",
				password: "password",
			},
			wantError: false,
			mockArg: mockArg{
				input: 1,
				output: UserInfo{
					ID:       1,
					Login:    "test",
					Password: "password",
				},
				err: nil,
			},
		},
		{
			name:      "repository failed",
			input:     1,
			output:    nil,
			wantError: true,
			mockArg: mockArg{
				input:  1,
				output: UserInfo{},
				err:    errors.New("repository failed"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := NewMockUserRepository(ctrl)

			r.EXPECT().Get(gomock.Any(), tt.mockArg.input).Return(tt.mockArg.output, tt.mockArg.err)

			ctx := context.Background()
			actual, err := getUserByID(ctx, tt.input, r)

			if tt.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.output.id, actual.id)
			assert.Equal(t, tt.output.password, actual.password)
			assert.Equal(t, tt.output.login, actual.login)
		})
	}

}
