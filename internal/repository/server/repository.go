package repository

import "context"

//SessionRepository is interface for work with users' sessions
type SessionRepository interface {
	Create(context.Context, CredentialCreate) (int, error)
	Get(ctx context.Context, id int) (CredentialInfo, error)
}

type CredentialCreate struct {
	UserID    int
	PublicKey []byte
}

type CredentialInfo struct {
	ID        int
	UserID    int
	PublicKey []byte
}
