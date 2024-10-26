package config

type DatabaseConnectionConfig struct {
	Name            string `json:"name" mapstructure:"NAME" default:""`
	Driver          string `json:"driver" mapstructure:"DRIVER" default:""`
	Type            string `json:"type" mapstructure:"TYPE" default:"gorm"`
	URI             string `json:"uri" mapstructure:"URI" default:""`
	MigrationsPath  string `json:"migrationsPath" mapstructure:"MIGRATIONS_PATH" default:"data/migrations"`
	MaxIdleConns    int    `json:"max_idle_conns" default:"64" mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns    int    `json:"max_open_conns" default:"64" mapstructure:"MAX_OPEN_CONNS"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" default:"60" mapstructure:"CONN_MAX_LIFETIME"`
}

type DatabaseConfig map[string]DatabaseConnectionConfig

// type DatabaseConfig struct {
// 	Connections map[string]DatabaseConnectionConfig `json:"connections" mapstructure:"CONNECTIONS"`
// }
