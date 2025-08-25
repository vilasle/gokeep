package sqlite

import (
	"context"
	"database/sql"
	"os"

	"github.com/huandu/go-sqlbuilder"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/repository/client"
)

type ClientRepository struct {
	db *sql.DB
}

func NewSQLiteClient(dbPath string) (*ClientRepository, error) {
	if _, err := os.Stat(dbPath); err != nil {
		if fd, err := os.Create(dbPath); err == nil {
			fd.Close()
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &ClientRepository{db: db}, nil
}

func (r *ClientRepository) CreateScheme(ctx context.Context) error {
	//enable checks for foreign keys
	if _, err := r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS private_data ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		external_id INTEGER, 
		type INTEGER, 
		dek TEXT, 
		data TEXT, 
		view TEXT
	);`); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		owner_id INTEGER, 
		key TEXT, 
		value TEXT, 
		FOREIGN KEY (owner_id) REFERENCES private_data (id));
	`); err != nil {
		return err
	}

	return nil
}

func (r *ClientRepository) Close() error {
	return r.db.Close()
}

func (r *ClientRepository) Save(ctx context.Context, tData model.Type, req client.SaveRequest) error {
	if req.ID > 0 {
		return r.update(ctx, tData, req, req.ID)
	}
	//search by external id and if there is record then update that
	id, err := r.searchByExternalID(req.ExternalID)
	if err != nil {
		return err
	}

	if id == 0 {
		return r.add(ctx, tData, req)
	} else {
		return r.update(ctx, tData, req, id)
	}
}

func (r *ClientRepository) add(ctx context.Context,
	tData model.Type, req client.SaveRequest) error {

	iq := sqlbuilder.InsertInto("private_data").
		Cols("external_id", "type", "dek", "data", "view").
		Values(req.ExternalID, tData, req.DEK, req.Data, req.View).
		Returning("id")
	iq.SetFlavor(sqlbuilder.SQLite)

	q, args := iq.Build()

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *ClientRepository) update(ctx context.Context,
	tData model.Type, req client.SaveRequest, id int) error {

	uq := sqlbuilder.Update("private_data")
	uq.Set(
		uq.Equal("dek", req.DEK),
		uq.Equal("data", req.Data),
		uq.Equal("view", req.View),
		uq.Equal("type", tData),
	)
	uq.SetFlavor(sqlbuilder.SQLite)

	q, args := uq.Where(uq.Equal("id", id)).Build()

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *ClientRepository) searchByExternalID(externalID int) (int, error) {
	query := sqlbuilder.
		Select("id").
		From("private_data")

	q, args := query.
		Where(query.Equal("external_id", externalID)).
		Build()

	row := r.db.QueryRow(q, args...)
	var id int
	if err := row.Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
	}
	return id, nil
}

func (r *ClientRepository) Get(ctx context.Context, tData model.Type, req ...client.GetRequest) ([]client.GetResponse, error) {
	sq := sqlbuilder.Select("id", "external_id", "dek", "data", "view").From("private_data")
	sq = sq.Where(sq.Equal("type", tData))

	for _, r := range req {
		if r.ID > 0 {
			sq.Where(sq.Equal("id", r.ID))
		}
	}

	q, args := sq.Build()

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]client.GetResponse, 0)

	for rows.Next() {
		var response client.GetResponse
		if err := rows.Scan(&response.ID, &response.ExternalID, &response.DEK, &response.Data, &response.View); err != nil {
			return nil, err
		}
		result = append(result, response)
	}

	return result, nil
}

func (r *ClientRepository) Delete(ctx context.Context, tData model.Type, req client.DeleteRequest) error {
	dq := sqlbuilder.DeleteFrom("private_data")
	dq = dq.Where(dq.Equal("type", tData))

	if req.ID > 0 {
		dq = dq.Where(dq.Equal("id", req.ID))
	}

	q, args := dq.Build()

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}
