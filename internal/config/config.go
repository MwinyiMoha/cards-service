package config

import (
	"cards-service/internal/core/ports"

	"github.com/go-playground/validator/v10"
	"github.com/mwinyimoha/commons/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME" validate:"required"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION" validate:"required"`
	AppID          string `mapstructure:"APP_ID"`
	Debug          bool   `mapstructure:"DEBUG"`
	ServerPort     int    `mapstructure:"SERVER_PORT" validate:"required,min=1,max=65535"`
	DefaultTimeout int    `mapstructure:"DEFAULT_TIMEOUT" validate:"required,min=1"`
}

func New(val ports.AppValidator) (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")

	v.SetDefault("SERVICE_NAME", "")
	v.SetDefault("SERVICE_VERSION", "0.1.0")
	v.SetDefault("APP_ID", "")
	v.SetDefault("DEBUG", true)
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("DEFAULT_TIMEOUT", 10)

	v.AutomaticEnv()

	v.AddConfigPath("./")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.WrapError(err, errors.Internal, "failed to load configuration file")
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.WrapError(err, errors.Internal, "failed to unmarshal config")
	}

	if err := cfg.validate(val); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate(v ports.AppValidator) error {
	if err := v.Struct(c); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			violations := errors.BuildViolations(verr)
			return errors.NewValidationError(violations)
		}

		return errors.WrapError(err, errors.Internal, "config validation failed")
	}

	return nil
}
