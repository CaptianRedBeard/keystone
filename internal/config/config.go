package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// ErrNoConfig is returned when the config file does not exist.
var ErrNoConfig = errors.New("config file not found")

// Config holds Keystone configuration values.
type Config struct {
	DBPath     string            `yaml:"db_path"`
	AgentsDir  string            `yaml:"agents_dir"`
	JSONOutput bool              `yaml:"json_output,omitempty"`
	Secrets    map[string]string `yaml:"secrets,omitempty"`
}

// New returns a config populated with defaults, optionally overridden by environment variables.
func New() *Config {
	cfg := &Config{
		DBPath:     "",
		AgentsDir:  "./agents",
		JSONOutput: false,
		Secrets:    make(map[string]string),
	}

	if envDir := os.Getenv("KEYSTONE_AGENTS_DIR"); envDir != "" {
		cfg.AgentsDir = envDir
	}

	return cfg
}

// Load reads and unmarshals a YAML config from disk.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoConfig
		}
		return nil, err
	}

	cfg := New()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save marshals and writes the config to the given path.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
