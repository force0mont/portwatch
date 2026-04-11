package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestNew_InvalidCIDR(t *testing.T) {
	_, err := filter.New([]filter.Rule{{CIDR: "not-a-cidr"}})
	if err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestSuppressed_ByPort(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Port: 8080}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !f.Suppressed(8080, "tcp", "0.0.0.0") {
		t.Error("expected port 8080 to be suppressed")
	}
	if f.Suppressed(9090, "tcp", "0.0.0.0") {
		t.Error("port 9090 should not be suppressed")
	}
}

func TestSuppressed_ByProtocol(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Port: 53, Protocol: "udp"}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !f.Suppressed(53, "udp", "0.0.0.0") {
		t.Error("expected udp/53 to be suppressed")
	}
	if f.Suppressed(53, "tcp", "0.0.0.0") {
		t.Error("tcp/53 should not be suppressed by udp rule")
	}
}

func TestSuppressed_ByCIDR(t *testing.T) {
	f, err := filter.New([]filter.Rule{{CIDR: "127.0.0.0/8"}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !f.Suppressed(22, "tcp", "127.0.0.1") {
		t.Error("127.0.0.1 should be suppressed by 127.0.0.0/8")
	}
	if f.Suppressed(22, "tcp", "10.0.0.1") {
		t.Error("10.0.0.1 should not be suppressed")
	}
}

func TestSuppressed_CombinedRule(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Port: 443, Protocol: "tcp", CIDR: "192.168.0.0/16"}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// all three must match
	if !f.Suppressed(443, "tcp", "192.168.1.50") {
		t.Error("expected combined rule to suppress")
	}
	if f.Suppressed(443, "tcp", "10.0.0.1") {
		t.Error("address outside CIDR should not be suppressed")
	}
}

func TestLen(t *testing.T) {
	f, _ := filter.New([]filter.Rule{{Port: 80}, {Port: 443}})
	if f.Len() != 2 {
		t.Errorf("expected Len 2, got %d", f.Len())
	}
}
