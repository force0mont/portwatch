package digest

import (
	"testing"
)

func entries(proto, addr string, port uint16) Entry {
	return Entry{Proto: proto, Address: addr, Port: port}
}

func TestChanged_FirstCall_AlwaysTrue(t *testing.T) {
	d := New()
	if !d.Changed([]Entry{entries("tcp", "0.0.0.0", 80)}) {
		t.Fatal("expected first call to return true")
	}
}

func TestChanged_SameEntries_ReturnsFalse(t *testing.T) {
	d := New()
	e := []Entry{entries("tcp", "0.0.0.0", 80)}
	d.Changed(e)
	if d.Changed(e) {
		t.Fatal("expected identical snapshot to return false")
	}
}

func TestChanged_DifferentPort_ReturnsTrue(t *testing.T) {
	d := New()
	d.Changed([]Entry{entries("tcp", "0.0.0.0", 80)})
	if !d.Changed([]Entry{entries("tcp", "0.0.0.0", 443)}) {
		t.Fatal("expected changed snapshot to return true")
	}
}

func TestChanged_OrderIndependent(t *testing.T) {
	d := New()
	a := []Entry{entries("tcp", "0.0.0.0", 80), entries("udp", "0.0.0.0", 53)}
	b := []Entry{entries("udp", "0.0.0.0", 53), entries("tcp", "0.0.0.0", 80)}
	d.Changed(a)
	if d.Changed(b) {
		t.Fatal("expected order-swapped entries to produce same digest")
	}
}

func TestChanged_EmptySlice_StableDigest(t *testing.T) {
	d := New()
	d.Changed([]Entry{})
	if d.Changed([]Entry{}) {
		t.Fatal("expected two empty snapshots to be equal")
	}
}

func TestLast_BeforeAnyCall_EmptyString(t *testing.T) {
	d := New()
	if d.Last() != "" {
		t.Fatalf("expected empty string, got %q", d.Last())
	}
}

func TestLast_AfterChanged_NonEmpty(t *testing.T) {
	d := New()
	d.Changed([]Entry{entries("tcp", "127.0.0.1", 8080)})
	if d.Last() == "" {
		t.Fatal("expected non-empty digest after Changed call")
	}
}

func TestChanged_AddEntry_ReturnsTrue(t *testing.T) {
	d := New()
	one := []Entry{entries("tcp", "0.0.0.0", 22)}
	two := []Entry{entries("tcp", "0.0.0.0", 22), entries("tcp", "0.0.0.0", 80)}
	d.Changed(one)
	if !d.Changed(two) {
		t.Fatal("expected added entry to trigger change")
	}
}
