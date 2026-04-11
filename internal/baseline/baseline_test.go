package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danvolchek/portwatch/internal/baseline"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

var sampleEntries = []baseline.Entry{
	{Protocol: "tcp", Address: "0.0.0.0", Port: 22},
	{Protocol: "tcp", Address: "127.0.0.1", Port: 8080},
	{Protocol: "udp", Address: "0.0.0.0", Port: 53},
}

func TestNew_EmptyWhenNoFile(t *testing.T) {
	b, err := baseline.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(b.All()); got != 0 {
		t.Errorf("expected 0 entries, got %d", got)
	}
}

func TestSet_And_Contains(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	if err := b.Set(sampleEntries); err != nil {
		t.Fatalf("Set: %v", err)
	}
	for _, e := range sampleEntries {
		if !b.Contains(e) {
			t.Errorf("expected baseline to contain %+v", e)
		}
	}
}

func TestContains_ReturnsFalseForUnknown(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Set(sampleEntries)
	unknown := baseline.Entry{Protocol: "tcp", Address: "0.0.0.0", Port: 9999}
	if b.Contains(unknown) {
		t.Error("expected Contains to return false for unknown entry")
	}
}

func TestSet_PersistsToDisk(t *testing.T) {
	path := tempPath(t)
	b, _ := baseline.New(path)
	_ = b.Set(sampleEntries)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("baseline file not written: %v", err)
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	path := tempPath(t)

	// Write baseline with first instance.
	b1, _ := baseline.New(path)
	_ = b1.Set(sampleEntries)

	// Load with a second instance.
	b2, err := baseline.New(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got := len(b2.All()); got != len(sampleEntries) {
		t.Errorf("expected %d entries after reload, got %d", len(sampleEntries), got)
	}
	for _, e := range sampleEntries {
		if !b2.Contains(e) {
			t.Errorf("reloaded baseline missing %+v", e)
		}
	}
}

func TestSet_ReplacesExistingEntries(t *testing.T) {
	b, _ := baseline.New(tempPath(t))
	_ = b.Set(sampleEntries)
	newEntries := []baseline.Entry{{Protocol: "tcp", Address: "0.0.0.0", Port: 443}}
	_ = b.Set(newEntries)

	if b.Contains(sampleEntries[0]) {
		t.Error("old entry should have been removed after Set")
	}
	if !b.Contains(newEntries[0]) {
		t.Error("new entry should be present after Set")
	}
}
