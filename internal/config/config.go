// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/portwatch/internal/filter"
)

// Config holds the complete portwatch runtime configuration.
type Config struct {
	// ScanInterval controls how often the port scanner runs.
	ScanInterval Duration `json:"scan_interval"`
	// LogFormat is either "text" or "json".
	LogFormat string `json:"log_format"`
	// ProcPath is the root used to locate /proc/net files (default "/proc").
	ProcPath string `json:"proc_path"`
	// SuppressRules lists ports/CIDRs that should never produce alerts.
	SuppressRules []filter.Rule `json:"suppress_rules,omitempty"`
}

// Duration is a time.Duration that marshals/unmarshals as a string (e.g. "30s").
type Duration struct{ time.Duration }

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("config: invalid duration %q: %w", s, err)
	}
	d.Duration = v
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ScanInterval: Duration{30 * time.Second},
		LogFormat:    "text",
		ProcPath:     "/proc",
	}
}

// Load reads a JSON config file from path, applying defaults for missing fields.
// If path is empty the default configuration is returned.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("config: decode: %w", err)
	}
	return cfg, validate(cfg)
}

func validate(cfg Config) error {
	if cfg.ScanInterval.Duration <= 0 {
		return errors.New("config: scan_interval must be positive")
	}
	if cfg.LogFormat != "text" && cfg.LogFormat != "json" {
		return fmt.Errorf("config: log_format must be \"text\" or \"json\", got %q", cfg.LogFormat)
	}
	return nil
}
