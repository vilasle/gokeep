package client

import "context"

type AuthService interface {
	CreateAccount(accountName, password string, publicKey []byte) error
	Login(accountName, password string) (credential []byte, err error)
}

type LoginPasswordDataService interface {
	Save(context.Context, LoginPasswordSaveRequest) LoginPasswordSaveResponse
	Get(context.Context, LoginPasswordGetRequest) LoginPasswordGetResponse
	Delete(context.Context, LoginPasswordDeleteRequest) LoginPasswordDeleteResponse
}

type BankCardDataService interface {
	Save(context.Context, BankCardSaveRequest) BankCardSaveResponse
	Get(context.Context, BankCardGetRequest) BankCardGetResponse
	Delete(context.Context, BankCardDeleteRequest) BankCardDeleteResponse
}

type TextDataDataService interface {
	Save(context.Context, TextDataSaveRequest) TextDataSaveResponse
	Get(context.Context, TextDataGetRequest) TextDataGetResponse
	Delete(context.Context, TextDataDeleteRequest) TextDataDeleteResponse
}

type BinaryDataDataService interface {
	Save(context.Context, BinaryDataSaveRequest) BinaryDataSaveResponse
	Get(context.Context, BinaryDataGetRequest) BinaryDataGetResponse
	Delete(context.Context, BinaryDataDeleteRequest) BinaryDataDeleteResponse
}
