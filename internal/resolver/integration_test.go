package resolver_test

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/yourusername/portwatch/internal/resolver"
)

// buildFakeRoot constructs a /proc-like tree with multiple PIDs.
func buildFakeRoot(t *testing.T, entries map[int]struct{ comm string; inode uint64 }) string {
	t.Helper()
	root := t.TempDir()
	for pid, info := range entries {
		pidStr := strconv.Itoa(pid)
		fdDir := filepath.Join(root, pidStr, "fd")
		if err := os.MkdirAll(fdDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, pidStr, "comm"), []byte(info.comm), 0o644); err != nil {
			t.Fatal(err)
		}
		target := "socket:[" + strconv.FormatUint(info.inode, 10) + "]"
		if err := os.Symlink(target, filepath.Join(fdDir, "4")); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestConcurrent_Lookup_NoRace(t *testing.T) {
	entries := map[int]struct{ comm string; inode uint64 }{
		100: {"nginx", 1001},
		101: {"sshd", 1002},
		102: {"postgres", 1003},
	}
	root := buildFakeRoot(t, entries)

	// Use the exported constructor via the internal test helper path.
	// We access the package via the same module path used in go.mod.
	_ = resolver.New() // ensure package is reachable

	// Rebuild with our fake root via a white-box approach:
	// We test concurrency using the exported New() + real /proc,
	// but guard against missing /proc by skipping gracefully.
	if _, err := os.Stat("/proc"); os.IsNotExist(err) {
		t.Skip("/proc not available")
	}
	_ = root // used in non-exported variant; keep build happy

	r := resolver.New()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// inode 0 will never match; we just exercise concurrency.
			_, _ = r.Lookup(0)
		}()
	}
	wg.Wait()
}
