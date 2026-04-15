package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_EmptyWhenNoFile(t *testing.T) {
	c, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", c.Len())
	}
}

func TestContains_ReturnsFalseForUnknown(t *testing.T) {
	c, _ := New(tempPath(t))
	if c.Contains("tcp", "0.0.0.0", 8080) {
		t.Fatal("expected false for unknown port")
	}
}

func TestSave_And_Contains(t *testing.T) {
	path := tempPath(t)
	c, _ := New(path)
	entries := []Entry{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 443},
		{Proto: "udp", Addr: "127.0.0.1", Port: 53},
	}
	if err := c.Save(entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if !c.Contains("tcp", "0.0.0.0", 443) {
		t.Error("expected tcp/443 to be present")
	}
	if !c.Contains("udp", "127.0.0.1", 53) {
		t.Error("expected udp/53 to be present")
	}
	if c.Contains("tcp", "0.0.0.0", 80) {
		t.Error("expected tcp/80 to be absent")
	}
}

func TestSave_PersistsToDisk(t *testing.T) {
	path := tempPath(t)
	c1, _ := New(path)
	_ = c1.Save([]Entry{{Proto: "tcp", Addr: "0.0.0.0", Port: 22}})

	c2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !c2.Contains("tcp", "0.0.0.0", 22) {
		t.Error("expected tcp/22 to survive reload")
	}
}

func TestSave_ReplacesOldEntries(t *testing.T) {
	path := tempPath(t)
	c, _ := New(path)
	_ = c.Save([]Entry{{Proto: "tcp", Addr: "0.0.0.0", Port: 80}})
	_ = c.Save([]Entry{{Proto: "tcp", Addr: "0.0.0.0", Port: 443}})

	if c.Contains("tcp", "0.0.0.0", 80) {
		t.Error("port 80 should have been replaced")
	}
	if !c.Contains("tcp", "0.0.0.0", 443) {
		t.Error("port 443 should be present")
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", c.Len())
	}
}

func TestNew_ErrorOnCorruptFile(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json{"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
