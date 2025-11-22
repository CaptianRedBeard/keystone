package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"keystone/internal/config"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "configtest")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "config.yaml")

	cfg := &config.Config{
		DBPath:    "foo.db",
		AgentsDir: "./agents",
	}

	if err := config.Save(path, cfg); err != nil {
		t.Fatal(err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.DBPath != "foo.db" {
		t.Fatalf("expected DBPath foo.db, got %s", loaded.DBPath)
	}
	if loaded.AgentsDir != "./agents" {
		t.Fatalf("expected AgentsDir ./agents, got %s", loaded.AgentsDir)
	}
}

func TestLoadMissing(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "missing.yaml")

	_, err := config.Load(path)
	if err == nil || err != config.ErrNoConfig {
		t.Fatalf("expected ErrNoConfig, got %v", err)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "bad.yaml")

	if err := os.WriteFile(path, []byte("{not valid"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSaveFail(t *testing.T) {
	// Writing to a directory should fail
	err := config.Save("/", &config.Config{})
	if err == nil {
		t.Fatal("expected error saving to invalid path")
	}
}

func TestDefaults(t *testing.T) {
	cfg := config.New()
	if cfg.AgentsDir != "./agents" {
		t.Errorf("expected default AgentsDir './agents', got %s", cfg.AgentsDir)
	}
}

func TestEnvOverride(t *testing.T) {
	tmp := tempDir(t)
	defer os.Unsetenv("KEYSTONE_AGENTS_DIR")
	os.Setenv("KEYSTONE_AGENTS_DIR", tmp)

	cfg := config.New()
	if cfg.AgentsDir != tmp {
		t.Errorf("expected AgentsDir '%s' from env, got %s", tmp, cfg.AgentsDir)
	}
}
