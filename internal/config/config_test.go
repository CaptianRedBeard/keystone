package config

import (
	"bytes"
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

func TestPrintConfigSecrets(t *testing.T) {
	Cfg.Providers = map[string]ProviderConfig{
		"venice": {APIKey: "SECRET123", BaseURL: "https://api.venice.com", Model: "default"},
	}
	Cfg.Default.Provider = "venice"
	Cfg.Default.Model = "default"
	Cfg.LogLevel = "info"
	Cfg.Storage = "/tmp/db.sqlite"

	// Capture stdout
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	PrintConfig(false)

	// Close writer and read output
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	output := buf.String()
	if output == "" {
		t.Error("PrintConfig produced no output")
	}
	if bytes.Contains([]byte(output), []byte("SECRET123")) {
		t.Error("PrintConfig revealed secret API key when showSecrets=false")
	}

	// Now test with secrets
	r2, w2, _ := os.Pipe()
	os.Stdout = w2

	PrintConfig(true)

	w2.Close()
	var buf2 bytes.Buffer
	buf2.ReadFrom(r2)
	os.Stdout = old

	output2 := buf2.String()
	if !bytes.Contains([]byte(output2), []byte("SECRET123")) {
		t.Error("PrintConfig did not reveal secret API key when showSecrets=true")
	}
}
