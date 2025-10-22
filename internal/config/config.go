package config

import (
	"github.com/mwinyimoha/commons/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`
	AppID          string `mapstructure:"APP_ID"`
	Debug          bool   `mapstructure:"DEBUG"`
	ServerPort     int    `mapstructure:"SERVER_PORT"`
	DefaultTimeout int    `mapstructure:"DEFAULT_TIMEOUT"`
}

func New() (*Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigType("env")

	viper.SetDefault("SERVICE_NAME", "")
	viper.SetDefault("SERVICE_VERSION", "0.1.0")
	viper.SetDefault("APP_ID", "")
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("DEFAULT_TIMEOUT", 10)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.WrapError(err, errors.Internal, "failed to load configuration variables")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.WrapError(err, errors.Internal, "failed to unmarshal config")
	}

	return &cfg, nil
}
