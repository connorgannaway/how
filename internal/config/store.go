package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

// Get the path to the config file
func GetConfigPath() (string, error) {
        // Prefer XDG_CONFIG_HOME if set
        if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
                howDir := filepath.Join(xdgConfig, "how")
                if err := os.MkdirAll(howDir, 0700); err != nil {
                        return "", err
                }
                return filepath.Join(howDir, "config.json"), nil
        }

		// No XDG_CONFIG_HOME, use platform defaults

        var configDir string

        // Use ~/.config on macOS. os.UserConfigDir returns ~/Library/Application Support on macOS
        if runtime.GOOS == "darwin" {
                homeDir, err := os.UserHomeDir()
                if err != nil {
                        return "", err
                }
                configDir = filepath.Join(homeDir, ".config")
        } else {
                // Use platform default for Windows (%appdata%), Linux ($HOME/.config) and others
                var err error
                configDir, err = os.UserConfigDir()
                if err != nil {
                        return "", err
                }
        }

        // Create how subdirectory
        howDir := filepath.Join(configDir, "how")
        if err := os.MkdirAll(howDir, 0700); err != nil {
                return "", err
        }

        return filepath.Join(howDir, "config.json"), nil
  }

// Load the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return a new empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return NewConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Validate current provider if set
	if config.CurrentProvider != "" {
		validProviders := GetProviders()
		valid := slices.Contains(validProviders, config.CurrentProvider)
		if !valid {
			return nil, fmt.Errorf("invalid provider: %s", config.CurrentProvider)
		}
	}

	return &config, nil
}

// Save the configuration to disk
func Save(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
