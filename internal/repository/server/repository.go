package repository

import "context"

type SessionRepository interface {
	Create(context.Context, CredentialCreate) (int, error)
	Get(ctx context.Context, id int) (CredentialInfo, error)
}

type CredentialCreate struct {
	UserId    int
	PublicKey []byte
}

type CredentialInfo struct {
	ID        int
	UserId    int
	PublicKey []byte
}
