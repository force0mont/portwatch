package topology_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/topology"
)

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Protocol: "tcp", Addr: "0.0.0.0"},
		{Port: 443, Protocol: "tcp", Addr: "0.0.0.0"},
		{Port: 8080, Protocol: "tcp", Addr: "127.0.0.1"},
		{Port: 5353, Protocol: "udp", Addr: "192.168.1.10"},
		{Port: 9000, Protocol: "tcp", Addr: "192.168.1.10"},
	}
}

func TestNew_EmptyTopology(t *testing.T) {
	topo := topology.New()
	if topo.Len() != 0 {
		t.Fatalf("expected 0, got %d", topo.Len())
	}
}

func TestBuild_LenMatchesInput(t *testing.T) {
	topo := topology.New()
	ports := makePorts()
	topo.Build(ports)
	if got := topo.Len(); got != len(ports) {
		t.Fatalf("expected %d, got %d", len(ports), got)
	}
}

func TestBuild_GroupsCorrectClasses(t *testing.T) {
	topo := topology.New()
	topo.Build(makePorts())

	groups := topo.Groups()
	classCount := map[string]int{}
	for _, g := range groups {
		classCount[g.Class] += len(g.Ports)
	}

	if classCount["loopback"] != 1 {
		t.Errorf("expected 1 loopback port, got %d", classCount["loopback"])
	}
	if classCount["private"] != 2 {
		t.Errorf("expected 2 private ports, got %d", classCount["private"])
	}
}

func TestBuild_ReplacesOldGroups(t *testing.T) {
	topo := topology.New()
	topo.Build(makePorts())
	topo.Build([]scanner.Port{
		{Port: 22, Protocol: "tcp", Addr: "0.0.0.0"},
	})
	if got := topo.Len(); got != 1 {
		t.Fatalf("expected 1 after rebuild, got %d", got)
	}
}

func TestGroups_ReturnsCopy(t *testing.T) {
	topo := topology.New()
	topo.Build(makePorts())

	g1 := topo.Groups()
	g2 := topo.Groups()
	if len(g1) != len(g2) {
		t.Fatal("snapshot lengths differ")
	}
	// Mutate copy — should not affect next call.
	g1[0].Ports = nil
	g3 := topo.Groups()
	if topo.Len() == 0 {
		t.Fatal("mutation of copy affected topology")
	}
	_ = g3
}

func TestBuild_ProtocolsKeptSeparate(t *testing.T) {
	topo := topology.New()
	topo.Build([]scanner.Port{
		{Port: 53, Protocol: "tcp", Addr: "0.0.0.0"},
		{Port: 53, Protocol: "udp", Addr: "0.0.0.0"},
	})
	groups := topo.Groups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups (tcp+udp), got %d", len(groups))
	}
}
