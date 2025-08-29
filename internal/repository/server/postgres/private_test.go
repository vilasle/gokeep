package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gokeep/internal/model"
)

func Test_NewPrivateDataRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS entity").WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = NewPrivateDataRepository(db)
	assert.NoError(t, err)

}

func Test_PrivateDataRepository_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO entity").
			WithArgs(req.UserID, req.Type, req.View).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entityId))
		mock.ExpectExec("INSERT INTO data_encrypted").
			WithArgs(entityId, req.Data, req.DEK).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(entityId, k, v).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit()

		id, err := r.Add(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("starting transaction failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin().WillReturnError(errors.New("err"))

		id, err := r.Add(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("adding entity failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO entity").
			WithArgs(req.UserID, req.Type, req.View).
			WillReturnError(errors.New("err"))

		mock.ExpectRollback()

		id, err := r.Add(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("adding encrypted data failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO entity").
			WithArgs(req.UserID, req.Type, req.View).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entityId))
		mock.ExpectExec("INSERT INTO data_encrypted").
			WithArgs(entityId, req.Data, req.DEK).
			WillReturnError(errors.New("err"))

		mock.ExpectRollback()

		id, err := r.Add(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("adding metadata failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO entity").
			WithArgs(req.UserID, req.Type, req.View).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entityId))
		mock.ExpectExec("INSERT INTO data_encrypted").
			WithArgs(entityId, req.Data, req.DEK).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(entityId, k, v).
				WillReturnError(errors.New("err"))
		}

		mock.ExpectRollback()

		id, err := r.Add(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("commit failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO entity").
			WithArgs(req.UserID, req.Type, req.View).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entityId))
		mock.ExpectExec("INSERT INTO data_encrypted").
			WithArgs(entityId, req.Data, req.DEK).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(entityId, k, v).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit().WillReturnError(errors.New("err"))
		mock.ExpectRollback()

		id, err := r.Add(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})
}

func Test_PrivateDataRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE entity").
			WithArgs(req.UserID, req.Type, req.View, req.ID).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))
		mock.ExpectExec("UPDATE data_encrypted").
			WithArgs(req.Data, req.DEK, req.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(req.ID, k, v).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit()

		id, err := r.Update(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("starting transaction failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin().WillReturnError(errors.New("err"))

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("update entity failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE entity").
			WithArgs(req.UserID, req.Type, req.View, req.ID).
			WillReturnError(errors.New("err"))

		mock.ExpectRollback()

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("update encrypted data failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE entity").
			WithArgs(req.UserID, req.Type, req.View, req.ID).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))
		mock.ExpectExec("UPDATE data_encrypted").
			WithArgs(req.Data, req.DEK, req.ID).
			WillReturnError(errors.New("err"))

		mock.ExpectRollback()

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("update metadata failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE entity").
			WithArgs(req.UserID, req.Type, req.View, req.ID).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))
		mock.ExpectExec("UPDATE data_encrypted").
			WithArgs(req.Data, req.DEK, req.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(entityId, k, v).
				WillReturnError(errors.New("err"))
		}

		mock.ExpectRollback()

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("commit failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			ID:       10,
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE entity").
			WithArgs(req.UserID, req.Type, req.View, req.ID).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))
		mock.ExpectExec("UPDATE data_encrypted").
			WithArgs(req.Data, req.DEK, req.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for k, v := range req.Metadata {
			mock.ExpectExec("INSERT INTO metadata").
				WithArgs(req.ID, k, v).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}

		mock.ExpectCommit().WillReturnError(errors.New("err"))
		mock.ExpectRollback()

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

	t.Run("empty req.ID", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()
		req := model.PrivateDataSave{
			UserID:   1,
			Type:     1,
			View:     "test",
			Data:     []byte("data"),
			DEK:      []byte("DEK"),
			Metadata: map[string]string{"key": "value"},
		}

		id, err := r.Update(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, entityId, id)
	})

}

func Test_PrivateDataRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM data_encrypted").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM entity").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectCommit()

		err = r.Delete(ctx, entityId)
		assert.NoError(t, err)
	})

	t.Run("empty id", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 0
		ctx := context.Background()

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})

	t.Run("begin transaction failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin().WillReturnError(errors.New("error"))

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})

	t.Run("delete data_encrypted failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM data_encrypted").
			WithArgs(entityId).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})

	t.Run("delete metadata failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM data_encrypted").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(entityId).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})

	t.Run("delete entity failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM data_encrypted").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM entity").
			WithArgs(entityId).
			WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})

	t.Run("commit failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM data_encrypted").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectExec("DELETE FROM entity").
			WithArgs(entityId).
			WillReturnResult(sqlmock.NewResult(int64(entityId), 1))

		mock.ExpectCommit().WillReturnError(errors.New("error"))
		mock.ExpectRollback()

		err = r.Delete(ctx, entityId)
		assert.Error(t, err)
	})
}

func Test_PrivateDataRepository_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		expected := model.PrivateDataInfo{
			ID:   entityId,
			View: "view",
			Data: []byte("data"),
			DEK:  []byte("dek"),
			Metadata: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}

		mock.ExpectQuery("SELECT").
			WithArgs(entityId).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "view", "data", "dek", "key", "value"}).
					AddRow(entityId, "view", "data", "dek", "key1", "value1").
					AddRow(entityId, "view", "data", "dek", "key2", "value2"))

		data, err := r.Get(ctx, entityId)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})

	t.Run("getting entity failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		expected := model.PrivateDataInfo{}

		mock.ExpectQuery("SELECT").
			WithArgs(entityId).
			WillReturnError(errors.New("error"))

		data, err := r.Get(ctx, entityId)
		assert.Error(t, err)
		assert.Equal(t, expected, data)
	})

	t.Run("error scanning", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		expected := model.PrivateDataInfo{}

		mock.ExpectQuery("SELECT").
			WithArgs(entityId).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "view", "data", "dek", "key", "value"}).
					AddRow("wrong_id", 1, 1, 1, 1, true))

		data, err := r.Get(ctx, entityId)
		assert.Error(t, err)
		assert.Equal(t, expected, data)
	})
}

func Test_PrivateDataRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		entityId := 10
		ctx := context.Background()

		user := &model.User{}

		expected := []model.PrivateDataInfo{
			{
				ID:   entityId,
				View: "view",
				Data: []byte("data"),
				DEK:  []byte("dek"),
				Metadata: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		}

		mock.ExpectQuery("SELECT").
			WithArgs(0, 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "view", "data", "dek", "key", "value"}).
					AddRow(entityId, "view", "data", "dek", "key1", "value1").
					AddRow(entityId, "view", "data", "dek", "key2", "value2"))

		data, err := r.List(ctx, 1, user)
		assert.NoError(t, err)
		assert.Equal(t, expected, data)
	})

	t.Run("getting failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := PrivateDataRepository{db}
		ctx := context.Background()

		user := &model.User{}

		mock.ExpectQuery("SELECT").
			WithArgs(0, 1).
			WillReturnError(errors.New("error"))

		data, err := r.List(ctx, 1, user)
		assert.Error(t, err)
		assert.Nil(t, data)
	})

}
