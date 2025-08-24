package postgres

import (
	"context"
	"database/sql"

	repository "github.com/vilasle/gokeep/internal/repository/server"
)

var _ repository.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) (*SessionRepository, error) {
	repository := &SessionRepository{db: db}

	return repository, repository.initSchema()
}

func (r *SessionRepository) Create(ctx context.Context, credential repository.CredentialCreate) (id int, err error) {
	txt := "INSERT INTO session (user_id, public_key) VALUES($1, $2) RETURNING id"
	err = r.db.QueryRowContext(ctx, txt, credential.UserId, credential.PublicKey).Scan(&id)
	return
}

func (r *SessionRepository) Get(ctx context.Context, id int) (repository.CredentialInfo, error) {
	txt := "SELECT id, user_id, public_key FROM session WHERE id = $1"

	row := r.db.QueryRowContext(ctx, txt, id)
	var credential repository.CredentialInfo
	err := row.Scan(&credential.ID, &credential.UserId, &credential.PublicKey)

	if err == sql.ErrNoRows {
		return repository.CredentialInfo{}, repository.ErrNotFound
	}

	return credential, err
}

func (r *SessionRepository) initSchema() error {
	txt := `
	CREATE TABLE IF NOT EXISTS session (
		id bigint GENERATED ALWAYS AS IDENTITY,
		user_id bigint NOT NULL,
		public_key BYTEA NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	CREATE INDEX IF NOT EXISTS session_user_idx ON users (id);
	`

	_, err := r.db.Exec(txt)
	return err
}
