package sqlite

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vilasle/gokeep/internal/model"
	"github.com/vilasle/gokeep/internal/repository/client"
)

// ClientRepository is a repository for clients
type ClientRepository struct {
	db *sql.DB
}

// NewSQLiteClient creates a new SQLite client
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


// CreateScheme creates the database scheme
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

// Close closes the database connection
func (r *ClientRepository) Close() error {
	return r.db.Close()
}

//Save - add new private data or update existed record in database
func (r *ClientRepository) Save(ctx context.Context, req client.SaveRequest) error {
	if req.ID > 0 {
		return r.update(ctx, req)
	}
	return r.add(ctx, req)
}

//Get - get private data from database by id and type of private data
func (r *ClientRepository) Get(ctx context.Context, req client.GetRequest) ([]client.GetResponse, error) {
	txt := `SELECT t1.id
	,t1.external_id
	,t1.dek
	,t1.data
	,t1.view
	,IFNULL(t2.key, '')
	,IFNULL(t2.value, '')
	FROM private_data AS t1 
		LEFT JOIN metadata AS t2 ON t1.id = t2.owner_id 
	WHERE t1.type = ? and t1.id = ?`

	rows, err := r.db.QueryContext(ctx, txt, req.Type, req.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readEntityFromRows(rows)

}

//All - get all private data from database by type of private data
func (r *ClientRepository) All(ctx context.Context, tData model.Type) ([]client.GetResponse, error) {
	txt := `
		SELECT t1.id
		,t1.external_id
		,t1.dek
		,t1.data
		,t1.view
		,IFNULL(t2.key, '') 
		,IFNULL(t2.value, '')
		FROM private_data AS t1 
			LEFT JOIN metadata AS t2 ON t1.id = t2.owner_id 
		WHERE t1.type = ?			
	`

	rows, err := r.db.QueryContext(ctx, txt, tData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readEntityFromRows(rows)
}

//Delete - delete private data from database by id and type of private data
func (r *ClientRepository) Delete(ctx context.Context, req client.DeleteRequest) error {
	if err := r.deleteMetadataByOwner(ctx, req.ID); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx,
		`DELETE FROM private_data WHERE type = ? AND id = ?`,
		req.Type, req.ID)
	return err
}


//Rewrite - clear all data in database and add new data
func (r *ClientRepository) Rewrite(ctx context.Context, req []client.SaveRequest) error {
	txt := `
		DELETE FROM metadata;
		DELETE FROM sqlite_sequence WHERE name = 'metadata';
		DELETE FROM private_data;
		DELETE FROM sqlite_sequence WHERE name = 'private_data';
	`
	if _, err := r.db.ExecContext(ctx, txt); err != nil {
		return err
	}

	txtEntity := `
		INSERT INTO private_data (external_id, type, dek, data, view) 
		VALUES (?, ?, ?, ?, ?) 
		RETURNING id
	`
	txtMetadata := `
		INSERT INTO metadata (owner_id, key, value) 
		VALUES (?, ?, ?) 
	`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, e := range req {
		r, err := tx.ExecContext(ctx, txtEntity, e.ExternalID, e.Type, e.DEK, e.Data, e.View)
		if err != nil {
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		for _, m := range e.Metadata {
			if _, err := tx.ExecContext(ctx, txtMetadata, id, m.Key, m.Value); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (r *ClientRepository) add(ctx context.Context, req client.SaveRequest) error {

	txt := `
	INSERT INTO private_data (external_id, type, dek, data, view) 
	VALUES (?, ?, ?, ?, ?) 
	RETURNING id`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := r.db.ExecContext(ctx, txt, req.ExternalID, req.Type, req.DEK, req.Data, req.View)
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

func (r *ClientRepository) update(ctx context.Context, req client.SaveRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txt := `UPDATE private_data SET dek = ?, data = ?, view = ?, type = ? WHERE id = ?`

	if _, err := tx.ExecContext(ctx, txt, req.DEK, req.Data, req.View, req.Type, req.ID); err != nil {
		return err
	}

	if err := r.addOrUpdateMetadataByOwner(ctx, tx, req.ID, req.Metadata); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ClientRepository) deleteMetadataByOwner(ctx context.Context, ownerId int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM metadata WHERE owner_id = ?`,
		ownerId)
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

	if err = tx.QueryRowContext(ctx, q, ownerId, key).Scan(&id); err == sql.ErrNoRows {
		err = nil
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

func readEntityFromRows(rows *sql.Rows) ([]client.GetResponse, error) {
	tmp := make(map[int]client.GetResponse)

	for rows.Next() {
		var response client.GetResponse
		var key, value string
		if err := rows.Scan(
			&response.ID, &response.ExternalID, &response.DEK,
			&response.Data, &response.View, &key, &value); err != nil {

			return nil, err
		}

		if r, ok := tmp[response.ID]; !ok {
			response.Metadata = make([]client.MetadataValue, 0)

			if key != "" && value != "" {
				response.Metadata = append(response.Metadata, client.MetadataValue{Key: key, Value: value})
			}

			tmp[response.ID] = response
		} else {
			r = tmp[response.ID]
			r.Metadata = append(r.Metadata, client.MetadataValue{Key: key, Value: value})
			tmp[r.ID] = r
		}
	}

	result := make([]client.GetResponse, 0)
	for _, r := range tmp {
		result = append(result, r)
	}

	return result, nil
}
