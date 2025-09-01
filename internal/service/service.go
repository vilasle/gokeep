package service

import (
	"context"

	"github.com/vilasle/gokeep/internal/encryption"
)

// AuthService - create new users and check existing users, and check tokens
type AuthService interface {
	//Login check user's login and password and return credential token
	Login(ctx context.Context, req RegisterLoginUser) (string, error)
	//Register create user if it does not exist	and return credential token
	Register(ctx context.Context, req RegisterLoginUser) (error)
	GetSessionByCredentialToken(ctx context.Context, token string) (SessionInfo, error)
}

//PrivateDataService - common interface for private data
type PrivateDataService interface {
	List(ctx context.Context, userID int, clientKey encryption.Encoder) (ListPrivateDataResponse, error)
	Get(ctx context.Context, req GetPrivateData, clientKey encryption.Encoder) (PrivateDataResponse, error)
	Delete(ctx context.Context, req DeletePrivateData) error
}

// Add(context.Context, PrivateData) error
// Update(context.Context, PrivateData) error
// Delete(context.Context, PrivateData) error
// Get(context.Context, int64) (PrivateData, error)

//LoginPasswordService - interface for write/update operation with LoginPassword data 
type LoginPasswordService interface {
	PrivateDataService
	Add(ctx context.Context, req AddLoginPassword, clientKey encryption.Encoder) (PrivateDataResponse, error)
	Update(ctx context.Context, req UpdateLoginPassword, clientKey encryption.Encoder) (PrivateDataResponse, error)
}

//BankCardService - interface for write/update operation with LoginPassword data
type BankCardService interface {
	PrivateDataService
	Add(ctx context.Context, req AddBankCard, clientKey encryption.Encoder) (PrivateDataResponse, error)
	Update(ctx context.Context, req UpdateBankCard, clientKey encryption.Encoder) (PrivateDataResponse, error)
}

//TextDataService - interface for write/update operation with LoginPassword data
type TextDataService interface {
	PrivateDataService
	Add(ctx context.Context, req AddTextData, clientKey encryption.Encoder) (PrivateDataResponse, error)
	Update(ctx context.Context, req UpdateTextData, clientKey encryption.Encoder) (PrivateDataResponse, error)
}

//BinaryDataService - interface for write/update operation with LoginPassword data
type BinaryDataService interface {
	PrivateDataService
	Add(ctx context.Context, req AddBinaryData, clientKey encryption.Encoder) (PrivateDataResponse, error)
	Update(ctx context.Context, req UpdateBinaryData, clientKey encryption.Encoder) (PrivateDataResponse, error)
}
