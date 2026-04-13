// Package resolver maps raw port numbers to process names by reading
// /proc/<pid>/net/tcp (and udp) alongside /proc/<pid>/fd symlinks.
package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// ProcessInfo holds the resolved owner of a listening socket.
type ProcessInfo struct {
	PID  int
	Name string
}

// Resolver looks up which process owns a given inode.
type Resolver struct {
	mu      sync.Mutex
	procDir string
}

// New returns a Resolver that reads from the real /proc filesystem.
func New() *Resolver {
	return &Resolver{procDir: "/proc"}
}

// newWithProcDir is used in tests to inject a fake /proc directory.
func newWithProcDir(dir string) *Resolver {
	return &Resolver{procDir: dir}
}

// Lookup returns the ProcessInfo for the process that owns inode, or an error
// if no matching process is found.
func (r *Resolver) Lookup(inode uint64) (ProcessInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entries, err := os.ReadDir(r.procDir)
	if err != nil {
		return ProcessInfo{}, fmt.Errorf("resolver: read proc: %w", err)
	}

	target := fmt.Sprintf("socket:[%d]", inode)

	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // skip non-numeric entries
		}
		fdDir := filepath.Join(r.procDir, e.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				name := r.readComm(pid)
				return ProcessInfo{PID: pid, Name: name}, nil
			}
		}
	}
	return ProcessInfo{}, fmt.Errorf("resolver: inode %d not found", inode)
}

func (r *Resolver) readComm(pid int) string {
	data, err := os.ReadFile(filepath.Join(r.procDir, strconv.Itoa(pid), "comm"))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}
