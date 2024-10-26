package migrations

import (
	"context"
	"database/sql"

	"github.com/pablor21/goms/app/data/database/migrations/scripts"
)

type CreateAuthTables struct{}

func (m *CreateAuthTables) Name() string {
	return "CreateAuthTables"
}

func (m *CreateAuthTables) Up(ctx context.Context, db *sql.DB) error {
	sqlScript, err := scripts.Scripts.ReadFile("20241024_2220_crate_auth_tables.up.sql")
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, string(sqlScript))
	if err != nil {
		return err
	}

	// create a default user
	_, err = db.ExecContext(ctx, `INSERT INTO users (first_name, last_name, role, status, email, super_admin, password) VALUES('Admin', 'User', 'admin', 'active', 'admin@local.dev', 1, '$2a$14$y4dQdjXozS0qlbHceBl75uIN65JOVHrCwIiuH4PlqsKtQFZLzb3Ga');`)

	return err
}

func (m *CreateAuthTables) Down(ctx context.Context, db *sql.DB) error {
	sqlScript, err := scripts.Scripts.ReadFile("20241024_2220_crate_auth_tables.down.sql")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(sqlScript))
	return err
}
