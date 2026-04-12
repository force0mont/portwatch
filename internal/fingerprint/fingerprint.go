// Package fingerprint derives a stable identity string for a listening port
// based on its address, port number, and protocol. The fingerprint can be used
// to correlate observations across scan cycles without relying on ephemeral
// process IDs or timestamps.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
)

// Fingerprint is a stable, opaque identity string for a port observation.
type Fingerprint string

// Entry holds the fields used to derive a fingerprint.
type Entry struct {
	Proto   string
	Addr    net.IP
	Port    uint16
}

// Deriver computes and caches fingerprints for port entries.
type Deriver struct {
	mu    sync.Mutex
	cache map[string]Fingerprint
}

// New returns a new Deriver with an empty cache.
func New() *Deriver {
	return &Deriver{
		cache: make(map[string]Fingerprint),
	}
}

// Derive returns the fingerprint for the given entry, computing and caching it
// on first access.
func (d *Deriver) Derive(e Entry) Fingerprint {
	k := cacheKey(e)

	d.mu.Lock()
	defer d.mu.Unlock()

	if fp, ok := d.cache[k]; ok {
		return fp
	}

	fp := compute(e)
	d.cache[k] = fp
	return fp
}

// Invalidate removes a cached fingerprint for the given entry, forcing
// recomputation on the next call to Derive.
func (d *Deriver) Invalidate(e Entry) {
	k := cacheKey(e)

	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.cache, k)
}

// Len returns the number of cached fingerprints.
func (d *Deriver) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.cache)
}

func compute(e Entry) Fingerprint {
	raw := fmt.Sprintf("%s|%s|%d", e.Proto, e.Addr.String(), e.Port)
	sum := sha256.Sum256([]byte(raw))
	return Fingerprint(hex.EncodeToString(sum[:8]))
}

func cacheKey(e Entry) string {
	return fmt.Sprintf("%s:%s:%d", e.Proto, e.Addr.String(), e.Port)
}
