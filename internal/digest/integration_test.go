package digest_test

import (
	"sync"
	"testing"

	"github.com/stevezaluk/portwatch/internal/digest"
)

// TestConcurrent_Changed_NoRace verifies that concurrent calls to Changed
// from multiple goroutines do not trigger the race detector.
func TestConcurrent_Changed_NoRace(t *testing.T) {
	d := digest.New()
	e := []digest.Entry{
		{Proto: "tcp", Address: "0.0.0.0", Port: 80},
		{Proto: "udp", Address: "0.0.0.0", Port: 53},
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Changed(e)
			_ = d.Last()
		}()
	}
	wg.Wait()
}

// TestDigest_RoundTrip_ConsistentAcrossInstances ensures two independent
// Digesters produce the same hash for identical entry sets.
func TestDigest_RoundTrip_ConsistentAcrossInstances(t *testing.T) {
	entries := []digest.Entry{
		{Proto: "tcp", Address: "0.0.0.0", Port: 443},
		{Proto: "tcp", Address: "127.0.0.1", Port: 22},
	}

	d1 := digest.New()
	d2 := digest.New()

	d1.Changed(entries)
	d2.Changed(entries)

	if d1.Last() != d2.Last() {
		t.Fatalf("digests differ: %q vs %q", d1.Last(), d2.Last())
	}
}
