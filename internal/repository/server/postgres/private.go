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

func (r *PrivateDataRepository) Get(ctx context.Context, id int) (model.PrivateDataInfo, error) {
	txt := `SELECT id, user_id, "type", view FROM entity WHERE id = $1`
	var data model.PrivateDataInfo
	err := r.db.
		QueryRowContext(ctx, txt, id).
		Scan(&data.ID, &data.UserID, &data.Type, &data.View)

	return data, err
}

func (r *PrivateDataRepository) List(ctx context.Context, modelType model.Type, owner *model.User) ([]model.PrivateDataInfo, error) {
	txt := `SELECT id, user_id, "type", view FROM entity WHERE "type" = $1 AND user_id = $2`
	rows, err := r.db.QueryContext(ctx, txt, modelType, owner.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.PrivateDataInfo
	for rows.Next() {
		var d model.PrivateDataInfo
		err := rows.Scan(&d.ID, &d.UserID, &d.Type, &d.View)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
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
	`
	_, err := r.db.Exec(txt)
	return err
}
