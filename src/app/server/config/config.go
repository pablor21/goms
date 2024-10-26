package config

type ServerConfig struct {
	Host string `json:"host" yaml:"host" MAPSTRUCTURE:"host"`
	Port int    `json:"port" yaml:"port" MAPSTRUCTURE:"port"`
}
