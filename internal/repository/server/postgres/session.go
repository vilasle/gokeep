package postgres

import (
	"context"
	"database/sql"

	repository "github.com/vilasle/gokeep/internal/repository/server"
)

var _ repository.SessionRepository = (*SessionRepository)(nil)

// SessionRepository is a repository for sessions
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *sql.DB) (*SessionRepository, error) {
	repository := &SessionRepository{db: db}

	return repository, repository.initSchema()
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, credential repository.CredentialCreate) (id int, err error) {
	txt := "INSERT INTO session (user_id, public_key) VALUES($1, $2) RETURNING id"
	err = r.db.QueryRowContext(ctx, txt, credential.UserID, credential.PublicKey).Scan(&id)
	return
}

// Get gets a session by id
func (r *SessionRepository) Get(ctx context.Context, id int) (credential repository.CredentialInfo, err error) {
	txt := "SELECT id, user_id, public_key FROM session WHERE id = $1"

	err = r.db.QueryRowContext(ctx, txt, id).Scan(&credential.ID, &credential.UserID, &credential.PublicKey)
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
