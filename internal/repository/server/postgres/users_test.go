package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vilasle/gokeep/internal/model"
)

func Test_NewUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = NewUserRepository(db)
	assert.NoError(t, err)
}

func Test_UserRepository_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	req := model.UserAdd{
		Login:    "test",
		Password: "test",
	}
	mock.ExpectQuery("INSERT INTO users").WithArgs(req.Login, req.Password).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	r := &UserRepository{db}
	assert.NoError(t, err)

	id, err := r.Add(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func Test_UserRepository_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.Background()
		req := 1
		expected := model.UserInfo{
			ID:       1,
			Login:    "test",
			Password: "test",
		}

		mock.ExpectQuery("SELECT id, login, password FROM users").
			WithArgs(req).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "login", "password"}).
					AddRow(1, "test", "test"),
			)

		r := &UserRepository{db}
		assert.NoError(t, err)

		user, err := r.Get(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, user)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.Background()
		req := 1
		expected := model.UserInfo{}

		mock.ExpectQuery("SELECT id, login, password FROM users").
			WithArgs(req).
			WillReturnError(sql.ErrNoRows)

		r := &UserRepository{db}
		assert.NoError(t, err)

		user, err := r.Get(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, expected, user)
	})
}

func Test_UserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	req := model.UserUpdate{
		ID:       1,
		Login:    "test",
		Password: "test",
	}
	mock.ExpectExec("UPDATE users").WithArgs(req.Login, req.Password, req.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	r := &UserRepository{db}
	assert.NoError(t, err)

	err = r.Update(ctx, req)
	assert.NoError(t, err)
}

func Test_UserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	req := 1
	mock.ExpectExec("DELETE FROM users").WithArgs(req).WillReturnResult(sqlmock.NewResult(1, 1))

	r := &UserRepository{db}
	assert.NoError(t, err)

	err = r.Delete(ctx, req)
	assert.NoError(t, err)
}

func Test_UserRepository_Find(t *testing.T) {
	t.Run("user is found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.Background()
		login := "test"
		expected := model.UserInfo{
			ID:       1,
			Login:    "test",
			Password: "test",
		}

		mock.ExpectQuery("SELECT id, login, password FROM users").
			WithArgs(login).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "login", "password"}).
					AddRow(1, "test", "test"),
			)

		r := &UserRepository{db}
		assert.NoError(t, err)

		ui, err := r.Find(ctx, login)
		assert.NoError(t, err)
		assert.Equal(t, expected, ui)
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		ctx := context.Background()
		login := "test"
		expected := model.UserInfo{}

		mock.ExpectQuery("SELECT id, login, password FROM users").
			WithArgs(login).
			WillReturnError(sql.ErrNoRows)

		r := &UserRepository{db}
		assert.NoError(t, err)

		ui, err := r.Find(ctx, login)
		assert.Error(t, err)
		assert.Equal(t, expected, ui)
	})

}
