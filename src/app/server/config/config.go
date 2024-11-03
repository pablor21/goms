package config

type CookieConfig struct {
	Name     string `mapstructure:"NAME" default:"session"`
	MaxAge   int    `mapstructure:"MAX_AGE" default:"2592000"`
	Secure   bool   `mapstructure:"SECURE" default:"true"`
	HttpOnly bool   `mapstructure:"HTTP_ONLY" default:"true"`
	SameSite string `mapstructure:"SAME_SITE" default:"strict"`
	Path     string `mapstructure:"PATH" default:"/"`
	Domain   string `mapstructure:"DOMAIN"`
}

type SessionStoreConfig struct {
	Type   string                 `mapstructure:"TYPE" default:"cookie"`
	Config map[string]interface{} `mapstructure:"CONFIG"`
}

type SessionConfig struct {
	Secret   string             `mapstructure:"SECRET"`
	Lifetime int                `mapstructure:"LIFETIME" default:"86400"`
	Store    SessionStoreConfig `mapstructure:"STORE"`
	Cookie   CookieConfig
}

type AppConfig struct {
	BasePath string `mapstructure:"BASE_PATH" default:"/api"`
	Domain   string `mapstructure:"DOMAIN" default:""`
	Key      string `mapstructure:"KEY"`
}

type ApiConfig struct {
	AppConfig
}

type FrontendConfig struct {
	AppConfig
}

type AdminConfig struct {
	AppConfig
}

type AppsConfig struct {
	Api      ApiConfig      `mapstructure:"API"`
	Frontend FrontendConfig `mapstructure:"FRONTEND"`
	Admin    AdminConfig    `mapstructure:"ADMIN"`
}

type ServerConfig struct {
	Host string                   `json:"host" yaml:"host" MAPSTRUCTURE:"host"`
	Port int                      `json:"port" yaml:"port" MAPSTRUCTURE:"port"`
	Apps AppsConfig               `json:"apps" yaml:"apps" mapstructure:"APPS"`
	Auth map[string]SessionConfig `json:"auth" yaml:"auth" mapstructure:"AUTH"`
}
