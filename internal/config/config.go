// Package config handles loading and parsing of portwatch configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/rules"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	// IntervalSeconds is how often the scanner runs (default: 10).
	IntervalSeconds int `json:"interval_seconds"`
	// LogFormat is either "json" or "text" (default: "json").
	LogFormat string `json:"log_format"`
	// Rules defines the allow/alert rules evaluated against open ports.
	Rules []rules.Rule `json:"rules"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		IntervalSeconds: 10,
		LogFormat:       "json",
		Rules:           []rules.Rule{},
	}
}

// Load reads a JSON config file from the given path and returns a Config.
// Missing optional fields are filled with defaults.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the config values are acceptable.
func (c *Config) Validate() error {
	if c.IntervalSeconds <= 0 {
		return fmt.Errorf("config: interval_seconds must be > 0, got %d", c.IntervalSeconds)
	}
	if c.LogFormat != "json" && c.LogFormat != "text" {
		return fmt.Errorf("config: log_format must be \"json\" or \"text\", got %q", c.LogFormat)
	}
	for i, r := range c.Rules {
		if err := rules.Validate(r); err != nil {
			return fmt.Errorf("config: rule[%d]: %w", i, err)
		}
	}
	return nil
}
