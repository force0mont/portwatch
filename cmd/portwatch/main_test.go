package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMain_Version verifies that the binary prints a version string and exits 0.
func TestMain_Version(t *testing.T) {
	if os.Getenv("PORTWATCH_RUN_INTEGRATION") == "" {
		t.Skip("set PORTWATCH_RUN_INTEGRATION=1 to run binary integration tests")
	}

	binary := buildBinary(t)
	cmd := exec.Command(binary, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v\noutput: %s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty version output")
	}
}

// TestMain_InvalidConfig verifies that a bad config path causes exit 1.
func TestMain_InvalidConfig(t *testing.T) {
	if os.Getenv("PORTWATCH_RUN_INTEGRATION") == "" {
		t.Skip("set PORTWATCH_RUN_INTEGRATION=1 to run binary integration tests")
	}

	binary := buildBinary(t)
	cmd := exec.Command(binary, "-config", "/nonexistent/path/config.json")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	out := filepath.Join(tmp, "portwatch")
	cmd := exec.Command("go", "build", "-o", out, ".")
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, b)
	}
	return out
}
