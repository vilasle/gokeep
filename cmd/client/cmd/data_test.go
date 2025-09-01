package cmd

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/huandu/go-assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryData_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testFile := "test_file"
		content := []byte("some content")

		err := os.WriteFile(testFile, content, os.ModePerm)
		require.NoError(t, err)
		defer os.Remove(testFile)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryAdd.file = testFile
		binaryAdd.name = testFile

		ctx := context.Background()

		meta = append(meta, "key=value")

		client.EXPECT().SaveBinaryData(ctx, content, testFile, 0, map[string]string{"key": "value"}).Return(nil)

		assert.Equal(t, 0, binaryAddHandle(ctx, client))
	})

	t.Run("not filled file path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		ctx := context.Background()

		binaryAdd.file = ""
		binaryAdd.name = ""

		assert.Equal(t, reasonNotFillRequiredArgs, binaryAddHandle(ctx, client))
	})

	t.Run("not filled name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryAdd.file = "test_path"
		binaryAdd.name = ""

		ctx := context.Background()

		assert.Equal(t, reasonNotFillRequiredArgs, binaryAddHandle(ctx, client))
	})

	t.Run("does not existed file", func(t *testing.T) {
		testFile := "test_file"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryAdd.file = testFile
		binaryAdd.name = testFile

		ctx := context.Background()

		assert.Equal(t, reasonInternalError, binaryAddHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		testFile := "test_file"
		content := []byte("some content")

		err := os.WriteFile(testFile, content, os.ModePerm)
		require.NoError(t, err)
		defer os.Remove(testFile)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryAdd.file = testFile
		binaryAdd.name = testFile

		ctx := context.Background()

		meta = append(meta, "key=value")

		client.EXPECT().
			SaveBinaryData(ctx, content, testFile, 0, map[string]string{"key": "value"}).
			Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, binaryAddHandle(ctx, client))
	})

}

func TestBinaryData_Edit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testFile := "test_file"
		content := []byte("some content")

		err := os.WriteFile(testFile, content, os.ModePerm)
		require.NoError(t, err)
		defer os.Remove(testFile)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryEdit.file = testFile
		binaryEdit.name = testFile
		binaryEdit.id = 1

		ctx := context.Background()

		meta = append(meta, "key=value")

		client.EXPECT().SaveBinaryData(ctx, content, testFile, 1, map[string]string{"key": "value"}).Return(nil)

		assert.Equal(t, 0, binaryEditHandle(ctx, client))
	})

	t.Run("not filled file path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		ctx := context.Background()

		binaryEdit.file = ""
		binaryEdit.name = "name"
		binaryEdit.id = 1

		assert.Equal(t, reasonNotFillRequiredArgs, binaryEditHandle(ctx, client))
	})

	t.Run("not filled id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		ctx := context.Background()

		binaryEdit.file = "test_path"
		binaryEdit.name = "name"
		binaryEdit.id = 0

		assert.Equal(t, reasonNotFillRequiredArgs, binaryEditHandle(ctx, client))
	})

	t.Run("not filled name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryEdit.file = "test_path"
		binaryEdit.name = ""
		binaryEdit.id = 1

		ctx := context.Background()

		assert.Equal(t, reasonNotFillRequiredArgs, binaryEditHandle(ctx, client))
	})

	t.Run("does not existed file", func(t *testing.T) {
		testFile := "test_file"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryEdit.file = testFile
		binaryEdit.name = testFile
		binaryEdit.id = 1

		ctx := context.Background()

		assert.Equal(t, reasonInternalError, binaryEditHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		testFile := "test_file"
		content := []byte("some content")

		err := os.WriteFile(testFile, content, os.ModePerm)
		require.NoError(t, err)
		defer os.Remove(testFile)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryEdit.file = testFile
		binaryEdit.name = testFile
		binaryEdit.id = 1

		ctx := context.Background()

		meta = append(meta, "key=value")

		client.EXPECT().
			SaveBinaryData(ctx, content, testFile, 1, map[string]string{"key": "value"}).
			Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, binaryEditHandle(ctx, client))
	})
}

func TestBinaryData_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryGet.id = 1

		ctx := context.Background()

		client.EXPECT().GetBinaryData(ctx, 1).Return(nil)

		assert.Equal(t, 0, binaryGetHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryGet.id = 1

		ctx := context.Background()

		client.EXPECT().GetBinaryData(ctx, 1).Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, binaryGetHandle(ctx, client))
	})
}

func TestBinaryData_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryDelete.id = 1

		ctx := context.Background()

		client.EXPECT().DeleteBinaryData(ctx, 1).Return(nil)

		assert.Equal(t, 0, binaryDeleteHandle(ctx, client))
	})

	t.Run("not filled id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryDelete.id = 0

		ctx := context.Background()

		assert.Equal(t, reasonNotFillRequiredArgs, binaryDeleteHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryDelete.id = 1

		ctx := context.Background()

		client.EXPECT().DeleteBinaryData(ctx, 1).Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, binaryDeleteHandle(ctx, client))
	})
}

func TestTextData_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		textGet.id = 1

		ctx := context.Background()

		client.EXPECT().GetTextData(ctx, 1).Return(nil)

		assert.Equal(t, 0, textGetHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		binaryGet.id = 1

		ctx := context.Background()

		client.EXPECT().GetTextData(ctx, 1).Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, textGetHandle(ctx, client))
	})
}

func TestTextData_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		textDelete.id = 1

		ctx := context.Background()

		client.EXPECT().DeleteTextData(ctx, 1).Return(nil)

		assert.Equal(t, 0, textDeleteHandle(ctx, client))
	})

	t.Run("not filled id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		textDelete.id = 0

		ctx := context.Background()

		assert.Equal(t, reasonNotFillRequiredArgs, textDeleteHandle(ctx, client))
	})

	t.Run("client failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)

		textDelete.id = 1

		ctx := context.Background()

		client.EXPECT().DeleteTextData(ctx, 1).Return(errors.New("error"))

		assert.Equal(t, reasonInternalError, textDeleteHandle(ctx, client))
	})
}
