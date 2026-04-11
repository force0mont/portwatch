// Package baseline persists a known-good set of listening ports to disk
// so that portwatch can distinguish truly new listeners from ports that were
// already open when the daemon first started.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single baselined port.
type Entry struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
}

// Baseline holds the persisted set of known-good entries.
type Baseline struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
	SavedAt time.Time `json:"saved_at"`
}

type diskFormat struct {
	SavedAt time.Time `json:"saved_at"`
	Entries []Entry   `json:"entries"`
}

// New creates an empty Baseline backed by the given file path.
// If the file already exists its contents are loaded.
func New(path string) (*Baseline, error) {
	b := &Baseline{
		entries: make(map[string]Entry),
		path:    path,
	}
	if _, err := os.Stat(path); err == nil {
		if err := b.load(); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Contains reports whether the entry is part of the baseline.
func (b *Baseline) Contains(e Entry) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[keyFor(e)]
	return ok
}

// Set replaces the entire baseline with the provided entries and persists it.
func (b *Baseline) Set(entries []Entry) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = make(map[string]Entry, len(entries))
	for _, e := range entries {
		b.entries[keyFor(e)] = e
	}
	b.SavedAt = time.Now().UTC()
	return b.save()
}

// All returns a snapshot of all baselined entries.
func (b *Baseline) All() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

func (b *Baseline) save() error {
	df := diskFormat{SavedAt: b.SavedAt, Entries: b.All()}
	data, err := json.MarshalIndent(df, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o600)
}

func (b *Baseline) load() error {
	data, err := os.ReadFile(b.path)
	if err != nil {
		return err
	}
	var df diskFormat
	if err := json.Unmarshal(data, &df); err != nil {
		return err
	}
	b.SavedAt = df.SavedAt
	for _, e := range df.Entries {
		b.entries[keyFor(e)] = e
	}
	return nil
}

func keyFor(e Entry) string {
	return e.Protocol + "|" + e.Address + "|" + string(rune(e.Port))
}
