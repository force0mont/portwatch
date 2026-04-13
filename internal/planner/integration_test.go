package planner_test

import (
	"sync"
	"testing"
	"time"

	"github.com/iamcaleberic/portwatch/internal/planner"
)

func TestConcurrent_MarkAndNext_NoRace(t *testing.T) {
	p := planner.New(100*time.Millisecond, 10*time.Millisecond)
	var wg sync.WaitGroup
	now := time.Now()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			p.Mark(now.Add(time.Duration(offset) * 100 * time.Millisecond))
			_ = p.Next()
			_ = p.Missed()
		}(i)
	}
	wg.Wait()
}

func TestPlanner_MissedAccumulates(t *testing.T) {
	p := planner.New(10*time.Second, 0)
	base := time.Now()

	p.Mark(base)
	p.Mark(base.Add(20 * time.Second)) // missed
	p.Mark(base.Add(40 * time.Second)) // missed
	p.Mark(base.Add(52 * time.Second)) // within 1.5× → not missed

	if got := p.Missed(); got != 2 {
		t.Fatalf("expected 2 missed scans, got %d", got)
	}
}
