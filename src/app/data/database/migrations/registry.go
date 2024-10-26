package migrations

import "github.com/pablor21/goms/pkg/database/migrator"

// MigrationsRegistry is the list of migrations that the migrator will use
// add your migrations instances here
func GetMigrations() []migrator.Migration {
	return []migrator.Migration{
		&CreateAuthTables{},
	}
}
