package decay

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestAdd_NewKey_SetsScore(t *testing.T) {
	tr := newWithClock(time.Minute, fixedClock(epoch))
	tr.Add("tcp:8080", 5.0)
	if got := tr.Score("tcp:8080"); got != 5.0 {
		t.Fatalf("expected 5.0, got %f", got)
	}
}

func TestScore_UnknownKey_ReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	if got := tr.Score("tcp:9999"); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestScore_AfterOneHalfLife_HalvesScore(t *testing.T) {
	halfLife := time.Hour
	tr := newWithClock(halfLife, fixedClock(epoch))
	tr.Add("tcp:80", 8.0)

	// advance clock by exactly one half-life
	tr.now = fixedClock(epoch.Add(halfLife))

	got := tr.Score("tcp:80")
	want := 4.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected ~%f, got %f", want, got)
	}
}

func TestScore_AfterTwoHalfLives_QuartersScore(t *testing.T) {
	halfLife := time.Hour
	tr := newWithClock(halfLife, fixedClock(epoch))
	tr.Add("tcp:443", 16.0)

	tr.now = fixedClock(epoch.Add(2 * halfLife))

	got := tr.Score("tcp:443")
	want := 4.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected ~%f, got %f", want, got)
	}
}

func TestAdd_AccumulatesAfterDecay(t *testing.T) {
	halfLife := time.Hour
	tr := newWithClock(halfLife, fixedClock(epoch))
	tr.Add("tcp:22", 8.0)

	// advance one half-life, score is now ~4.0
	tr.now = fixedClock(epoch.Add(halfLife))
	tr.Add("tcp:22", 2.0)

	got := tr.Score("tcp:22")
	want := 6.0
	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("expected ~%f, got %f", want, got)
	}
}

func TestReset_ClearsScore(t *testing.T) {
	tr := newWithClock(time.Minute, fixedClock(epoch))
	tr.Add("tcp:3306", 10.0)
	tr.Reset("tcp:3306")
	if got := tr.Score("tcp:3306"); got != 0 {
		t.Fatalf("expected 0 after reset, got %f", got)
	}
}

func TestAdd_IndependentKeys(t *testing.T) {
	tr := newWithClock(time.Hour, fixedClock(epoch))
	tr.Add("tcp:80", 3.0)
	tr.Add("udp:53", 7.0)

	if got := tr.Score("tcp:80"); got != 3.0 {
		t.Fatalf("tcp:80 expected 3.0, got %f", got)
	}
	if got := tr.Score("udp:53"); got != 7.0 {
		t.Fatalf("udp:53 expected 7.0, got %f", got)
	}
}
