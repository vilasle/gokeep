package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// type mockReader

func Test_register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		username, password := "account", "password"
		client.EXPECT().CreateAccount(ctx, username, password)
		exitCode := register(ctx, client, password, []string{username})
		assert.Equal(t, 0, exitCode)
	})

	t.Run("empty login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		password := "password"

		exitCode := register(ctx, client, password, []string{})
		assert.Equal(t, reasonNotFillRequiredArgs, exitCode)
	})

	t.Run("creating failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		username, password := "account", "password"

		client.EXPECT().CreateAccount(ctx, username, password).Return(errors.New("error"))
		exitCode := register(ctx, client, password, []string{username})
		assert.Equal(t, reasonInternalError, exitCode)
	})
}

func Test_login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		username, password := "account", "password"
		client.EXPECT().Login(ctx, username, password)
		exitCode := login(ctx, client, password, []string{username})
		assert.Equal(t, 0, exitCode)
	})

	t.Run("empty login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		password := "password"

		exitCode := login(ctx, client, password, []string{})
		assert.Equal(t, reasonNotFillRequiredArgs, exitCode)
	})

	t.Run("creating failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := NewMockClient(ctrl)

		ctx := context.Background()

		username, password := "account", "password"

		client.EXPECT().Login(ctx, username, password).Return(errors.New("error"))
		exitCode := login(ctx, client, password, []string{username})
		assert.Equal(t, reasonInternalError, exitCode)
	})
}
