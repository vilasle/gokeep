package postgres

import (
	"context"
	"database/sql"

	"github.com/vilasle/gokeep/internal/model"
)

type EntityType int

const (
	EntityTypeUsepass EntityType = iota + 1
	EntityTypeBankCard
	EntityTypeTextData
	EntityTypeBinaryData
)

type PrivateDataRepository struct {
	db *sql.DB
}

func NewPrivateDataRepository(db *sql.DB) (*PrivateDataRepository, error) {
	repository := &PrivateDataRepository{db: db}
	return repository, repository.initSchema()
}

func (r *PrivateDataRepository) Add(ctx context.Context, data model.PrivateDataSave) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := r.addEntity(ctx, tx, data)
	if err != nil {
		return 0, err
	}

	err = r.addEncryptedData(ctx, tx, id, data.Data, data.DEK)
	if err != nil {
		return 0, err
	}

	if err := r.saveMetadata(ctx, tx, data.ID, data.Metadata); err != nil {
		return 0, err
	}

	err = tx.Commit()

	return id, err
}

func (r *PrivateDataRepository) addEntity(ctx context.Context, tx *sql.Tx, data model.PrivateDataSave) (int, error) {
	txt := `INSERT INTO entity (user_id, "type", view) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := tx.
		QueryRowContext(ctx, txt, data.UserID, data.Type, data.View).
		Scan(&id)

	return id, err
}

func (r *PrivateDataRepository) addEncryptedData(ctx context.Context, tx *sql.Tx, id int, data, dek []byte) error {
	txt := `INSERT INTO data_encrypted (entity_id, data, dek) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, txt, id, data, dek)
	return err
}

func (r *PrivateDataRepository) Update(ctx context.Context, data model.PrivateDataSave) (int, error) {
	if data.ID == 0 {
		return 0, model.ErrEmptyID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err := r.updateEntity(ctx, tx, data); err != nil {
		return 0, err
	}

	if err := r.updateEncryptedData(ctx, tx, data.ID, data.Data, data.DEK); err != nil {
		return 0, err
	}

	if err := r.saveMetadata(ctx, tx, data.ID, data.Metadata); err != nil {
		return 0, err
	}

	err = tx.Commit()

	return data.ID, err
}

func (r *PrivateDataRepository) updateEntity(ctx context.Context, tx *sql.Tx, data model.PrivateDataSave) error {
	txt := `UPDATE entity SET user_id = $1, "type" = $2, view = $3, updated_at = NOW() WHERE id = $4`
	_, err := tx.ExecContext(ctx, txt, data.UserID, data.Type, data.View, data.ID)
	return err
}

func (r *PrivateDataRepository) updateEncryptedData(ctx context.Context, tx *sql.Tx, id int, data, dek []byte) error {
	txt := `UPDATE data_encrypted SET data = $1, dek = $2 WHERE entity_id = $3`
	_, err := tx.ExecContext(ctx, txt, data, dek, id)
	return err
}

func (r *PrivateDataRepository) saveMetadata(ctx context.Context, tx *sql.Tx, id int, metadata map[string]string) error {
	txt := `
	INSERT INTO metadata (entity_id , key, value) 
	VALUES ($1, $2, $3) w
	ON CONFLICT (entity_id,key) DO 
	UPDATE SET value = EXCLUDED.value;
	`
	for k, v := range metadata {
		if _, err := tx.ExecContext(ctx, txt, id, k, v); err != nil {
			return err
		}
	}
	return nil
}

func (r *PrivateDataRepository) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return model.ErrEmptyID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txt := []string{
		`DELETE FROM data_encrypted WHERE entity_id = $1`,
		`DELETE FROM entity WHERE id = $1`,
	}

	for _, txt := range txt {
		if _, err := r.db.ExecContext(ctx, txt, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PrivateDataRepository) Get(ctx context.Context, id int) (response model.PrivateDataInfo, err error) {
	txt := `
	SELECT t1.id 
		,t1.view
		,t2.data
		,t2.dek 
		,t3.key
		,t3.value
	FROM entity AS t1 
		LEFT JOIN data_encrypted AS t2 ON t1.id = t2.entity_id
		LEFT JOIN metadata AS t3 ON t1.id = t3.entity_id
	WHERE t1.id = $1`

	rows, err := r.db.QueryContext(ctx, txt, id)
	if err != nil {
		return
	}
	defer rows.Close()

	result, err := readGettingRows(rows)
	for k := range result {
		response = result[k]
		break
	}

	return response, err
}

func (r *PrivateDataRepository) List(ctx context.Context,
	modelType model.Type, owner *model.User) (response []model.PrivateDataInfo, err error) {

	txt := `
	SELECT t1.id 
		,t1.view
		,t2.data
		,t2.dek 
		,t3.key
		,t3.value
	FROM entity AS t1 
		LEFT JOIN data_encrypted AS t2 ON t1.id = t2.entity_id
		LEFT JOIN metadata AS t3 ON t1.id = t3.entity_id
	WHERE t1.user_id = $1 AND t1.type = $2`

	rows, err := r.db.QueryContext(ctx, txt, owner.ID, modelType)
	if err != nil {
		return
	}
	defer rows.Close()

	result, err := readGettingRows(rows)
	response = make([]model.PrivateDataInfo, 0, len(result))
	for k := range result {
		response = append(response, result[k])
	}

	return response, err
}

func readGettingRows(rows *sql.Rows) (result map[int]model.PrivateDataInfo, err error) {
	result = make(map[int]model.PrivateDataInfo)
	for rows.Next() {
		var d model.PrivateDataInfo
		var key, value string
		err := rows.Scan(&d.ID, &d.View, &d.Data, &d.DEK, &key, &value)
		if err != nil {
			return result, err
		}
		if _, ok := result[d.ID]; !ok {
			d.Metadata = make(map[string]string)
			if key != "" {
				d.Metadata[key] = value
			}
			result[d.ID] = d
		} else {
			if key != "" {
				d.Metadata[key] = value
			}
		}
	}
	return result, err
}

func (r *PrivateDataRepository) initSchema() error {
	txt := `
	CREATE TABLE IF NOT EXISTS entity (
		id bigint GENERATED ALWAYS AS IDENTITY UNIQUE,
		user_id bigint NOT NULL,
		"type" int NOT NULL,	
		view TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	CREATE INDEX IF NOT EXISTS entity_user_id_idx ON entity (user_id);
	CREATE INDEX IF NOT EXISTS entity_type_idx ON entity ("type");
	
	CREATE TABLE IF NOT EXISTS data_encrypted (
		id bigint GENERATED ALWAYS AS IDENTITY UNIQUE,
		entity_id bigint NOT NULL,
		data BYTEA NOT NULL,
		dek BYTEA NOT NULL,
		FOREIGN KEY (entity_id) REFERENCES entity (id)
	);
	CREATE INDEX IF NOT EXISTS data_encrypted_entity_id_idx ON data_encrypted (entity_id);

	CREATE TABLE IF NOT EXISTS metadata (
		id bigint GENERATED ALWAYS AS IDENTITY UNIQUE,
		entity_id bigint NOT NULL,
		key TEXT,
		value TEXT,
		FOREIGN KEY (entity_id) REFERENCES entity (id),
		CONSTRAINT metadata_entity_id_key_unique UNIQUE (entity_id, key)
	);
	CREATE INDEX IF NOT EXISTS metadata_entity_id_idx ON metadata (entity_id);
	CREATE INDEX IF NOT EXISTS metadata_key_idx ON metadata (key);
	`

	_, err := r.db.Exec(txt)
	return err
}
