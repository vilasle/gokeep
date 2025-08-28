package client

import (
	"context"

	"github.com/vilasle/gokeep/internal/model"
)

type ClientRepository interface {
	CreateScheme(ctx context.Context) error
	Close() error
	Save(ctx context.Context, req SaveRequest) error
	Get(ctx context.Context, req GetRequest) ([]GetResponse, error)
	All(ctx context.Context, tData model.Type) ([]GetResponse, error)
	Delete(ctx context.Context, req DeleteRequest) error
	Rewrite(ctx context.Context, req []SaveRequest) error
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
	Type       int
	Data       string
	Metadata   []MetadataValue
}

type SaveResponse struct {
	ID int
}

type GetRequest struct {
	ID   int
	Type int
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
	ID   int
	Type int
}
