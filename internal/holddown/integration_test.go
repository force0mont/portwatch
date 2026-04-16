package holddown_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/portwatch/internal/holddown"
)

func TestConcurrent_Seen_NoRace(t *testing.T) {
	h := holddown.New(10 * time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		port := uint16(8000 + i)
		go func(p uint16) {
			defer wg.Done()
			h.Seen(p, "tcp")
			h.Seen(p, "tcp")
			h.Gone(p, "tcp")
		}(port)
	}
	wg.Wait()
}

func TestConcurrent_PruneAndSeen_NoRace(t *testing.T) {
	h := holddown.New(1 * time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			h.Seen(uint16(9000+idx), "udp")
			time.Sleep(2 * time.Millisecond)
			h.Prune()
		}(i)
	}
	wg.Wait()
}
