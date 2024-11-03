package config

type OtpConfig struct {
	RetryPeriod     int `json:"retryPeriod" mapstructure:"RETRY_PERIOD" default:"15"`
	Period          int `json:"period" mapstructure:"PERIOD" default:"15"`
	Digits          int `json:"digits" mapstructure:"DIGITS" default:"6"`
	RetryMultiplier int `json:"retryMultiplier" mapstructure:"RETRY_MULTIPLIER" default:"1"`
}

type EncryptionConfig struct {
	Algorithm string `json:"algorithm" mapstructure:"ALGORITHM" default:"aes"`
	Key       string `json:"key" mapstructure:"KEY" default:""`
}

type SecurityConfig struct {
	Otp        OtpConfig        `json:"otp" mapstructure:"OTP"`
	Encryption EncryptionConfig `json:"encryption" mapstructure:"ENCRYPTION"`
}
