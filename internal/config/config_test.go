package config

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/mwinyimoha/commons/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockValidator struct{}

func (m *mockValidator) Struct(s interface{}) error {
	return errors.NewErrorf(errors.Internal, "mock validation error")
}

func newValidator() *validator.Validate {
	return validator.New()
}

func resetEnv() {
	os.Unsetenv("SERVICE_NAME")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DEFAULT_TIMEOUT")
	os.Unsetenv("APP_ID")
	os.Unsetenv("DEBUG")
}

func TestNew(t *testing.T) {

	t.Run("Success With Defaults", func(t *testing.T) {
		defer resetEnv()
		v := newValidator()

		os.Setenv("SERVICE_NAME", "TestService")

		cfg, err := New(v)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, "TestService", cfg.ServiceName)
		assert.Equal(t, 8080, cfg.ServerPort)        // default
		assert.Equal(t, 10, cfg.DefaultTimeout)      // default
		assert.Equal(t, "0.1.0", cfg.ServiceVersion) // default
	})
}

func TestDefaultenvOverrides(t *testing.T) {
	defer resetEnv()
	v := newValidator()

	os.Setenv("SERVICE_NAME", "CustomService")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DEFAULT_TIMEOUT", "20")
	os.Setenv("DEBUG", "false")

	cfg, err := New(v)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "CustomService", cfg.ServiceName)
	assert.Equal(t, 9090, cfg.ServerPort)
	assert.Equal(t, 20, cfg.DefaultTimeout)
	assert.Equal(t, false, cfg.Debug)
}

func TestValidationFails(t *testing.T) {

	t.Run("Missing Required Field", func(t *testing.T) {
		defer resetEnv()
		v := newValidator()

		os.Setenv("SERVER_PORT", "8080")
		os.Setenv("DEBUG", "false")

		cfg, err := New(v)
		assert.Nil(t, cfg)
		require.Error(t, err)

		_, ok := err.(*errors.Error)
		assert.True(t, ok, "expected custom error type")
	})

	t.Run("Invalid Field Value", func(t *testing.T) {
		defer resetEnv()
		v := newValidator()

		os.Setenv("SERVICE_NAME", "TestService")
		os.Setenv("SERVER_PORT", "70000") // invalid port
		os.Setenv("DEBUG", "false")

		cfg, err := New(v)
		assert.Nil(t, cfg)
		require.Error(t, err)

		_, ok := err.(*errors.Error)
		assert.True(t, ok, "expected custom error type")
	})

	t.Run("Generic Error", func(t *testing.T) {
		defer resetEnv()

		os.Setenv("SERVICE_NAME", "TestService")
		os.Setenv("SERVER_PORT", "8080")
		os.Setenv("DEFAULT_TIMEOUT", "10")
		os.Setenv("DEBUG", "false")

		cfg, err := New(&mockValidator{})
		assert.Nil(t, cfg)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "config validation failed")
	})
}

func TestDebugMode(t *testing.T) {

	t.Run("No Config File", func(t *testing.T) {
		defer resetEnv()
		v := newValidator()

		os.Setenv("SERVICE_NAME", "TestService")
		os.Setenv("DEBUG", "true")

		cfg, err := New(v)
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("Invalid Config File", func(t *testing.T) {
		defer resetEnv()
		v := newValidator()

		os.Setenv("SERVICE_NAME", "TestService")
		os.Setenv("DEBUG", "true")

		tmp := "./config.env"
		os.WriteFile(tmp, []byte("INVALID :: = value"), 0644)
		defer os.Remove(tmp)

		cfg, err := New(v)
		assert.Nil(t, cfg)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "failed to load configuration file")
	})
}
