package buffer

import (
	"fmt"
	"testing"
)

func TestNew_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for capacity 0")
		}
	}()
	New[int](0)
}

func TestPush_And_Len(t *testing.T) {
	b := New[int](4)
	b.Push("a", 1)
	b.Push("b", 2)
	if b.Len() != 2 {
		t.Fatalf("want 2, got %d", b.Len())
	}
}

func TestAll_OrderIsInsertionOrder(t *testing.T) {
	b := New[int](4)
	for i := range 4 {
		b.Push(fmt.Sprintf("k%d", i), i)
	}
	entries := b.All()
	for i, e := range entries {
		if e.Value != i {
			t.Fatalf("index %d: want %d, got %d", i, i, e.Value)
		}
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	b := New[int](3)
	b.Push("a", 10)
	b.Push("b", 20)
	b.Push("c", 30)
	b.Push("d", 40) // evicts 10

	entries := b.All()
	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	if entries[0].Value != 20 {
		t.Fatalf("want oldest=20, got %d", entries[0].Value)
	}
	if entries[2].Value != 40 {
		t.Fatalf("want newest=40, got %d", entries[2].Value)
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b := New[string](4)
	b.Push("x", "hello")
	b.Push("y", "world")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("want 0 after reset, got %d", b.Len())
	}
	if got := b.All(); len(got) != 0 {
		t.Fatalf("want empty slice after reset, got %v", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	b := New[int](4)
	b.Push("a", 1)
	snap := b.All()
	snap[0].Value = 999
	if b.All()[0].Value == 999 {
		t.Fatal("All() should return a copy, not a reference")
	}
}
