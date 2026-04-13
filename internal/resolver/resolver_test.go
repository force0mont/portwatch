package resolver

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// buildFakeProc creates a minimal fake /proc tree for testing.
// It writes comm and creates a symlink under fd that points to the given inode.
func buildFakeProc(t *testing.T, pid int, comm string, inode uint64) string {
	t.Helper()
	root := t.TempDir()
	pidDir := filepath.Join(root, strconv.Itoa(pid))
	fdDir := filepath.Join(pidDir, "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := "socket:[" + strconv.FormatUint(inode, 10) + "]"
	if err := os.Symlink(target, filepath.Join(fdDir, "3")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestLookup_FindsProcess(t *testing.T) {
	const pid = 1234
	const inode = 99887766
	root := buildFakeProc(t, pid, "nginx", inode)
	r := newWithProcDir(root)

	info, err := r.Lookup(inode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != pid {
		t.Errorf("PID: got %d, want %d", info.PID, pid)
	}
	if info.Name != "nginx" {
		t.Errorf("Name: got %q, want %q", info.Name, "nginx")
	}
}

func TestLookup_UnknownInode_ReturnsError(t *testing.T) {
	root := buildFakeProc(t, 42, "sshd", 111)
	r := newWithProcDir(root)

	_, err := r.Lookup(999) // different inode
	if err == nil {
		t.Fatal("expected error for unknown inode, got nil")
	}
}

func TestLookup_MissingCommFile_ReturnsUnknown(t *testing.T) {
	const inode = 55443322
	root := buildFakeProc(t, 7, "sshd", inode)
	// Remove comm so readComm falls back to "unknown".
	if err := os.Remove(filepath.Join(root, "7", "comm")); err != nil {
		t.Fatal(err)
	}
	r := newWithProcDir(root)

	info, err := r.Lookup(inode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "unknown" {
		t.Errorf("Name: got %q, want \"unknown\"", info.Name)
	}
}

func TestNew_NotNil(t *testing.T) {
	if New() == nil {
		t.Fatal("New() returned nil")
	}
}
