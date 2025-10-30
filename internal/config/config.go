package config

import (
	"os"
	"strconv"

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
	v := viper.New()
	v.SetConfigType("env")

	v.SetDefault("SERVICE_NAME", "")
	v.SetDefault("SERVICE_VERSION", "0.1.0")
	v.SetDefault("APP_ID", "")
	v.SetDefault("DEBUG", true)
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("DEFAULT_TIMEOUT", 10)

	v.AutomaticEnv()

	debug := true
	if raw := os.Getenv("DEBUG"); raw != "" {
		val, err := strconv.ParseBool(raw)
		if err == nil {
			debug = val
		}
	}

	if debug {
		configPath := "./"
		v.AddConfigPath(configPath)

		if err := v.ReadInConfig(); err != nil {
			return nil, errors.WrapError(err, errors.Internal, "failed to load configuration file (DEBUG=true)")
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.WrapError(err, errors.Internal, "failed to unmarshal config")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.ServiceName == "" {
		return errors.NewErrorf(errors.Internal, "SERVICE_NAME must be set")
	}

	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return errors.NewErrorf(errors.Internal, "invalid SERVER_PORT value")
	}

	if c.DefaultTimeout <= 0 {
		return errors.NewErrorf(errors.Internal, "DEFAULT_TIMEOUT must be > 0")
	}

	return nil
}
