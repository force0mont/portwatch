package ratelimit

import (
	"testing"
	"time"
)

var (
	testWindow = 10 * time.Second
	testBurst  = 3
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_UnderBurst(t *testing.T) {
	now := time.Now()
	l := newWithClock(testWindow, testBurst, fixedClock(now))

	for i := 0; i < testBurst; i++ {
		if !l.Allow("tcp:8080") {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	now := time.Now()
	l := newWithClock(testWindow, testBurst, fixedClock(now))

	for i := 0; i < testBurst; i++ {
		l.Allow("tcp:8080")
	}
	if l.Allow("tcp:8080") {
		t.Fatal("expected Allow=false after burst exceeded")
	}
}

func TestAllow_WindowResets(t *testing.T) {
	now := time.Now()
	clock := fixedClock(now)
	l := newWithClock(testWindow, testBurst, clock)

	for i := 0; i < testBurst; i++ {
		l.Allow("tcp:9090")
	}
	if l.Allow("tcp:9090") {
		t.Fatal("expected suppression before window reset")
	}

	// Advance clock past the window.
	l.now = fixedClock(now.Add(testWindow + time.Millisecond))
	if !l.Allow("tcp:9090") {
		t.Fatal("expected Allow=true after window reset")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	now := time.Now()
	l := newWithClock(testWindow, 1, fixedClock(now))

	if !l.Allow("tcp:80") {
		t.Fatal("expected Allow=true for tcp:80")
	}
	if !l.Allow("udp:53") {
		t.Fatal("expected Allow=true for udp:53 (independent key)")
	}
	if l.Allow("tcp:80") {
		t.Fatal("expected Allow=false for tcp:80 after burst=1")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	now := time.Now()
	l := newWithClock(testWindow, 1, fixedClock(now))

	l.Allow("tcp:443")
	if l.Allow("tcp:443") {
		t.Fatal("expected suppression before reset")
	}
	l.Reset("tcp:443")
	if !l.Allow("tcp:443") {
		t.Fatal("expected Allow=true after Reset")
	}
}

func TestPurge_RemovesExpiredBuckets(t *testing.T) {
	now := time.Now()
	l := newWithClock(testWindow, testBurst, fixedClock(now))

	l.Allow("tcp:22")
	l.now = fixedClock(now.Add(testWindow + time.Millisecond))
	l.Purge()

	l.mu.Lock()
	_, exists := l.buckets["tcp:22"]
	l.mu.Unlock()

	if exists {
		t.Fatal("expected bucket to be purged after window expiry")
	}
}
