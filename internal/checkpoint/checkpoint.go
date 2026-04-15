// Package checkpoint persists and restores the last-known set of open ports
// so that portwatch can distinguish genuinely new listeners from ports that
// were already open before the daemon started.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry is a single persisted port record.
type Entry struct {
	Proto    string    `json:"proto"`
	Addr     string    `json:"addr"`
	Port     uint16    `json:"port"`
	SavedAt  time.Time `json:"saved_at"`
}

// Checkpoint loads and saves a snapshot of known ports to a JSON file.
type Checkpoint struct {
	mu   sync.RWMutex
	path string
	data map[string]Entry
}

// New returns a Checkpoint backed by the file at path.
// If the file does not exist the checkpoint starts empty.
func New(path string) (*Checkpoint, error) {
	c := &Checkpoint{path: path, data: make(map[string]Entry)}
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return c, nil
}

func key(proto, addr string, port uint16) string {
	return proto + "|" + addr + "|" + itoa(port)
}

func itoa(n uint16) string {
	return string(rune('0'+n/10000%10)) +
		string(rune('0'+n/1000%10)) +
		string(rune('0'+n/100%10)) +
		string(rune('0'+n/10%10)) +
		string(rune('0'+n%10))
}

// Contains reports whether the given port tuple was present at the last save.
func (c *Checkpoint) Contains(proto, addr string, port uint16) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key(proto, addr, port)]
	return ok
}

// Save replaces the persisted set with entries and writes to disk.
func (c *Checkpoint) Save(entries []Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().UTC()
	next := make(map[string]Entry, len(entries))
	for _, e := range entries {
		e.SavedAt = now
		next[key(e.Proto, e.Addr, e.Port)] = e
	}
	c.data = next
	return c.flush()
}

// Len returns the number of entries currently held.
func (c *Checkpoint) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *Checkpoint) load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer f.Close()
	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return err
	}
	for _, e := range entries {
		c.data[key(e.Proto, e.Addr, e.Port)] = e
	}
	return nil
}

func (c *Checkpoint) flush() error {
	entries := make([]Entry, 0, len(c.data))
	for _, e := range c.data {
		entries = append(entries, e)
	}
	f, err := os.CreateTemp("", "checkpoint-*")
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), c.path)
}
