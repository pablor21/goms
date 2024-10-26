package migrator

import (
	"context"
	"database/sql"
	"time"
)

type MigrationRecord struct {
	ID    int       `db:"id"`
	Name  string    `db:"name"`
	RunAt time.Time `db:"run_at"`
}

type Migration interface {
	Name() string
	Up(ctx context.Context, db *sql.DB) error
	Down(ctx context.Context, db *sql.DB) error
}

type Migrator interface {
	Up(ctx context.Context, steps int) ([]MigrationRecord, error)
	Down(ctx context.Context, steps int) ([]MigrationRecord, error)
	Reset(ctx context.Context) error
	Pending(ctx context.Context) ([]Migration, error)
	Applied(ctx context.Context) ([]MigrationRecord, error)
	Create(ctx context.Context, name string) error
	// RegisterMigration(m Migration)
}

// MigrationRegistry is a map of migrations (key is the database connection name)
type MigrationRegistry map[string][]Migration
