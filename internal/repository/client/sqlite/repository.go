package sqlite

import (
	"context"

	"github.com/vilasle/gokeep/internal/repository/client"
)

// TODO implement it
type ClientRepository struct {
	// db *sql.DB
	dbPath string
}

func NewSQLiteClient(dbPath string) *ClientRepository {
	return &ClientRepository{dbPath}
}

func (r *ClientRepository) CreateScheme() error {
	return nil
}

func (r *ClientRepository) Close() error {
	return nil
}

func (r *ClientRepository) Save(ctx context.Context, tData client.PrivateDataType, req client.SaveRequest) client.SaveResponse {
	return client.SaveResponse{
		ID: 1,
	}
}

func (r *ClientRepository) Get(ctx context.Context, tData client.PrivateDataType, req ...client.GetRequest) ([]client.GetResponse, error) {
	return []client.GetResponse{
		{
			ID: 1,
		},
	}, nil
}

func (r *ClientRepository) Delete(ctx context.Context, tData client.PrivateDataType, req client.DeleteRequest) error {
	return nil
}
