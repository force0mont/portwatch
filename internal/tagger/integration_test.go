package tagger_test

import (
	"sync"
	"testing"

	"github.com/your-org/portwatch/internal/tagger"
)

// TestConcurrent_TagAndOverride_NoRace verifies that concurrent reads
// and writes do not trigger the race detector.
func TestConcurrent_TagAndOverride_NoRace(t *testing.T) {
	tg := tagger.New()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			port := uint16(1024 + n)
			tg.Override(port, "svc")
			_ = tg.Tag(port)
			_ = tg.Known(port)
		}(i)
	}
	wg.Wait()
}

// TestWellKnown_FullTable checks that every port in the built-in table
// is reported as known without any overrides.
func TestWellKnown_FullTable(t *testing.T) {
	knownPorts := []uint16{22, 25, 53, 80, 110, 143, 443, 3306, 5432, 6379, 8080, 8443, 27017}
	tg := tagger.New()

	for _, p := range knownPorts {
		if !tg.Known(p) {
			t.Errorf("port %d should be known", p)
		}
		if tg.Tag(p) == "unknown" {
			t.Errorf("port %d should not be tagged unknown", p)
		}
	}
}
