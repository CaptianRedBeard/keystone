package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Providers map[string]ProviderConfig `mapstructure:"providers"`
	Default   DefaultConfig             `mapstructure:"default"`
	LogLevel  string                    `mapstructure:"logLevel"`
	Storage   string                    `mapstructure:"storage"`
}

type ProviderConfig struct {
	APIKey  string `mapstructure:"apiKey"`
	BaseURL string `mapstructure:"baseURL"`
	Model   string `mapstructure:"model"`
}

type DefaultConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
}

// Global in-memory config
var Cfg Config

func LoadConfig(cfgPath string) error {
	v := viper.New()

	if cfgPath != "" {
		v.SetConfigFile(cfgPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to locate home directory: %w", err)
		}
		v.AddConfigPath(filepath.Join(home, ".keystone"))
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Environment variable support: KEYSTONE_DEFAULT_PROVIDER, etc.
	v.SetEnvPrefix("KEYSTONE")
	v.AutomaticEnv()

	v.SetDefault("logLevel", "info")
	v.SetDefault("storage", filepath.Join(os.TempDir(), "keystone_usage.db"))
	v.SetDefault("default.provider", "venice")
	v.SetDefault("default.model", "default")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("⚠️ No config file found, using defaults only.")
	}

	if err := v.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func SaveConfig(cfgPath string) error {
	v := viper.New()
	v.Set("providers", Cfg.Providers)
	v.Set("default", Cfg.Default)
	v.Set("logLevel", Cfg.LogLevel)
	v.Set("storage", Cfg.Storage)

	if cfgPath == "" {
		home, _ := os.UserHomeDir()
		cfgPath = filepath.Join(home, ".keystone", "config.yaml")
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return v.WriteConfigAs(cfgPath)
}

func PrintConfig(showSecrets bool) {
	fmt.Println("Keystone Configuration:")
	fmt.Printf("  Log Level: %s\n", Cfg.LogLevel)
	fmt.Printf("  Storage:   %s\n", Cfg.Storage)
	fmt.Printf("  Default Provider: %s\n", Cfg.Default.Provider)
	fmt.Printf("  Default Model:    %s\n", Cfg.Default.Model)
	fmt.Println("  Providers:")

	for name, p := range Cfg.Providers {
		apiKey := "****"
		if showSecrets {
			apiKey = p.APIKey
		}
		fmt.Printf("   - %s\n", name)
		fmt.Printf("     BaseURL: %s\n", p.BaseURL)
		fmt.Printf("     Model:   %s\n", p.Model)
		fmt.Printf("     APIKey:  %s\n", apiKey)
	}
}
