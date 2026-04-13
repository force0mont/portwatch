package observer

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return epoch }

func makePort(proto, addr string) scanner.Port {
	return scanner.Port{Protocol: proto, Address: addr}
}

func TestObserve_NewPort_CreatesEntry(t *testing.T) {
	o := newWithClock(fixedClock)
	o.Observe([]scanner.Port{makePort("tcp", "0.0.0.0:80")})

	e, ok := o.Get(makePort("tcp", "0.0.0.0:80"))
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.SeenCount != 1 {
		t.Errorf("SeenCount = %d, want 1", e.SeenCount)
	}
	if e.MissCount != 0 {
		t.Errorf("MissCount = %d, want 0", e.MissCount)
	}
}

func TestObserve_SecondScan_IncrementsSeenCount(t *testing.T) {
	o := newWithClock(fixedClock)
	p := makePort("tcp", "0.0.0.0:443")
	o.Observe([]scanner.Port{p})
	o.Observe([]scanner.Port{p})

	e, _ := o.Get(p)
	if e.SeenCount != 2 {
		t.Errorf("SeenCount = %d, want 2", e.SeenCount)
	}
}

func TestObserve_MissingPort_IncrementsMissCount(t *testing.T) {
	o := newWithClock(fixedClock)
	p := makePort("tcp", "0.0.0.0:8080")
	o.Observe([]scanner.Port{p})
	o.Observe([]scanner.Port{}) // port absent

	e, _ := o.Get(p)
	if e.MissCount != 1 {
		t.Errorf("MissCount = %d, want 1", e.MissCount)
	}
}

func TestStabilityScore_AllPresent_ReturnsOne(t *testing.T) {
	o := newWithClock(fixedClock)
	p := makePort("udp", "0.0.0.0:53")
	o.Observe([]scanner.Port{p})
	o.Observe([]scanner.Port{p})
	o.Observe([]scanner.Port{p})

	e, _ := o.Get(p)
	if got := e.StabilityScore(); got != 1.0 {
		t.Errorf("StabilityScore = %.2f, want 1.00", got)
	}
}

func TestStabilityScore_HalfMissed_ReturnsHalf(t *testing.T) {
	o := newWithClock(fixedClock)
	p := makePort("tcp", "0.0.0.0:22")
	o.Observe([]scanner.Port{p})
	o.Observe([]scanner.Port{})

	e, _ := o.Get(p)
	if got := e.StabilityScore(); got != 0.5 {
		t.Errorf("StabilityScore = %.2f, want 0.50", got)
	}
}

func TestGet_UnknownPort_ReturnsFalse(t *testing.T) {
	o := newWithClock(fixedClock)
	_, ok := o.Get(makePort("tcp", "0.0.0.0:9999"))
	if ok {
		t.Error("expected ok=false for unknown port")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	o := newWithClock(fixedClock)
	o.Observe([]scanner.Port{makePort("tcp", "0.0.0.0:80"), makePort("udp", "0.0.0.0:53")})

	all := o.All()
	if len(all) != 2 {
		t.Errorf("len(All) = %d, want 2", len(all))
	}
}

func TestLen_MatchesUniquePortCount(t *testing.T) {
	o := newWithClock(fixedClock)
	o.Observe([]scanner.Port{makePort("tcp", "0.0.0.0:80")})
	o.Observe([]scanner.Port{makePort("tcp", "0.0.0.0:80"), makePort("udp", "0.0.0.0:53")})

	if got := o.Len(); got != 2 {
		t.Errorf("Len = %d, want 2", got)
	}
}
