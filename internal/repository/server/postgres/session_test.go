package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	repository "github.com/vilasle/gokeep/internal/repository/server"
)

func Test_NewSessionRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS session").WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = NewSessionRepository(db)

	assert.NoError(t, err)
}

func Test_SessionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := SessionRepository{db}
	ctx := context.Background()
	req := repository.CredentialCreate{
		UserID:    1,
		PublicKey: []byte("test"),
	}
	expected := 2

	mock.ExpectQuery("INSERT INTO session").
		WithArgs(req.UserID, req.PublicKey).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	resp, err := r.Create(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func Test_SessionRepository_Get(t *testing.T) {
	t.Run("session is found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := SessionRepository{db}
		ctx := context.Background()
		req := 2
		expected := repository.CredentialInfo{
			ID:        2,
			UserID:    1,
			PublicKey: []byte("test"),
		}

		mock.ExpectQuery("SELECT id, user_id, public_key FROM session").
			WithArgs(req).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "public_key"}).AddRow(2, 1, "test"))

		resp, err := r.Get(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("session is not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		r := SessionRepository{db}
		ctx := context.Background()
		req := 2
		expected := repository.CredentialInfo{}

		mock.ExpectQuery("SELECT id, user_id, public_key FROM session").
			WithArgs(req).WillReturnError(sql.ErrNoRows)

		resp, err := r.Get(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, expected, resp)
	})

}
