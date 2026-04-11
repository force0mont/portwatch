package state_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func port(proto, addr string, port uint16, pid int) scanner.Port {
	return scanner.Port{Proto: proto, Addr: addr, Port: port, PID: pid}
}

func TestDiff_FirstScan_AllAppeared(t *testing.T) {
	tr := state.New()
	ports := []scanner.Port{port("tcp", "0.0.0.0", 80, 100)}
	changes := tr.Diff(ports)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != state.Appeared {
		t.Errorf("expected Appeared, got %s", changes[0].Kind)
	}
}

func TestDiff_NoChange_NoEvents(t *testing.T) {
	tr := state.New()
	ports := []scanner.Port{port("tcp", "0.0.0.0", 443, 200)}
	tr.Diff(ports)
	changes := tr.Diff(ports)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_PortDisappears(t *testing.T) {
	tr := state.New()
	ports := []scanner.Port{port("tcp", "0.0.0.0", 8080, 300)}
	tr.Diff(ports)
	changes := tr.Diff([]scanner.Port{})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != state.Disappeared {
		t.Errorf("expected Disappeared, got %s", changes[0].Kind)
	}
}

func TestDiff_NewAndGone(t *testing.T) {
	tr := state.New()
	old := []scanner.Port{port("tcp", "0.0.0.0", 22, 1)}
	tr.Diff(old)
	new := []scanner.Port{port("udp", "0.0.0.0", 53, 2)}
	changes := tr.Diff(new)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
}

func TestSnapshot_ReturnsCurrentBaseline(t *testing.T) {
	tr := state.New()
	ports := []scanner.Port{
		port("tcp", "127.0.0.1", 3306, 500),
		port("tcp", "0.0.0.0", 5432, 501),
	}
	tr.Diff(ports)
	snap := tr.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected snapshot length 2, got %d", len(snap))
	}
}

func TestNew_EmptyBaseline(t *testing.T) {
	tr := state.New()
	if snap := tr.Snapshot(); len(snap) != 0 {
		t.Errorf("expected empty baseline, got %d entries", len(snap))
	}
}
