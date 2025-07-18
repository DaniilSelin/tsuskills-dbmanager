package config

import (
    "github.com/spf13/viper"
    "github.com/opensearch-project/opensearch-go"
    "log"
    "fmt"
    "os"
)

const ConfigFilePath = "config/config.local.yml"

type Config struct {
	Search opensearch.Config `mapstructure:"opensearch"`
}

func LoadConfig() (Config, error) {
	if cfgFile := os.Getenv("CONFIG_FILE_PATH"); cfgFile != "" {
	    viper.SetConfigFile(cfgFile)
	} else {
	    viper.SetConfigFile(ConfigFilePath)
		log.Println("failed to read CONFIG_FILE_PATH, using default path")
	}

	if err := viper.ReadInConfig(); err != nil {
	    return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return Config{}, fmt.Errorf("unable to decode into struct: %w", err)
    }

	return cfg, nil
}