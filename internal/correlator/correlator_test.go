package correlator

import (
	"testing"
	"time"

	"github.com/iamcalledrob/portwatch/internal/scanner"
)

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedAt(t time.Time) func() time.Time { return func() time.Time { return t } }

func makePort(port uint16, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Addr: "127.0.0.1"}
}

func TestRecord_CreatesNewGroup(t *testing.T) {
	c := newWithClock(5*time.Second, fixedAt(t0))
	g := c.Record("scan", makePort(80, "tcp"))
	if g == nil {
		t.Fatal("expected non-nil group")
	}
	if len(g.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(g.Ports))
	}
}

func TestRecord_SameKey_WithinWindow_AppendsPort(t *testing.T) {
	c := newWithClock(5*time.Second, fixedAt(t0))
	c.Record("scan", makePort(80, "tcp"))
	g := c.Record("scan", makePort(443, "tcp"))
	if len(g.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(g.Ports))
	}
}

func TestRecord_SameKey_AfterWindow_StartsNewGroup(t *testing.T) {
	now := t0
	c := newWithClock(5*time.Second, func() time.Time { return now })
	c.Record("scan", makePort(80, "tcp"))
	now = t0.Add(10 * time.Second)
	g := c.Record("scan", makePort(443, "tcp"))
	if len(g.Ports) != 1 {
		t.Fatalf("expected new group with 1 port, got %d", len(g.Ports))
	}
}

func TestRecord_DifferentKeys_IndependentGroups(t *testing.T) {
	c := newWithClock(5*time.Second, fixedAt(t0))
	c.Record("alpha", makePort(80, "tcp"))
	c.Record("beta", makePort(22, "tcp"))
	if c.Len() != 2 {
		t.Fatalf("expected 2 groups, got %d", c.Len())
	}
}

func TestFlush_ReturnsExpiredGroups(t *testing.T) {
	now := t0
	c := newWithClock(5*time.Second, func() time.Time { return now })
	c.Record("scan", makePort(80, "tcp"))
	now = t0.Add(10 * time.Second)
	expired := c.Flush()
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired group, got %d", len(expired))
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 active groups after flush, got %d", c.Len())
	}
}

func TestFlush_RetainsActiveGroups(t *testing.T) {
	c := newWithClock(5*time.Second, fixedAt(t0))
	c.Record("scan", makePort(80, "tcp"))
	expired := c.Flush()
	if len(expired) != 0 {
		t.Fatalf("expected 0 expired groups, got %d", len(expired))
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 active group, got %d", c.Len())
	}
}
