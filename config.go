package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/HomayoonAlimohammadi/reviewer/pkg/modelvendor"
)

type Config struct {
	ModelVendor modelvendor.Config `json:"model_vendor" yaml:"model_vendor" mapstructure:"model_vendor"`
}

func loadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	// default values
	viper.SetDefault("model_vendor.model_name", "llama3.2:3b")
	viper.SetDefault("model_vendor.vendor_name", modelvendor.Ollama)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	// var config map[string]any
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
