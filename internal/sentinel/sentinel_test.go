package sentinel_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/sentinel"
)

func makePorts(specs ...struct {
	port  uint16
	proto string
	addr  string
}) []scanner.Port {
	out := make([]scanner.Port, len(specs))
	for i, s := range specs {
		out[i] = scanner.Port{Port: s.port, Protocol: s.proto, Addr: s.addr}
	}
	return out
}

func TestCheck_NoEntries_ReturnsEmpty(t *testing.T) {
	s := sentinel.New(nil)
	ports := makePorts(struct {
		port  uint16
		proto string
		addr  string
	}{4444, "tcp", "0.0.0.0"})
	if got := s.Check(ports); len(got) != 0 {
		t.Fatalf("expected no matches, got %d", len(got))
	}
}

func TestCheck_MatchingPort_ReturnsMatch(t *testing.T) {
	s := sentinel.New([]sentinel.Entry{{Port: 4444, Protocol: "tcp"}})
	ports := makePorts(struct {
		port  uint16
		proto string
		addr  string
	}{4444, "tcp", "0.0.0.0"})
	matches := s.Check(ports)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Port != 4444 || matches[0].Protocol != "tcp" {
		t.Errorf("unexpected match content: %+v", matches[0])
	}
}

func TestCheck_ProtocolMismatch_NoMatch(t *testing.T) {
	s := sentinel.New([]sentinel.Entry{{Port: 4444, Protocol: "tcp"}})
	ports := makePorts(struct {
		port  uint16
		proto string
		addr  string
	}{4444, "udp", "0.0.0.0"})
	if got := s.Check(ports); len(got) != 0 {
		t.Fatalf("expected no matches, got %d", len(got))
	}
}

func TestAdd_RegistersNewEntry(t *testing.T) {
	s := sentinel.New(nil)
	s.Add(sentinel.Entry{Port: 31337, Protocol: "tcp"})
	if s.Len() != 1 {
		t.Fatalf("expected Len 1, got %d", s.Len())
	}
}

func TestRemove_UnregistersEntry(t *testing.T) {
	e := sentinel.Entry{Port: 31337, Protocol: "tcp"}
	s := sentinel.New([]sentinel.Entry{e})
	s.Remove(e)
	if s.Len() != 0 {
		t.Fatalf("expected Len 0 after Remove, got %d", s.Len())
	}
}

func TestMatch_String_ContainsPortAndProtocol(t *testing.T) {
	m := sentinel.Match{
		Entry: sentinel.Entry{Port: 4444, Protocol: "tcp"},
		Addr:  "192.168.1.1",
	}
	got := m.String()
	for _, want := range []string{"4444", "tcp", "192.168.1.1"} {
		if len(got) == 0 {
			t.Fatalf("String() returned empty")
		}
		if !contains(got, want) {
			t.Errorf("String() %q missing %q", got, want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
