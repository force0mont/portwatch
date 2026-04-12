package ledger

import (
	"sync"
	"testing"
)

func TestConcurrent_RecordAppeared_NoRace(t *testing.T) {
	l := New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(port uint16) {
			defer wg.Done()
			l.RecordAppeared(port, "tcp")
			l.RecordDisappeared(port, "tcp")
		}(uint16(8000 + i))
	}

	wg.Wait()

	all := l.All()
	if len(all) != 50 {
		t.Fatalf("expected 50 entries, got %d", len(all))
	}
	for _, e := range all {
		if e.Appeared != 1 {
			t.Errorf("key %s: expected Appeared=1, got %d", e.Key, e.Appeared)
		}
		if e.Disappeared != 1 {
			t.Errorf("key %s: expected Disappeared=1, got %d", e.Key, e.Disappeared)
		}
	}
}

func TestConcurrent_GetAndReset_NoRace(t *testing.T) {
	l := New()
	l.RecordAppeared(80, "tcp")

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			l.Get(80, "tcp")
		}()
		go func() {
			defer wg.Done()
			l.All()
		}()
	}
	wg.Wait()
}
