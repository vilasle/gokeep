package postgres

import (
	"context"
	"database/sql"

	"github.com/vilasle/gokeep/internal/model"
)

var _ model.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	repository := &UserRepository{db: db}

	return repository, repository.initSchema()
}

func (r *UserRepository) Add(ctx context.Context, user model.UserAdd) (id int, err error) {
	err = r.db.QueryRowContext(
		ctx, `INSERT INTO users (login, password) VALUES($1, $2) RETURNING id`,
		user.Login, user.Password).
		Scan(&id)

	return id, err
}

func (r *UserRepository) Get(ctx context.Context, id int) (user model.UserInfo, err error) {
	err = r.db.QueryRowContext(
		ctx, `SELECT id, login, password FROM users WHERE id = $1`,
		id).
		Scan(&user.ID, &user.Login, &user.Password)

	if err == sql.ErrNoRows {
		return model.UserInfo{}, model.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user model.UserUpdate) error {
	_, err := r.db.ExecContext(
		ctx, `UPDATE users SET login = $1, password = $2 WHERE id = $3`,
		user.Login, user.Password, user.ID)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(
		ctx, `DELETE FROM users WHERE id = $1`,
		id)

	return err
}

func (r *UserRepository) Find(ctx context.Context, login string) (user model.UserInfo, err error) {
	err = r.db.QueryRowContext(
		ctx, `SELECT id, login, password FROM users WHERE login = $1`,
		login).
		Scan(&user.ID, &user.Login, &user.Password)

	if err == sql.ErrNoRows {
		return model.UserInfo{}, model.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) initSchema() error {
	_, err := r.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id bigint GENERATED ALWAYS AS IDENTITY UNIQUE,
		login VARCHAR(100) NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON users (login);`)
	return err
}
