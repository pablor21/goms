package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"text/template"

	"github.com/pablor21/goms/app/data/database/migrations"
	"github.com/pablor21/goms/pkg/database/config"
	"github.com/pablor21/goms/pkg/database/migrator"
	"github.com/pablor21/goms/pkg/logger"
)

type CreateMigrationRequest struct {
	Name string `json:"name"`
}

var migrationTemplateText = `package migrations

import (
	"database/sql"
	"context"
)

type {{.Name}} struct {}

func (m *{{.Name}}) Name() string {
	return "{{.Name}}"
}

func (m *{{.Name}}) Up(ctx context.Context, db *sql.DB) error {	
	sqlScript := ` + "`CREATE TABLE IF NOT EXISTS {{.Name}} (id SERIAL PRIMARY KEY NOT NULL)`" + `
	_,err:= db.ExecContext(ctx, sqlScript)
	return err
}

func (m *{{.Name}}) Down(ctx context.Context, db *sql.DB) error {
	// Write your rollback here
	_,err:= db.ExecContext(ctx, "DROP TABLE IF EXISTS ` + "`{{.Name}}`" + `")
	return err
}
`

var migrationsTableName = "system_migrations"

var migrationsTableCreateSqlite = `CREATE TABLE IF NOT EXISTS ` + migrationsTableName + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
var migrationsTableCreatePostgres = `CREATE TABLE IF NOT EXISTS ` + migrationsTableName + ` (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
var migrationsTableCreateMysql = `CREATE TABLE IF NOT EXISTS ` + migrationsTableName + ` (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

var ErrMigrationExists = os.ErrExist
var migrationTemplate = template.Must(template.New("migration").Parse(migrationTemplateText))

type DefaultMigrator struct {
	db         *sql.DB
	config     config.DatabaseConnectionConfig
	migrations []migrator.Migration
}

func NewDefaultMigrator(config config.DatabaseConnectionConfig, db *sql.DB) migrator.Migrator {
	return &DefaultMigrator{
		db:     db,
		config: config,
	}

}

func (m *DefaultMigrator) getRegisteredMigrations() []migrator.Migration {
	if m.migrations == nil {
		m.migrations = migrations.GetMigrations()
	}
	return m.migrations
}

func (m *DefaultMigrator) EnsureMigrationsTable(ctx context.Context) error {
	var migrationsTableCreateSQL string
	switch m.config.Driver {
	case "sqlite":
	case "sqlite3":
		migrationsTableCreateSQL = migrationsTableCreateSqlite
	case "postgres":
		migrationsTableCreateSQL = migrationsTableCreatePostgres
	case "mysql":
		migrationsTableCreateSQL = migrationsTableCreateMysql
	default:
		return ErrUnsupportedDriver
	}
	_, err := m.db.ExecContext(ctx, migrationsTableCreateSQL)
	return err
}

func (m *DefaultMigrator) Up(ctx context.Context, steps int) (migrations []migrator.MigrationRecord, err error) {
	toApply, err := m.Pending(ctx)
	if err != nil {
		return
	}

	if steps > 0 && steps < len(toApply) {
		toApply = toApply[:steps]
	}

	if len(toApply) == 0 {
		logger.Warn().Msg("No pending migrations")
		return
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	for _, migration := range toApply {
		logger.Info().Msg("Applying migration " + migration.Name())
		err = migration.Up(ctx, m.db)
		if err != nil {
			return
		}

		_, err = m.db.ExecContext(ctx, "INSERT INTO "+migrationsTableName+" (name) VALUES (?)", migration.Name())
		if err != nil {
			return
		}

		migrations = append(migrations, migrator.MigrationRecord{Name: migration.Name()})
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (m *DefaultMigrator) Down(ctx context.Context, steps int) (migrations []migrator.MigrationRecord, err error) {
	applied, err := m.Applied(ctx)
	if err != nil {
		return
	}

	if steps > 0 && steps < len(applied) {
		applied = applied[len(applied)-steps:]
	}

	if len(applied) == 0 {
		logger.Warn().Msg("No applied migrations")
		return
	}

	var tx *sql.Tx
	tx, err = m.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	for _, record := range applied {
		for _, migration := range m.getRegisteredMigrations() {
			if migration.Name() == record.Name {
				err = migration.Down(ctx, m.db)
				if err != nil {
					return
				}

				_, err = m.db.ExecContext(ctx, "DELETE FROM "+migrationsTableName+" WHERE name = ?", record.Name)
				if err != nil {
					return
				}

				migrations = append(migrations, record)
				break
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (m *DefaultMigrator) Reset(ctx context.Context) (err error) {
	err = m.EnsureMigrationsTable(ctx)
	if err != nil {
		return err
	}

	// down all applied migrations
	_, err = m.Down(ctx, 0)
	if err != nil {
		return err
	}

	// up all pending migrations
	_, err = m.Up(ctx, 0)
	if err != nil {
		return err
	}

	return
}

func (m *DefaultMigrator) Create(ctx context.Context, name string) error {
	// check if the migration already exists
	os.MkdirAll(m.config.MigrationsPath, os.ModePerm)

	fmt.Println(m.config.MigrationsPath + "/" + name + ".go")

	// if the migration already exists, return an error
	_, err := os.Stat(m.config.MigrationsPath + "/" + name + ".go")
	if err == nil {
		return ErrMigrationExists
	}

	file, err := os.Create(m.config.MigrationsPath + "/" + name + ".go")
	if err != nil {
		return err
	}
	defer file.Close()

	err = migrationTemplate.Execute(file, CreateMigrationRequest{Name: name})
	if err != nil {
		return err
	}

	return nil
}

func (m *DefaultMigrator) Pending(ctx context.Context) ([]migrator.Migration, error) {
	migrations := m.getRegisteredMigrations()
	applied, err := m.Applied(ctx)
	if err != nil {
		return nil, err
	}

	var pending []migrator.Migration
	for _, migration := range migrations {
		found := false
		for _, record := range applied {
			if record.Name == migration.Name() {
				found = true
				break
			}
		}
		if !found {
			pending = append(pending, migration)
		}
	}
	return pending, nil
}

func (m *DefaultMigrator) Applied(ctx context.Context) (ret []migrator.MigrationRecord, err error) {
	err = m.EnsureMigrationsTable(ctx)
	if err != nil {
		return
	}

	rows, err := m.db.QueryContext(ctx, "SELECT id, name, run_at FROM "+migrationsTableName+" ORDER BY run_at")
	if err != nil {
		return
	}

	for rows.Next() {
		var record migrator.MigrationRecord
		err = rows.Scan(&record.ID, &record.Name, &record.RunAt)
		if err != nil {
			return
		}
		ret = append(ret, record)
	}

	return
}
