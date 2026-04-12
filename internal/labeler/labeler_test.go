package labeler_test

import (
	"testing"

	"github.com/patrickward/portwatch/internal/labeler"
	"github.com/patrickward/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.PortEntry {
	return scanner.PortEntry{Protocol: proto, Port: port}
}

func TestLabel_MatchingRule_ReturnsLabel(t *testing.T) {
	l, err := labeler.New([]labeler.Rule{
		{Port: 80, Protocol: "tcp", Label: "http"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := l.Label(makePort("tcp", 80))
	if got != "http" {
		t.Errorf("want http, got %q", got)
	}
}

func TestLabel_NoMatch_ReturnsEmpty(t *testing.T) {
	l, _ := labeler.New(nil)
	if got := l.Label(makePort("tcp", 9999)); got != "" {
		t.Errorf("want empty, got %q", got)
	}
}

func TestLabel_ProtocolMismatch_ReturnsEmpty(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Port: 53, Protocol: "tcp", Label: "dns-tcp"},
	})
	if got := l.Label(makePort("udp", 53)); got != "" {
		t.Errorf("want empty for protocol mismatch, got %q", got)
	}
}

func TestNew_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Port: 22, Protocol: "tcp", Label: ""},
	})
	if err == nil {
		t.Error("expected error for empty label")
	}
}

func TestNew_UnknownProtocol_ReturnsError(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Port: 22, Protocol: "icmp", Label: "ping"},
	})
	if err == nil {
		t.Error("expected error for unknown protocol")
	}
}

func TestAdd_AddsRuleAtRuntime(t *testing.T) {
	l, _ := labeler.New(nil)
	if err := l.Add(labeler.Rule{Port: 443, Protocol: "tcp", Label: "https"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := l.Label(makePort("tcp", 443)); got != "https" {
		t.Errorf("want https, got %q", got)
	}
}

func TestRemove_DeletesRule(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Port: 8080, Protocol: "tcp", Label: "alt-http"},
	})
	l.Remove("tcp", 8080)
	if got := l.Label(makePort("tcp", 8080)); got != "" {
		t.Errorf("want empty after remove, got %q", got)
	}
}
