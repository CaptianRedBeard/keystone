package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfgPath := filepath.Join(os.TempDir(), "nonexistent_config.yaml")
	err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if Cfg.Default.Provider != "venice" {
		t.Errorf("expected default provider 'venice', got '%s'", Cfg.Default.Provider)
	}
	if Cfg.Default.Model != "default" {
		t.Errorf("expected default model 'default', got '%s'", Cfg.Default.Model)
	}
	if Cfg.LogLevel != "info" {
		t.Errorf("expected logLevel 'info', got '%s'", Cfg.LogLevel)
	}
	if Cfg.Storage == "" {
		t.Error("expected Storage to have a default value")
	}
}

func TestSaveAndReloadConfig(t *testing.T) {
	// Modify in-memory config
	Cfg.Default.Provider = "mock_provider"
	Cfg.Default.Model = "mock_model"
	Cfg.LogLevel = "debug"

	tempPath := filepath.Join(os.TempDir(), "keystone_test_config.yaml")
	defer os.Remove(tempPath)

	// Save
	if err := SaveConfig(tempPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Reset in-memory config
	Cfg.Default.Provider = ""
	Cfg.Default.Model = ""
	Cfg.LogLevel = ""

	// Reload
	if err := LoadConfig(tempPath); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if Cfg.Default.Provider != "mock_provider" {
		t.Errorf("expected provider 'mock_provider', got '%s'", Cfg.Default.Provider)
	}
	if Cfg.Default.Model != "mock_model" {
		t.Errorf("expected model 'mock_model', got '%s'", Cfg.Default.Model)
	}
	if Cfg.LogLevel != "debug" {
		t.Errorf("expected logLevel 'debug', got '%s'", Cfg.LogLevel)
	}
}

func TestEnvOverrides(t *testing.T) {
	// Set environment variable
	os.Setenv("KEYSTONE_DEFAULT_PROVIDER", "env_provider")
	defer os.Unsetenv("KEYSTONE_DEFAULT_PROVIDER")

	cfgPath := filepath.Join(os.TempDir(), "nonexistent_config.yaml")
	if err := LoadConfig(cfgPath); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if Cfg.Default.Provider != "env_provider" {
		t.Errorf("expected provider from env 'env_provider', got '%s'", Cfg.Default.Provider)
	}
}
