package config

type AssetSectionConfig struct {
	DefaultQuality int      `json:"default_quality" yaml:"default_quality" MAPSTRUCTURE:"default_quality" default:"80"`
	ValidSizes     []string `json:"valid_sizes" yaml:"valid_sizes" MAPSTRUCTURE:"valid_sizes"`
}

type AssetImageConfig struct {
	DefaultQuality int                           `json:"default_quality" yaml:"default_quality" MAPSTRUCTURE:"default_quality" default:"80"`
	ValidQualities []int                         `json:"valid_qualities" yaml:"valid_qualities" MAPSTRUCTURE:"valid_qualities"`
	ValidSizes     []string                      `json:"valid_sizes" yaml:"valid_sizes" MAPSTRUCTURE:"valid_sizes"`
	ValidMimeTypes []string                      `json:"valid_mime_types" yaml:"valid_mime_types" MAPSTRUCTURE:"valid_mime_types"`
	Sections       map[string]AssetSectionConfig `json:"sections" yaml:"sections" MAPSTRUCTURE:"sections"`
}

type AssetConfig struct {
	BasePath string           `json:"base_path" yaml:"base_path" MAPSTRUCTURE:"base_path" default:"/assets"`
	Image    AssetImageConfig `json:"image" yaml:"image" MAPSTRUCTURE:"image"`
}
