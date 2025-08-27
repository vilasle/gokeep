package client

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
)

type ClientRepository interface {
	CreateScheme(ctx context.Context) error
	Close() error
	Save(ctx context.Context, tData model.Type, req SaveRequest) error
	Get(ctx context.Context, tData model.Type, req ...GetRequest) ([]GetResponse, error)
	All(ctx context.Context, tData model.Type) ([]GetResponse, error)
	Delete(ctx context.Context, tData model.Type, req DeleteRequest) error
}

type MetadataValue struct {
	Key   string
	Value string
}

type SaveRequest struct {
	ID         int
	View       string
	ExternalID int
	DEK        string
	Data       string
	Metadata   []MetadataValue
}

type SaveResponse struct {
	ID int
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
	Metadata   []MetadataValue
}

type DeleteRequest struct {
	ID int
}
