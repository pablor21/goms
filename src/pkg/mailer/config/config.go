package config

type MailerConnectionConfig struct {
	Driver       string `json:"driver" mapstructure:"DRIVER" default:"smtp"`
	Host         string `json:"host" mapstructure:"HOST" default:"localhost"`
	Port         int    `json:"port" mapstructure:"PORT" default:"587"`
	Username     string `json:"username" mapstructure:"USERNAME" default:""`
	Password     string `json:"password" mapstructure:"PASSWORD" default:""`
	From         string `json:"from" mapstructure:"FROM" default:""`
	SmtpAuthType string `json:"auth_type" mapstructure:"AUTH_TYPE" default:"LOGIN"`
}

type MailerConfig map[string]MailerConnectionConfig
