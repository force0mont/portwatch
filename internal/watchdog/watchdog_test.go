package watchdog

import (
	"context"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	mu := &struct{ v time.Time }{v: t}
	return func() time.Time { return mu.v }
}

func TestBeat_ResetsTimer(t *testing.T) {
	now := time.Now()
	clock := fixedClock(now)
	w := newWithClock(100*time.Millisecond, clock)
	w.Beat()
	w.mu.Lock()
	got := w.lastBeat
	w.mu.Unlock()
	if !got.Equal(now) {
		t.Fatalf("expected lastBeat %v, got %v", now, got)
	}
}

func TestRun_FiresWhenStuck(t *testing.T) {
	base := time.Now()
	advanced := base.Add(200 * time.Millisecond)
	calls := 0
	clock := func() time.Time {
		if calls == 0 {
			calls++
			return base
		}
		return advanced
	}
	w := newWithClock(100*time.Millisecond, clock)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go w.Run(ctx)
	select {
	case <-w.Stuck():
		// expected
	case <-ctx.Done():
		t.Fatal("watchdog did not fire before deadline")
	}
}

func TestRun_DoesNotFireAfterBeat(t *testing.T) {
	w := New(200 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Run(ctx)
	// keep beating
	done := time.After(300 * time.Millisecond)
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			w.Beat()
		case <-w.Stuck():
			t.Fatal("watchdog fired despite regular beats")
		case <-done:
			return
		}
	}
}

func TestRun_CancelsCleanly(t *testing.T) {
	w := New(500 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not exit after context cancel")
	}
}
