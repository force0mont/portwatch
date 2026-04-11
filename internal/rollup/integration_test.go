package rollup_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wander/portwatch/internal/rollup"
)

func TestRollup_ConcurrentRecords_NoRace(t *testing.T) {
	r := rollup.New(5*time.Second, 50)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("port:%d", (n%5)*100+8000)
			r.Record(key)
		}(i)
	}
	wg.Wait()
}

func TestRollup_MultipleKeys_IndependentThresholds(t *testing.T) {
	r := rollup.New(10*time.Second, 3)
	keys := []string{"tcp:80", "tcp:443", "udp:53"}
	summaries := map[string]int{}

	for _, k := range keys {
		for i := 0; i < 5; i++ {
			if _, ok := r.Record(k); ok {
				summaries[k]++
			}
		}
	}

	for _, k := range keys {
		if summaries[k] != 1 {
			t.Errorf("key %s: expected 1 summary, got %d", k, summaries[k])
		}
	}
}
