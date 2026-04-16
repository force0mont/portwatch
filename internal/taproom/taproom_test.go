package taproom

import (
	"testing"
)

func TestNew_PanicsOnZeroMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	New(0)
}

func TestAdd_And_Contains(t *testing.T) {
	tr := New(10)
	if err := tr.Add(8080, "tcp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.Contains(8080, "tcp") {
		t.Fatal("expected port to be contained")
	}
}

func TestAdd_ErrAlreadyTapped(t *testing.T) {
	tr := New(10)
	_ = tr.Add(443, "tcp")
	if err := tr.Add(443, "tcp"); err != ErrAlreadyTapped {
		t.Fatalf("expected ErrAlreadyTapped, got %v", err)
	}
}

func TestAdd_ErrTapFull(t *testing.T) {
	tr := New(2)
	_ = tr.Add(80, "tcp")
	_ = tr.Add(443, "tcp")
	if err := tr.Add(8080, "tcp"); err != ErrTapFull {
		t.Fatalf("expected ErrTapFull, got %v", err)
	}
}

func TestRemove_KnownPort_ReturnsTrue(t *testing.T) {
	tr := New(10)
	_ = tr.Add(22, "tcp")
	if !tr.Remove(22, "tcp") {
		t.Fatal("expected Remove to return true")
	}
	if tr.Contains(22, "tcp") {
		t.Fatal("expected port to be removed")
	}
}

func TestRemove_UnknownPort_ReturnsFalse(t *testing.T) {
	tr := New(10)
	if tr.Remove(9999, "tcp") {
		t.Fatal("expected Remove to return false for unknown port")
	}
}

func TestProtocol_Independence(t *testing.T) {
	tr := New(10)
	_ = tr.Add(53, "tcp")
	_ = tr.Add(53, "udp")
	if tr.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", tr.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := New(10)
	_ = tr.Add(80, "tcp")
	_ = tr.Add(443, "tcp")
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
	// Mutating the copy should not affect internal state.
	all[0] = Entry{Port: 9999, Protocol: "udp"}
	if !tr.Contains(80, "tcp") {
		t.Fatal("original entry should still be present")
	}
}
