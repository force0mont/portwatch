package window_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stevenlawton/portwatch/internal/window"
)

func TestConcurrent_Add_NoRace(t *testing.T) {
	w := window.New(time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "k"
			if i%2 == 0 {
				key = "j"
			}
			w.Add(key)
			w.Count(key)
		}(i)
	}
	wg.Wait()
}

func TestWindow_RealTime_CountDecaysOverTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time test in short mode")
	}

	w := window.New(100 * time.Millisecond)
	w.Add("k")
	w.Add("k")

	if got := w.Count("k"); got != 2 {
		t.Fatalf("want 2 before expiry, got %d", got)
	}

	time.Sleep(150 * time.Millisecond)

	if got := w.Count("k"); got != 0 {
		t.Fatalf("want 0 after expiry, got %d", got)
	}
}
