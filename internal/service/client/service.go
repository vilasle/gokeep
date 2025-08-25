package client

import "context"

type AuthService interface {
	CreateAccount(ctx context.Context, accountName, password string) error
	Login(ctx context.Context, accountName, password string, publicKey []byte) (credential []byte, err error)
}

type LoginPasswordDataService interface {
	Save(context.Context, LoginPasswordSaveRequest) (SaveResponse, error)
	Get(context.Context, GetRequest) ([]EncryptedData, error)
	Delete(context.Context, DeleteRequest) error
}

type BankCardDataService interface {
	Save(context.Context, BankCardSaveRequest) (SaveResponse, error)
	Get(context.Context, GetRequest) ([]EncryptedData, error)
	Delete(context.Context, DeleteRequest) error
}

type TextDataDataService interface {
	Save(context.Context, TextDataSaveRequest) (SaveResponse, error)
	Get(context.Context, GetRequest) ([]EncryptedData, error)
	Delete(context.Context, DeleteRequest) error
}

type BinaryDataDataService interface {
	Save(context.Context, BinaryDataSaveRequest) (SaveResponse, error)
	Get(context.Context, GetRequest) ([]EncryptedData, error)
	Delete(context.Context, DeleteRequest) error
}
