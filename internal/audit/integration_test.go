package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/your/portwatch/_CreatesFile(t *testing.T)  := t.TempDir()
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Record(audit.ActionPortAppeared, 9090, "udp", "0.0.0.0", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	var entry audit.Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if entry.Port != 9090 {
		t.Errorf("port = %d, want 9090", entry.Port)
	}
}

func TestNew_AppendsToExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	for i := 0; i < 3; i++ {
		l, err := audit.New(path)
		if err != nil {
			t.Fatalf("New iteration %d: %v", i, err)
		}
		_ = l.Record(audit.ActionAlertSent, uint16(8000+i), "tcp", "0.0.0.0", "")
	}

	f, _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}

func TestRecord_ConcurrentWrites_NoRace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	l, _ := audit.New(path)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = l.Record(audit.ActionPortAppeared, uint16(1024+i), "tcp", "0.0.0.0", "")
		}(i)
	}
	wg.Wait()
}
