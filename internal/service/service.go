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

type PrivateDataService interface {
	List(ctx context.Context, userID int) (ListPrivateDataResponse, error)
	Get(ctx context.Context, req GetPrivateData) (PrivateDataResponse, error)
	Delete(ctx context.Context, req DeletePrivateData) error
}

// Add(context.Context, PrivateData) error
// Update(context.Context, PrivateData) error
// Delete(context.Context, PrivateData) error
// Get(context.Context, int64) (PrivateData, error)

type LoginPasswordService interface {
	PrivateDataService
	Add(ctx context.Context, req AddLoginPassword, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
	Update(ctx context.Context, req UpdateLoginPassword, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
}

type BankCardService interface {
	PrivateDataService
	Add(ctx context.Context, req AddBankCard, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
	Update(ctx context.Context, req UpdateBankCard, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
}

type TextDataService interface {
	PrivateDataService
	Add(ctx context.Context, req AddTextData, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
	Update(ctx context.Context, req UpdateTextData, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
}

type BinaryDataService interface {
	PrivateDataService
	Add(ctx context.Context, req AddBinaryData, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
	Update(ctx context.Context, req UpdateBinaryData, clientKey encryption.Encoder) (AddingUpdatePrivateDataResponse, error)
}
