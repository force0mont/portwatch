package dedupe_test

import (
	"sync"
	"testing"
	"time"

	"github.com/iamcalledned/portwatch/internal/dedupe"
	"github.com/iamcalledned/portwatch/internal/scanner"
)

func TestDedupe_ConcurrentAccess_NoRace(t *testing.T) {
	d := dedupe.New(100 * time.Millisecond)
	p := scanner.Port{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.IsDuplicate(p, "alert")
		}()
	}
	wg.Wait()
}

func TestDedupe_MultiplePortsUnderLoad(t *testing.T) {
	d := dedupe.New(500 * time.Millisecond)

	ports := []scanner.Port{
		{Port: 80, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 443, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 53, Protocol: "udp", Address: "0.0.0.0"},
	}

	// First round — none should be duplicates.
	for _, p := range ports {
		if d.IsDuplicate(p, "alert") {
			t.Errorf("port %d: expected false on first call", p.Port)
		}
	}

	// Second round — all should be duplicates (window still open).
	for _, p := range ports {
		if !d.IsDuplicate(p, "alert") {
			t.Errorf("port %d: expected true within window", p.Port)
		}
	}

	// After reset — none should be duplicates.
	d.Reset()
	for _, p := range ports {
		if d.IsDuplicate(p, "alert") {
			t.Errorf("port %d: expected false after Reset", p.Port)
		}
	}
}
