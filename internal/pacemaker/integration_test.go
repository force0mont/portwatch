package pacemaker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pacemaker"
)

func TestConcurrent_Beat_NoRace(t *testing.T) {
	p := pacemaker.New(100 * time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Beat()
			_ = p.Missed()
		}()
	}
	wg.Wait()
}

func TestPacemaker_RealTime_DetectsSlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time test in short mode")
	}
	p := pacemaker.New(50 * time.Millisecond)
	p.Beat()
	time.Sleep(120 * time.Millisecond)
	ok := p.Beat()
	if ok {
		t.Fatal("expected slow beat to be flagged")
	}
	if p.Missed() < 1 {
		t.Fatal("expected at least one missed beat")
	}
}
