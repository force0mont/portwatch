package quorum_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/quorum"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makePort(proto, addr string) scanner.Port {
	return scanner.Port{Protocol: proto, Address: addr}
}

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero threshold")
		}
	}()
	quorum.New(0)
}

func TestObserve_BelowThreshold_ReturnsFalse(t *testing.T) {
	q := quorum.New(3)
	p := makePort("tcp", "0.0.0.0:8080")
	if q.Observe(p) {
		t.Fatal("expected false before threshold is reached")
	}
	if q.Observe(p) {
		t.Fatal("expected false on second observation")
	}
}

func TestObserve_AtThreshold_ReturnsTrue(t *testing.T) {
	q := quorum.New(3)
	p := makePort("tcp", "0.0.0.0:8080")
	q.Observe(p)
	q.Observe(p)
	if !q.Observe(p) {
		t.Fatal("expected true at threshold")
	}
}

func TestObserve_AboveThreshold_StaysTrue(t *testing.T) {
	q := quorum.New(2)
	p := makePort("tcp", "0.0.0.0:9000")
	q.Observe(p)
	if !q.Observe(p) {
		t.Fatal("expected true at threshold")
	}
	if !q.Observe(p) {
		t.Fatal("expected true above threshold")
	}
}

func TestEvict_ResetsCounter(t *testing.T) {
	q := quorum.New(2)
	p := makePort("tcp", "0.0.0.0:443")
	q.Observe(p)
	q.Evict(p)
	if q.Count(p) != 0 {
		t.Fatalf("expected count 0 after evict, got %d", q.Count(p))
	}
	if q.Observe(p) {
		t.Fatal("expected false after evict resets counter")
	}
}

func TestObserve_IndependentPorts(t *testing.T) {
	q := quorum.New(2)
	a := makePort("tcp", "0.0.0.0:80")
	b := makePort("tcp", "0.0.0.0:443")
	q.Observe(a)
	q.Observe(a)
	if q.Observe(b) {
		t.Fatal("port b should not be confirmed after one observation")
	}
}

func TestObserve_ProtocolDistinct(t *testing.T) {
	q := quorum.New(2)
	tcp := makePort("tcp", "0.0.0.0:53")
	udp := makePort("udp", "0.0.0.0:53")
	q.Observe(tcp)
	q.Observe(tcp)
	if q.Observe(udp) {
		t.Fatal("udp port should be tracked independently from tcp")
	}
}

func TestCount_ReturnsCurrentCount(t *testing.T) {
	q := quorum.New(5)
	p := makePort("udp", "127.0.0.1:514")
	for i := 1; i <= 3; i++ {
		q.Observe(p)
		if got := q.Count(p); got != i {
			t.Fatalf("expected count %d, got %d", i, got)
		}
	}
}
