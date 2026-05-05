package clock_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/clock"
)

func TestReal_ReturnsNonZeroTime(t *testing.T) {
	t.Parallel()
	now := clock.Real()
	if now.IsZero() {
		t.Fatal("expected non-zero time from Real clock")
	}
}

func TestFixed_AlwaysReturnsSameTime(t *testing.T) {
	t.Parallel()
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	clk := clock.Fixed(base)

	for i := 0; i < 5; i++ {
		got := clk()
		if !got.Equal(base) {
			t.Fatalf("call %d: got %v, want %v", i, got, base)
		}
	}
}

func TestAdvance_ShiftsBaseByDuration(t *testing.T) {
	t.Parallel()
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	clk := clock.Fixed(base)

	advanced := clock.Advance(clk, 5*time.Minute)
	want := base.Add(5 * time.Minute)

	if got := advanced(); !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAdvance_DoesNotMutateOriginal(t *testing.T) {
	t.Parallel()
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	clk := clock.Fixed(base)

	_ = clock.Advance(clk, 10*time.Second)

	if got := clk(); !got.Equal(base) {
		t.Fatalf("original clock mutated: got %v, want %v", got, base)
	}
}
