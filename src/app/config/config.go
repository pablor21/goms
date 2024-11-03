package config

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	server_config "github.com/pablor21/goms/app/server/config"
	database_config "github.com/pablor21/goms/pkg/database/config"
	"github.com/pablor21/goms/pkg/logger"
	mailer_config "github.com/pablor21/goms/pkg/mailer/config"
	security_config "github.com/pablor21/goms/pkg/security/config"
	storage_config "github.com/pablor21/goms/pkg/storage/config"
	"github.com/spf13/viper"
)

//go:embed default.yml
var defaultConfig embed.FS

type AppConfig struct {
	Name        string `json:"name" yaml:"name" MAPSTRUCTURE:"name"`
	Version     string `json:"version" yaml:"version" MAPSTRUCTURE:"version"`
	Environment string `json:"environment" yaml:"environment" MAPSTRUCTURE:"environment"`
	Description string `json:"description" yaml:"description" MAPSTRUCTURE:"description"`
}

type Config struct {
	App      AppConfig                      `json:"app" yaml:"app" MAPSTRUCTURE:"app"`
	Logger   logger.LoggerConfig            `json:"logger" yaml:"logger" MAPSTRUCTURE:"logger"`
	Server   server_config.ServerConfig     `json:"server" yaml:"server" MAPSTRUCTURE:"server"`
	Database database_config.DatabaseConfig `json:"database" yaml:"database" MAPSTRUCTURE:"database"`
	Storage  storage_config.StorageConfig   `json:"storage" yaml:"storage" MAPSTRUCTURE:"storage"`
	Assets   AssetConfig                    `json:"assets" yaml:"assets" MAPSTRUCTURE:"assets"`
	OTP      OTPConfig                      `json:"otp" yaml:"otp" MAPSTRUCTURE:"otp"`
	Mailer   mailer_config.MailerConfig     `json:"mailer" yaml:"mailer" MAPSTRUCTURE:"mailer"`
	Security security_config.SecurityConfig `json:"security" yaml:"security" MAPSTRUCTURE:"security"`
	Viper    *viper.Viper
	Values   map[string]interface{}
}

// type AssetConfig struct {
// 	BasePath string `json:"base_path" yaml:"base_path" MAPSTRUCTURE:"base_path" default:"/assets"`
// 	Image    struct {
// 		DefaultQuality int `mapstructure:"DEFAULT_QUALITY" default:"80"`
// 	} `mapstructure:"IMAGE"`
// }

var config *Config

func GetConfig() *Config {
	return config
}

func InitConfig(files []string) {
	cfg := &Config{}
	if err := defaults.Set(cfg); err != nil {
		panic(fmt.Errorf("unable to set defaults: %v", err))
	}

	cfg.App.Environment = os.Getenv("ENV")
	if cfg.App.Environment == "" {
		cfg.App.Environment = os.Getenv("ENVIRONMENT")
		if cfg.App.Environment == "" {
			cfg.App.Environment = "production"
		}
	}
	log.Default().Printf("Loading configuration")
	v := viper.New()
	cfg.Viper = v
	// v.SetEnvPrefix(viper.GetString("ENV"))
	v.AutomaticEnv()

	defaultConfigContent, err := defaultConfig.Open("default.yml")
	if err == nil {
		v.SetConfigType("yaml")
		v.MergeConfig(defaultConfigContent)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	configPath := filepath.Join(exPath, "config")

	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	v.AddConfigPath("")
	v.SetConfigType("yaml")

	for _, f := range files {
		v.SetConfigFile(f)
		err = v.MergeInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); (!ok) && err.Error() != "open "+f+": no such file or directory" {
				log.Default().Fatalf("Cannot read cofiguration: %s", err)
			}
		}
	}

	v.SetConfigName("")
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); (!ok) && err.Error() != "open .env: no such file or directory" {
			log.Default().Fatalf("Cannot read cofiguration: %s", err)
		}
	}
	err = v.Unmarshal(&cfg)
	if err != nil {
		log.Default().Fatalf("environment can't be loaded: %s", err)
	}

	// add the config.[env].yml file
	env := cfg.App.Environment
	log.Default().Printf("Loading environment specific configuration for: %s", env)
	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); (!ok) && err.Error() != "open config."+env+".yml: no such file or directory" {
			log.Default().Fatalf("Cannot read cofiguration: %s", err)
		}
	}

	v.SetConfigName("")
	v.SetConfigFile(".env." + env)
	v.SetConfigType("env")
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); (!ok) && err.Error() != "open .env."+env+": no such file or directory" {
			log.Default().Fatalf("Cannot read cofiguration: %s", err)
		}
	}

	err = v.Unmarshal(&cfg)

	if err != nil {
		log.Default().Fatalf("environment can't be loaded: %s", err)
	}
	cfg.Values = v.AllSettings()
	config = cfg
}
