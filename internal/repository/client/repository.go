package client

import "context"

type PrivateDataType int

const (
	TypeLoginPassword PrivateDataType = iota + 1
	TypeBankCard
	TypeTextData
	TypeBinaryData
)

type ClientRepository interface {
	CreateScheme() error
	Save(ctx context.Context, tData PrivateDataType, req SaveRequest) SaveResponse
	Get(ctx context.Context, tData PrivateDataType, req ...GetRequest) ([]GetResponse, error)
	Delete(ctx context.Context, tData PrivateDataType, req DeleteRequest) error
}

type SaveRequest struct {
	ExternalID int
	DEK        []byte
	Data       []byte
}

type SaveResponse struct {
	ID    int
	Error string
}

type GetRequest struct {
	ID int
}

type GetResponse struct {
	ID         int
	ExternalID int
	DEK        []byte
	Data       []byte
	View       string
}

type DeleteRequest struct {
	ID int
}