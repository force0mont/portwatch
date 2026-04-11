package snapshot

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto, addr string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Address: addr, Port: port}
}

func TestNew_LenMatchesInput(t *testing.T) {
	ports := []scanner.Port{
		makePort("tcp", "0.0.0.0", 80),
		makePort("tcp", "0.0.0.0", 443),
	}
	s := New(ports)
	if got := s.Len(); got != 2 {
		t.Fatalf("expected len 2, got %d", got)
	}
}

func TestContains_KnownPort(t *testing.T) {
	p := makePort("tcp", "0.0.0.0", 8080)
	s := New([]scanner.Port{p})
	if !s.Contains(p) {
		t.Fatal("expected snapshot to contain port 8080")
	}
}

func TestContains_UnknownPort(t *testing.T) {
	s := New([]scanner.Port{makePort("tcp", "0.0.0.0", 80)})
	if s.Contains(makePort("tcp", "0.0.0.0", 9999)) {
		t.Fatal("expected snapshot NOT to contain port 9999")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	ports := []scanner.Port{makePort("udp", "0.0.0.0", 53)}
	s := New(ports)
	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	// Mutating the returned slice must not affect the snapshot.
	all[0].Port.Port = 9999
	if !s.Contains(makePort("udp", "0.0.0.0", 53)) {
		t.Fatal("snapshot was mutated via All() return value")
	}
}

func TestCapturedAt_IsSet(t *testing.T) {
	before := time.Now()
	s := New(nil)
	after := time.Now()
	if s.CapturedAt().Before(before) || s.CapturedAt().After(after) {
		t.Fatalf("capturedAt %v not in expected range", s.CapturedAt())
	}
}

func TestDiff_Appeared(t *testing.T) {
	old := New([]scanner.Port{makePort("tcp", "0.0.0.0", 80)})
	next := New([]scanner.Port{
		makePort("tcp", "0.0.0.0", 80),
		makePort("tcp", "0.0.0.0", 443),
	})
	appeared, disappeared := old.Diff(next)
	if len(appeared) != 1 || appeared[0].Port != 443 {
		t.Fatalf("expected port 443 appeared, got %v", appeared)
	}
	if len(disappeared) != 0 {
		t.Fatalf("expected no disappeared ports, got %v", disappeared)
	}
}

func TestDiff_Disappeared(t *testing.T) {
	old := New([]scanner.Port{
		makePort("tcp", "0.0.0.0", 80),
		makePort("tcp", "0.0.0.0", 443),
	})
	next := New([]scanner.Port{makePort("tcp", "0.0.0.0", 80)})
	_, disappeared := old.Diff(next)
	if len(disappeared) != 1 || disappeared[0].Port != 443 {
		t.Fatalf("expected port 443 disappeared, got %v", disappeared)
	}
}

func TestDiff_NoChange(t *testing.T) {
	ports := []scanner.Port{makePort("tcp", "0.0.0.0", 80)}
	old := New(ports)
	next := New(ports)
	appeared, disappeared := old.Diff(next)
	if len(appeared) != 0 || len(disappeared) != 0 {
		t.Fatalf("expected no diff, got appeared=%v disappeared=%v", appeared, disappeared)
	}
}

func TestDiff_ProtocolDistinct(t *testing.T) {
	old := New([]scanner.Port{makePort("tcp", "0.0.0.0", 53)})
	next := New([]scanner.Port{makePort("udp", "0.0.0.0", 53)})
	appeared, disappeared := old.Diff(next)
	if len(appeared) != 1 || appeared[0].Protocol != "udp" {
		t.Fatalf("expected udp:53 appeared, got %v", appeared)
	}
	if len(disappeared) != 1 || disappeared[0].Protocol != "tcp" {
		t.Fatalf("expected tcp:53 disappeared, got %v", disappeared)
	}
}
