// Package config provides layered configuration for linden-cli.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds the resolved configuration.
type Config struct {
	BaseURL   string `json:"base_url"`
	AccountID string `json:"account_id"`
	Format    string `json:"format"`

	// Sources tracks where each value came from (not serialized).
	Sources map[string]string `json:"-"`
}

// FlagOverrides holds command-line flag values.
type FlagOverrides struct {
	Account string
	Format  string
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		BaseURL: "https://api.mylinden.family",
		Sources: make(map[string]string),
	}
}

// GlobalConfigDir returns the global config directory.
func GlobalConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "linden")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".linden"
	}
	return filepath.Join(home, ".config", "linden")
}

// Load builds the merged Config by layering all sources.
// Precedence (lowest → highest): default → global → local → env → flags
func Load(overrides FlagOverrides) (*Config, error) {
	cfg := Default()

	// Global config
	globalPath := filepath.Join(GlobalConfigDir(), "config.json")
	_ = loadFromFile(cfg, globalPath, "global")

	// Local config: .linden/config.json in CWD
	localPath := filepath.Join(".linden", "config.json")
	_ = loadFromFile(cfg, localPath, "local")

	// Environment variables
	LoadFromEnv(cfg)

	// Flag overrides
	ApplyOverrides(cfg, overrides)

	return cfg, nil
}

// LoadFromEnv applies environment variable overrides.
func LoadFromEnv(cfg *Config) {
	if v := os.Getenv("LINDEN_BASE_URL"); v != "" {
		cfg.BaseURL = v
		cfg.Sources["base_url"] = "env"
	}
	if v := os.Getenv("LINDEN_ACCOUNT"); v != "" {
		cfg.AccountID = v
		cfg.Sources["account_id"] = "env"
	}
}

// ApplyOverrides applies flag value overrides.
func ApplyOverrides(cfg *Config, overrides FlagOverrides) {
	if overrides.Account != "" {
		cfg.AccountID = overrides.Account
		cfg.Sources["account_id"] = "flag"
	}
	if overrides.Format != "" {
		cfg.Format = overrides.Format
		cfg.Sources["format"] = "flag"
	}
}

// PersistValue saves a config key/value to the chosen scope.
// scope must be "global" or "local".
func PersistValue(key, value, scope string) error {
	var configPath string
	switch scope {
	case "global":
		configPath = filepath.Join(GlobalConfigDir(), "config.json")
	case "local":
		configPath = filepath.Join(".linden", "config.json")
	default:
		return fmt.Errorf("invalid scope %q: must be \"global\" or \"local\"", scope)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data := make(map[string]any)
	if raw, err := os.ReadFile(configPath); err == nil { //nolint:gosec
		_ = json.Unmarshal(raw, &data)
	}

	data[key] = value

	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	return atomicWriteFile(configPath, append(encoded, '\n'))
}

// loadFromFile merges a JSON config file into cfg, tracking source.
func loadFromFile(cfg *Config, path, source string) error {
	raw, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return err // file not found is normal
	}

	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	if v, ok := data["base_url"].(string); ok && v != "" {
		cfg.BaseURL = v
		cfg.Sources["base_url"] = source
	}
	if v, ok := data["account_id"].(string); ok && v != "" {
		cfg.AccountID = v
		cfg.Sources["account_id"] = source
	}
	if v, ok := data["format"].(string); ok && v != "" {
		cfg.Format = v
		cfg.Sources["format"] = source
	}

	return nil
}

// atomicWriteFile writes data to path using a temp file + rename.
func atomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		if runtime.GOOS == "windows" {
			_ = os.Remove(path)
			return os.Rename(tmpPath, path)
		}
		os.Remove(tmpPath)
		return err
	}
	return nil
}
