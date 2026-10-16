package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func setTestEnvironment(t *testing.T) {
	t.Helper()

	values, err := godotenv.Read(filepath.Join("..", ".env.example"))
	require.NoError(t, err)

	for key, value := range values {
		t.Setenv(key, value)
	}

	t.Setenv("APP_NAME", "environment-app")
	t.Chdir(t.TempDir())
}

//nolint:paralleltest // Configuration tests change the process environment and working directory.
func TestConfigWithoutDotEnv(t *testing.T) {
	setTestEnvironment(t)

	cfg, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "environment-app", cfg.App().Name)
}

//nolint:paralleltest // Configuration tests change the process environment and working directory.
func TestEnvironmentOverridesDotEnv(t *testing.T) {
	setTestEnvironment(t)
	require.NoError(t, os.WriteFile(".env", []byte("APP_NAME=file-app\n"), 0o600))

	cfg, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "environment-app", cfg.App().Name)
}

//nolint:paralleltest // Configuration tests change the process environment and working directory.
func TestMalformedDotEnvReturnsError(t *testing.T) {
	setTestEnvironment(t)
	require.NoError(t, os.WriteFile(".env", []byte("APP_NAME=\"unterminated\n"), 0o600))

	_, err := NewConfig()
	require.Error(t, err)
	require.ErrorContains(t, err, "config error")
}
