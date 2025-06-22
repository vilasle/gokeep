package model

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewModelManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockRepositoryCollector(ctrl)

	ur := NewMockUserRepository(ctrl)

	m.EXPECT().User().Return(ur)

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
		Login:    "test",
		Password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
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
		Login:    "test",
		Password: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		r:        repository,
	}

	repository.EXPECT().Get("test").Return(expected, nil)

	manager := userManager{
		repository: repository,
	}

	u, _ := manager.FindByLogin("test")
	assert.Equal(t, expected, u)
}
