package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Environment          string        `mapstructure:"ENV"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationsURL        string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // Can be made configurable if needed
	viper.SetConfigType("env") // Supports JSON, YAML, etc.

	viper.AutomaticEnv() // Automatically override with environment variables

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config

	// Unmarshal config values into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
