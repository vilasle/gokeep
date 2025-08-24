package client

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
)





type ClientRepository interface {
	CreateScheme(ctx context.Context) error
	Close() error
	Save(ctx context.Context, tData model.Type, req SaveRequest) SaveResponse
	Get(ctx context.Context, tData model.Type, req ...GetRequest) ([]GetResponse, error)
	Delete(ctx context.Context, tData model.Type, req DeleteRequest) error
}

type SaveRequest struct {
	ID         int
	View       string
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
