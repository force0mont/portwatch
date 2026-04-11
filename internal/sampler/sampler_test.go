package sampler

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func entry(addr string, port int, proto string) scanner.Entry {
	return scanner.Entry{Addr: addr, Port: port, Protocol: proto}
}

func TestSample_N1_AlwaysPasses(t *testing.T) {
	s := New(1)
	e := entry("0.0.0.0", 8080, "tcp")
	for i := 0; i < 5; i++ {
		if !s.Sample(e) {
			t.Fatalf("call %d: expected true with n=1", i+1)
		}
	}
}

func TestSample_N3_PassesEveryThird(t *testing.T) {
	s := New(3)
	e := entry("0.0.0.0", 9000, "tcp")
	expected := []bool{true, false, false, true, false, false}
	for i, want := range expected {
		got := s.Sample(e)
		if got != want {
			t.Errorf("call %d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestSample_IndependentKeys(t *testing.T) {
	s := New(2)
	a := entry("0.0.0.0", 80, "tcp")
	b := entry("0.0.0.0", 443, "tcp")
	// first call for each key should pass
	if !s.Sample(a) {
		t.Error("first call for port 80 should pass")
	}
	if !s.Sample(b) {
		t.Error("first call for port 443 should pass")
	}
	// second call for each key should be suppressed
	if s.Sample(a) {
		t.Error("second call for port 80 should be suppressed")
	}
	if s.Sample(b) {
		t.Error("second call for port 443 should be suppressed")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s := New(3)
	e := entry("127.0.0.1", 22, "tcp")
	s.Sample(e) // count=1 → passes
	s.Sample(e) // count=2
	s.Reset()
	if !s.Sample(e) { // count should restart at 1 → passes
		t.Error("expected true after reset")
	}
}

func TestLen_TracksDistinctKeys(t *testing.T) {
	s := New(1)
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Sample(entry("0.0.0.0", 80, "tcp"))
	s.Sample(entry("0.0.0.0", 443, "tcp"))
	s.Sample(entry("0.0.0.0", 80, "udp"))
	if s.Len() != 3 {
		t.Fatalf("expected 3, got %d", s.Len())
	}
}

func TestNew_LessThanOne_TreatedAsOne(t *testing.T) {
	s := New(0)
	e := entry("0.0.0.0", 8080, "tcp")
	for i := 0; i < 4; i++ {
		if !s.Sample(e) {
			t.Errorf("call %d: expected true when n clamped to 1", i+1)
		}
	}
}
