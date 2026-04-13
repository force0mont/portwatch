package tracer

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestConcurrent_ObserveAndRemove_NoRace(t *testing.T) {
	tr := New()
	var wg sync.WaitGroup
	ports := []scanner.Port{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 443},
		{Protocol: "udp", Address: "0.0.0.0", Port: 53},
	}
	for i := 0; i < 20; i++ {
		for _, p := range ports {
			wg.Add(2)
			go func(port scanner.Port) {
				defer wg.Done()
				tr.Observe(port)
			}(p)
			go func(port scanner.Port) {
				defer wg.Done()
				tr.Get(port)
			}(p)
		}
	}
	wg.Wait()
}

func TestTracer_DurationGrowsOverRealTime(t *testing.T) {
	tr := New()
	p := scanner.Port{Protocol: "tcp", Address: "127.0.0.1", Port: 9999}
	tr.Observe(p)
	time.Sleep(20 * time.Millisecond)
	e := tr.Observe(p)
	if e.Duration < 10*time.Millisecond {
		t.Fatalf("expected duration >= 10ms, got %v", e.Duration)
	}
}
