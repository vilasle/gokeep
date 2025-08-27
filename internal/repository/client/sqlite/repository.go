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

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//save main entity
	iq := sqlbuilder.InsertInto("private_data").
		Cols("external_id", "type", "dek", "data", "view").
		Values(req.ExternalID, tData, req.DEK, req.Data, req.View).
		Returning("id")
	iq.SetFlavor(sqlbuilder.SQLite)

	q, args := iq.Build()

	result, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if err := r.addOrUpdateMetadataByOwner(ctx, tx, int(id), req.Metadata); err != nil {
		return err
	}

	return tx.Commit()
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

func (r *ClientRepository) All(ctx context.Context, tData model.Type) ([]client.GetResponse, error) {
	sq := sqlbuilder.Select("id", "external_id", "view").From("private_data")
	sq = sq.Where(sq.Equal("type", tData))

	q, args := sq.Build()

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]client.GetResponse, 0)

	for rows.Next() {
		var response client.GetResponse
		if err := rows.Scan(&response.ID, &response.ExternalID, &response.View); err != nil {
			return nil, err
		}
		result = append(result, response)
	}

	return result, nil
}

func (r *ClientRepository) Delete(ctx context.Context, tData model.Type, req client.DeleteRequest) error {

	if err := r.deleteMetadataByOwner(ctx, req.ID); err != nil {
		return err
	}

	dq := sqlbuilder.DeleteFrom("private_data")
	dq = dq.Where(dq.Equal("type", tData))

	if req.ID > 0 {
		dq = dq.Where(dq.Equal("id", req.ID))
	}

	q, args := dq.Build()

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *ClientRepository) deleteMetadataByOwner(ctx context.Context, ownerId int) error {
	dq := sqlbuilder.DeleteFrom("metadata")
	dq = dq.Where(dq.Equal("owner_id", ownerId))

	q, args := dq.Build()

	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *ClientRepository) addOrUpdateMetadataByOwner(ctx context.Context, tx *sql.Tx, ownerId int, metadata []client.MetadataValue) error {
	for _, m := range metadata {
		id, err := getMetadata(ctx, tx, ownerId, m.Key)
		if err != nil {
			return err
		}
		if id == 0 {
			if err := addMetadata(ctx, tx, ownerId, m.Key, m.Value); err != nil {
				return err
			}
		} else {
			if err := updateMetadata(ctx, tx, ownerId, m.Key, m.Value); err != nil {
				return err
			}
		}
	}
	return nil
}

func getMetadata(ctx context.Context, tx *sql.Tx, ownerId int, key string) (id int, err error) {
	q := `SELECT id FROM metadata WHERE owner_id = $1 AND key = $2`

	err = tx.QueryRowContext(ctx, q, ownerId, key).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, err
}

func addMetadata(ctx context.Context, tx *sql.Tx, ownerId int, key, value string) error {
	q := `INSERT INTO metadata (owner_id, key, value) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, q, ownerId, key, value)
	return err
}

func updateMetadata(ctx context.Context, tx *sql.Tx, ownerId int, key, value string) error {
	q := `UPDATE metadata SET value = $1 WHERE owner_id = $2 AND id = $3`

	_, err := tx.ExecContext(ctx, q, value, ownerId, key)
	return err
}
