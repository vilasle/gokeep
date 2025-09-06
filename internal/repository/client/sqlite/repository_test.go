package sqlite

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vilasle/gokeep/internal/repository/client"
)

func Test_NewSQLiteClient(t *testing.T) {

	localDbPath := "local.db"

	client, err := NewSQLiteClient(localDbPath)

	require.NoError(t, err)

	require.NoError(t, client.Close())

	stat, err := os.Stat(localDbPath)
	require.NoError(t, err)
	require.NotNil(t, stat)

	os.RemoveAll(localDbPath)
}

func Test_CreateSchema(t *testing.T) {
	t.Run("create schema", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectExec("PRAGMA").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		client := &ClientRepository{db: db}

		ctx := context.Background()

		err = client.CreateScheme(ctx)
		require.NoError(t, err)
	})

	t.Run("error pragma", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectExec("PRAGMA").WillReturnError(errors.New("error"))
		client := &ClientRepository{db: db}

		ctx := context.Background()

		err = client.CreateScheme(ctx)
		require.Error(t, err)
	})

	t.Run("error create private_data", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectExec("PRAGMA").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS private_data").WillReturnError(errors.New("error"))
		client := &ClientRepository{db: db}

		ctx := context.Background()

		err = client.CreateScheme(ctx)
		require.Error(t, err)
	})

	t.Run("error create metadata", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectExec("PRAGMA").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS metadata").WillReturnError(errors.New("error"))
		client := &ClientRepository{db: db}

		ctx := context.Background()

		err = client.CreateScheme(ctx)
		require.Error(t, err)
	})
}

func Test_Save(t *testing.T) {
	t.Run("add new entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.SaveRequest{
			View:       "test",
			ExternalID: 1,
			DEK:        "key",
			Type:       1,
			Data:       "data",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery("SELECT id FROM metadata").
			WithArgs(1, "key").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO metadata").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rep.Save(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("update existing entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.SaveRequest{
			ID:         1,
			View:       "test",
			ExternalID: 1,
			DEK:        "key",
			Type:       1,
			Data:       "data",
			Metadata: []client.MetadataValue{
				{
					Key:   "key",
					Value: "value",
				},
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery("SELECT id FROM metadata").
			WithArgs(1, "key").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec("UPDATE metadata").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = rep.Save(ctx, req)
		assert.NoError(t, err)
	})
}

func Test_Delete(t *testing.T) {
	t.Run("delete entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.DeleteRequest{
			ID:   1,
			Type: 1,
		}

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("DELETE FROM private_data").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = rep.Delete(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("error deleting metadata", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.DeleteRequest{
			ID:   1,
			Type: 1,
		}

		mock.ExpectExec("DELETE FROM metadata").
			WithArgs(1).WillReturnError(errors.New("error"))

		err = rep.Delete(ctx, req)
		assert.Error(t, err)
	})

}

func Test_Rewrite(t *testing.T) {
	t.Run("success rewriting entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := []client.SaveRequest{
			{
				ID:   1,
				Type: 1,
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.ExpectExec("DELETE FROM metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = rep.Rewrite(ctx, req)
		assert.NoError(t, err)
	})
	
	t.Run("deleting all data failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := []client.SaveRequest{
			{
				ID:   1,
				Type: 1,
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.ExpectExec("DELETE FROM metadata").WillReturnError(errors.New("error"))

		err = rep.Rewrite(ctx, req)
		assert.Error(t, err)
	})

	t.Run("opening transaction failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := []client.SaveRequest{
			{
				ID:   1,
				Type: 1,
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.ExpectExec("DELETE FROM metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectBegin().WillReturnError(errors.New("error"))
		mock.ExpectRollback()

		err = rep.Rewrite(ctx, req)
		assert.Error(t, err)
	})

	t.Run("insert private data failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := []client.SaveRequest{
			{
				ID:   1,
				Type: 1,
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.ExpectExec("DELETE FROM metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO private_data").WillReturnError(errors.New("error"))
		
		mock.ExpectRollback()

		err = rep.Rewrite(ctx, req)
		assert.Error(t, err)
	})

	t.Run("insert metadata failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := []client.SaveRequest{
			{
				ID:   1,
				Type: 1,
				Metadata: []client.MetadataValue{
					{
						Key:   "key",
						Value: "value",
					},
				},
			},
		}

		mock.ExpectExec("DELETE FROM metadata").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO private_data").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO metadata").WillReturnError(errors.New("error"))

		mock.ExpectRollback()

		err = rep.Rewrite(ctx, req)
		assert.Error(t, err)
	})



}

func Test_Get(t *testing.T) {
	t.Run("success getting entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.GetRequest{
			ID:   1,
			Type: 1,
		}

		mock.ExpectQuery(`SELECT`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "external_id", "dek", "data", "view", "key", "value"}).
					AddRow(1, 1, "dek", "data", "view", "key", "value"))

		res, err := rep.Get(ctx, req)

		require.NoError(t, err)

		ent := res[0]
		assert.Equal(t, 1, ent.ID)
		assert.Equal(t, 1, ent.ExternalID)
		assert.Equal(t, []byte("data"), ent.Data)
		assert.Equal(t, []byte("dek"), ent.DEK)
		assert.Equal(t, "view", ent.View)
		assert.Equal(t, []client.MetadataValue{{Key: "key", Value: "value"}}, ent.Metadata)
	})

	t.Run("error getting entity", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()
		req := client.GetRequest{
			ID:   1,
			Type: 1,
		}

		mock.ExpectQuery(`SELECT`).WillReturnError(errors.New("error"))

		_, err = rep.Get(ctx, req)

		require.Error(t, err)

	})

}

func Test_All(t *testing.T) {

	t.Run("success getting all entities", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()

		mock.ExpectQuery(`SELECT`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "external_id", "dek", "data", "view", "key", "value"}).
					AddRow(1, 1, "dek", "data", "view", "key", "value").
					AddRow(1, 1, "dek", "data", "view", "key1", "value1"))

		res, err := rep.All(ctx, 1)

		require.NoError(t, err)

		ent := res[0]
		assert.Equal(t, 1, ent.ID)
		assert.Equal(t, 1, ent.ExternalID)
		assert.Equal(t, []byte("data"), ent.Data)
		assert.Equal(t, []byte("dek"), ent.DEK)
		assert.Equal(t, "view", ent.View)
		assert.Equal(t, []client.MetadataValue{{Key: "key", Value: "value"}, {Key: "key1", Value: "value1"}}, ent.Metadata)
	})
	t.Run("error getting all entities", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rep := &ClientRepository{db: db}
		ctx := context.Background()

		mock.ExpectQuery(`SELECT`).WillReturnError(errors.New("error"))

		_, err = rep.All(ctx, 1)

		require.Error(t, err)

	})

}
