package sigcache

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestMatch_EmptyCache_ReturnsFalse(t *testing.T) {
	c := newWithClock(time.Minute, fixedClock(epoch))
	if c.Match("abc") {
		t.Fatal("expected false for empty cache")
	}
}

func TestMatch_AfterSet_ReturnsTrue(t *testing.T) {
	c := newWithClock(time.Minute, fixedClock(epoch))
	c.Set("abc123")
	if !c.Match("abc123") {
		t.Fatal("expected true after Set with same sig")
	}
}

func TestMatch_DifferentSig_ReturnsFalse(t *testing.T) {
	c := newWithClock(time.Minute, fixedClock(epoch))
	c.Set("abc123")
	if c.Match("xyz999") {
		t.Fatal("expected false for different sig")
	}
}

func TestMatch_ExpiredTTL_ReturnsFalse(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	c := newWithClock(30*time.Second, clock)
	c.Set("abc123")

	now = epoch.Add(31 * time.Second)
	if c.Match("abc123") {
		t.Fatal("expected false after TTL expiry")
	}
}

func TestMatch_ZeroTTL_NeverExpires(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	c := newWithClock(0, clock)
	c.Set("abc123")

	now = epoch.Add(24 * time.Hour)
	if !c.Match("abc123") {
		t.Fatal("expected true with zero TTL regardless of elapsed time")
	}
}

func TestInvalidate_ClearsCache(t *testing.T) {
	c := newWithClock(time.Minute, fixedClock(epoch))
	c.Set("abc123")
	c.Invalidate()
	if c.Match("abc123") {
		t.Fatal("expected false after Invalidate")
	}
}

func TestSet_OverwritesPreviousSig(t *testing.T) {
	c := newWithClock(time.Minute, fixedClock(epoch))
	c.Set("first")
	c.Set("second")
	if c.Match("first") {
		t.Fatal("expected old sig to be replaced")
	}
	if !c.Match("second") {
		t.Fatal("expected new sig to match")
	}
}
