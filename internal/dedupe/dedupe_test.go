package dedupe

import (
	"testing"
	"time"

	"github.com/iamcalledned/portwatch/internal/scanner"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func fixedClock(offset time.Duration) Clock {
	current := fixedNow
	return func() time.Time {
		t := current
		current = current.Add(offset)
		return t
	}
}

func makePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Address: "0.0.0.0"}
}

func TestIsDuplicate_FirstCall_NotDuplicate(t *testing.T) {
	d := newWithClock(5*time.Second, func() time.Time { return fixedNow })
	p := makePort(8080, "tcp")
	if d.IsDuplicate(p, "alert") {
		t.Fatal("expected false on first call")
	}
}

func TestIsDuplicate_SecondCall_WithinWindow_IsDuplicate(t *testing.T) {
	clock := func() time.Time { return fixedNow } // time never advances
	d := newWithClock(5*time.Second, clock)
	p := makePort(8080, "tcp")
	d.IsDuplicate(p, "alert")
	if !d.IsDuplicate(p, "alert") {
		t.Fatal("expected true on second call within window")
	}
}

func TestIsDuplicate_AfterWindow_NotDuplicate(t *testing.T) {
	clock := fixedClock(6 * time.Second) // each call advances 6 s
	d := newWithClock(5*time.Second, clock)
	p := makePort(8080, "tcp")
	d.IsDuplicate(p, "alert") // t=0  → recorded
	if d.IsDuplicate(p, "alert") { // t=6s → evicted, should be fresh
		t.Fatal("expected false after window expired")
	}
}

func TestIsDuplicate_DifferentKinds_Independent(t *testing.T) {
	clock := func() time.Time { return fixedNow }
	d := newWithClock(5*time.Second, clock)
	p := makePort(443, "tcp")
	d.IsDuplicate(p, "alert")
	if d.IsDuplicate(p, "info") {
		t.Fatal("different kind should not be considered duplicate")
	}
}

func TestIsDuplicate_DifferentPorts_Independent(t *testing.T) {
	clock := func() time.Time { return fixedNow }
	d := newWithClock(5*time.Second, clock)
	d.IsDuplicate(makePort(80, "tcp"), "alert")
	if d.IsDuplicate(makePort(443, "tcp"), "alert") {
		t.Fatal("different port should not be considered duplicate")
	}
}

func TestReset_ClearsState(t *testing.T) {
	clock := func() time.Time { return fixedNow }
	d := newWithClock(5*time.Second, clock)
	p := makePort(8080, "tcp")
	d.IsDuplicate(p, "alert")
	d.Reset()
	if d.IsDuplicate(p, "alert") {
		t.Fatal("expected false after Reset")
	}
}
