package config

type OTPConfig struct {
	Length int `json:"length" mapstructure:"LENGHT" default:"6"`
	// Duration in seconds
	Lifetime    int `json:"lifetime" mapstructure:"LIFETIME" default:"1800"` // 30 minutes
	MaxAttempts int `json:"maxAttempts" mapStructure:"MAX_ATTEMPTS" default:"3"`
	ResendDelay int `json:"resendDelay" mapstructure:"RESEND_DELAY" default:"60"` // 1 minute
}
