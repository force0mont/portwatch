package correlator_test

import (
	"sync"
	"testing"
	"time"

	"github.com/iamcalledrob/portwatch/internal/correlator"
	"github.com/iamcalledrob/portwatch/internal/scanner"
)

func TestConcurrent_Record_NoRace(t *testing.T) {
	c := correlator.New(100 * time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p := scanner.Port{Port: uint16(8000 + i), Protocol: "tcp", Addr: "0.0.0.0"}
			c.Record("burst", p)
		}(i)
	}
	wg.Wait()
	if c.Len() == 0 {
		t.Fatal("expected at least one active group")
	}
}

func TestFlush_AfterRealWindow_RemovesGroup(t *testing.T) {
	c := correlator.New(50 * time.Millisecond)
	c.Record("key", scanner.Port{Port: 9090, Protocol: "tcp", Addr: "0.0.0.0"})
	time.Sleep(120 * time.Millisecond)
	expired := c.Flush()
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired group after real sleep, got %d", len(expired))
	}
}
