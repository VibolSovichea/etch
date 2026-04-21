package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".scripture"
	configFile = "config.yaml"
)

type Config struct {
	VaultPath string `yaml:"vault_path"`
}

// configPath returns the path to the global config file.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

// Load reads the config from disk. Returns nil if no config exists yet.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Init creates the vault directory structure and writes the config.
func Init(vaultPath string) (*Config, error) {
	// Expand ~ if needed
	if len(vaultPath) >= 2 && vaultPath[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		vaultPath = filepath.Join(home, vaultPath[2:])
	}

	// Create vault directories
	dirs := []string{
		vaultPath,
		filepath.Join(vaultPath, "notes"),
		filepath.Join(vaultPath, "daily"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, err
		}
	}

	// Create internal dir for trash
	internal := filepath.Join(vaultPath, ".scripture")
	if err := os.MkdirAll(filepath.Join(internal, "trash"), 0755); err != nil {
		return nil, err
	}

	cfg := &Config{VaultPath: vaultPath}

	// Write global config
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	return cfg, nil
}
