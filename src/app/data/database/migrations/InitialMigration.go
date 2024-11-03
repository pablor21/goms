package migrations

import (
	"context"
	"database/sql"

	"github.com/pablor21/goms/app/data/database/migrations/scripts"
)

type InitialMigration struct{}

func (m *InitialMigration) Name() string {
	return "CreateAuthTables"
}

func (m *InitialMigration) Up(ctx context.Context, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlScript, err := scripts.Scripts.ReadFile("initial_migration.up.sql")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, string(sqlScript))
	if err != nil {
		return err
	}

	// create a default user
	_, err = tx.ExecContext(ctx, `INSERT INTO users (first_name, last_name, role, status, email, super_admin, password) VALUES('Admin', 'User', 'admin', 'active', 'admin@local.dev', 1, '$2a$14$y4dQdjXozS0qlbHceBl75uIN65JOVHrCwIiuH4PlqsKtQFZLzb3Ga');`)
	if err != nil {
		return err
	}
	// create a media library
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_libraries (name, description) VALUES('Default Library', 'Default Media Library');`)
	if err != nil {
		return err
	}
	// create a media folder
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_folders (library_id, name, path) VALUES(1, '.', '/');`)
	if err != nil {
		return err
	}
	// set the root folder
	_, err = tx.ExecContext(ctx, `UPDATE asset_libraries SET root_folder_id = 1 WHERE id = 1;`)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (m *InitialMigration) Down(ctx context.Context, db *sql.DB) error {
	sqlScript, err := scripts.Scripts.ReadFile("initial_migration.down.sql")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(sqlScript))
	return err
}
