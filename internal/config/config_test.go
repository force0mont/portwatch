package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/filter"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ScanInterval.Duration != 30*time.Second {
		t.Errorf("default interval: got %v", cfg.ScanInterval.Duration)
	}
	if cfg.LogFormat != "text" {
		t.Errorf("default log_format: got %q", cfg.LogFormat)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTemp(t, `{"scan_interval":"10s","log_format":"json"}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ScanInterval.Duration != 10*time.Second {
		t.Errorf("interval: got %v", cfg.ScanInterval.Duration)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("log_format: got %q", cfg.LogFormat)
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	p := writeTemp(t, `{"scan_interval":"-5s"}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestLoad_InvalidLogFormat(t *testing.T) {
	p := writeTemp(t, `{"scan_interval":"5s","log_format":"xml"}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid log_format")
	}
}

func TestLoad_SuppressRules(t *testing.T) {
	rules := []filter.Rule{
		{Port: 22, Protocol: "tcp"},
		{CIDR: "127.0.0.0/8"},
	}
	b, _ := json.Marshal(map[string]interface{}{
		"scan_interval":  "15s",
		"log_format":     "text",
		"suppress_rules": rules,
	})
	p := writeTemp(t, string(b))
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.SuppressRules) != 2 {
		t.Errorf("expected 2 suppress rules, got %d", len(cfg.SuppressRules))
	}
	if cfg.SuppressRules[0].Port != 22 {
		t.Errorf("first rule port: got %d", cfg.SuppressRules[0].Port)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
