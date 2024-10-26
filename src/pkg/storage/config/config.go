package config

type StorageBucketConfig struct {
	Uri string `mapstructure:"URI"`
}

type StorageConfig map[string]StorageBucketConfig
