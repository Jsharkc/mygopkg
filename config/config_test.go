package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Define a test config struct
type TestConfig struct {
	Name   string `mapstructure:"name"`
	Value  int    `mapstructure:"value"`
	Nested struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"nested"`
}

func TestInit(t *testing.T) {
	// Reset viper global instance before each test to avoid state pollution
	defer viper.Reset()

	t.Run("success with existing config file", func(t *testing.T) {
		viper.Reset()
		// Create a temporary config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test_config.yaml")
		configContent := `name: test
value: 42
nested:
  enabled: true`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		var cfg TestConfig
		_, err = Init("test_config", "yaml", []string{tempDir}, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, "test", cfg.Name)
		assert.Equal(t, 42, cfg.Value)
		assert.True(t, cfg.Nested.Enabled)
	})

	t.Run("success with non-existing config file", func(t *testing.T) {
		viper.Reset()
		tempDir := t.TempDir()
		var cfg TestConfig
		_, err := Init("nonexistent", "yaml", []string{tempDir}, &cfg)
		assert.NoError(t, err)
		// Should have zero values
		assert.Equal(t, "", cfg.Name)
		assert.Equal(t, 0, cfg.Value)
		assert.False(t, cfg.Nested.Enabled)
	})

	t.Run("error with invalid config file", func(t *testing.T) {
		viper.Reset()
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "invalid_config.yaml")
		invalidContent := `invalid: yaml: content: [`
		err := os.WriteFile(configFile, []byte(invalidContent), 0644)
		assert.NoError(t, err)

		var cfg TestConfig
		_, err = Init("invalid_config", "yaml", []string{tempDir}, &cfg)
		assert.Error(t, err)
		t.Log(err)
	})

	t.Run("error with nil rawVal", func(t *testing.T) {
		viper.Reset()
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "invalid_config.yaml")
		configContent := `name: test
value: 42
nested:
  enabled: true`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		_, err = Init("invalid_config", "yaml", []string{tempDir}, nil)
		assert.Error(t, err)
		t.Log(err)
	})
}

func TestIniWithEnv(t *testing.T) {
	t.Run("success with config file and env vars", func(t *testing.T) {
		// Set environment variables (note: viper replaces "." with "_" in env keys)
		os.Setenv("TEST_NAME", "from_env")
		os.Setenv("TEST_VALUE", "123")
		os.Setenv("TEST_NESTED_ENABLED", "true")
		defer func() {
			os.Unsetenv("TEST_NAME")
			os.Unsetenv("TEST_VALUE")
			os.Unsetenv("TEST_NESTED_ENABLED")
		}()

		// Create a temporary config file with partial config
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test_config.yaml")
		configContent := `name: from_file
nested:
  enabled: false`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		var cfg TestConfig
		_, err = IniWithEnv("test_config", "yaml", []string{tempDir}, &cfg, "test")

		assert.NoError(t, err)
		// Config file takes precedence over env vars, but env vars provide missing values
		assert.Equal(t, "from_env", cfg.Name) // config file takes precedence
		assert.Equal(t, 123, cfg.Value)       // from env (not in config)
		assert.True(t, cfg.Nested.Enabled)    // config file takes precedence
	})

	t.Run("success with only env vars", func(t *testing.T) {
		// Set environment variables (uppercase with prefix)
		os.Setenv("TEST_NAME", "env_only")
		os.Setenv("TEST_VALUE", "456")
		os.Setenv("TEST_NESTED_ENABLED", "true")
		defer func() {
			os.Unsetenv("TEST_NAME")
			os.Unsetenv("TEST_VALUE")
			os.Unsetenv("TEST_NESTED_ENABLED")
		}()

		tempDir := t.TempDir()
		var cfg TestConfig
		_, err := IniWithEnv("nonexistent", "yaml", []string{tempDir}, &cfg, "test")
		assert.NoError(t, err)
		assert.Equal(t, "env_only", cfg.Name)
		assert.Equal(t, 456, cfg.Value)
		assert.True(t, cfg.Nested.Enabled)
	})

	t.Run("success with multiple paths", func(t *testing.T) {
		tempDir1 := t.TempDir()
		tempDir2 := t.TempDir()

		// Create config in second directory
		configFile := filepath.Join(tempDir2, "test_config.yaml")
		configContent := `name: from_second_dir`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		var cfg TestConfig
		_, err = IniWithEnv("test_config", "yaml", []string{tempDir1, tempDir2}, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, "from_second_dir", cfg.Name)
	})
}
