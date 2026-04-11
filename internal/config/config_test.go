package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTemp(t, `{"interval_seconds":5,"log_format":"text","rules":[]}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSeconds != 5 {
		t.Errorf("expected interval 5, got %d", cfg.IntervalSeconds)
	}
	if cfg.LogFormat != "text" {
		t.Errorf("expected log_format text, got %q", cfg.LogFormat)
	}
}

func TestLoad_Defaults(t *testing.T) {
	p := writeTemp(t, `{}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.IntervalSeconds != 10 {
		t.Errorf("expected default interval 10, got %d", cfg.IntervalSeconds)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("expected default log_format json, got %q", cfg.LogFormat)
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	p := writeTemp(t, `{"interval_seconds":0}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for interval_seconds=0")
	}
}

func TestLoad_InvalidLogFormat(t *testing.T) {
	p := writeTemp(t, `{"log_format":"xml"}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for log_format=xml")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_WithRules(t *testing.T) {
	p := writeTemp(t, `{"rules":[{"port":22,"protocol":"tcp","action":"allow"}]}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(cfg.Rules))
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.IntervalSeconds != 10 {
		t.Errorf("expected 10, got %d", cfg.IntervalSeconds)
	}
}
