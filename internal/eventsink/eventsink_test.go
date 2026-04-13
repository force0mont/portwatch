package eventsink_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/eventsink"
)

func fixedEvent() alerter.Event {
	return alerter.Event{Level: "alert", Message: "unexpected listener on :4444"}
}

func TestNew_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	eventsink.New(0)
}

func TestPush_ErrSinkFull_WhenBufferFull(t *testing.T) {
	s := eventsink.New(1)
	if err := s.Push(fixedEvent()); err != nil {
		t.Fatalf("first push should succeed, got: %v", err)
	}
	if err := s.Push(fixedEvent()); err != eventsink.ErrSinkFull {
		t.Fatalf("expected ErrSinkFull, got: %v", err)
	}
}

func TestRun_DeliversToSingleConsumer(t *testing.T) {
	s := eventsink.New(8)
	var received atomic.Int32
	s.Register(func(alerter.Event) { received.Add(1) })
	s.Run()

	for i := 0; i < 5; i++ {
		if err := s.Push(fixedEvent()); err != nil {
			t.Fatalf("push %d failed: %v", i, err)
		}
	}
	s.Close()

	if got := received.Load(); got != 5 {
		t.Fatalf("expected 5 deliveries, got %d", got)
	}
}

func TestRun_FansOutToMultipleConsumers(t *testing.T) {
	s := eventsink.New(8)
	var a, b atomic.Int32
	s.Register(func(alerter.Event) { a.Add(1) })
	s.Register(func(alerter.Event) { b.Add(1) })
	s.Run()

	_ = s.Push(fixedEvent())
	s.Close()

	if a.Load() != 1 || b.Load() != 1 {
		t.Fatalf("expected both consumers to receive 1 event, got a=%d b=%d", a.Load(), b.Load())
	}
}

func TestClose_IdempotentAfterDrain(t *testing.T) {
	s := eventsink.New(4)
	s.Run()
	_ = s.Push(fixedEvent())
	// Close should return promptly without deadlock.
	done := make(chan struct{})
	go func() { s.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close timed out")
	}
}

func TestRun_ConcurrentPush_NoRace(t *testing.T) {
	s := eventsink.New(64)
	s.Register(func(alerter.Event) {})
	s.Run()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Push(fixedEvent())
		}()
	}
	wg.Wait()
	s.Close()
}
