package gatekeeper

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedNow }

func defaultCfg() Config {
	return Config{
		MaxPerWindow:      3,
		Window:            time.Minute,
		CooldownAfterDeny: 30 * time.Second,
	}
}

func TestAdmit_UnderLimit_Admitted(t *testing.T) {
	g := newWithClock(defaultCfg(), fixedClock)
	for i := 0; i < 3; i++ {
		d := g.Admit("tcp:8080")
		if !d.Admitted {
			t.Fatalf("call %d: expected admitted, got denied (%s)", i+1, d.Reason)
		}
	}
}

func TestAdmit_ExceedsLimit_Denied(t *testing.T) {
	g := newWithClock(defaultCfg(), fixedClock)
	for i := 0; i < 3; i++ {
		g.Admit("tcp:8080")
	}
	d := g.Admit("tcp:8080")
	if d.Admitted {
		t.Fatal("expected denial after exceeding rate limit")
	}
	if d.Reason != "rate_limit" {
		t.Fatalf("expected reason rate_limit, got %s", d.Reason)
	}
}

func TestAdmit_CooldownBlocksAfterDenial(t *testing.T) {
	now := fixedNow
	g := newWithClock(defaultCfg(), func() time.Time { return now })
	for i := 0; i < 4; i++ {
		g.Admit("tcp:9090") // 4th triggers denial + sets deniedAt
	}
	// Advance by less than cooldown.
	now = now.Add(10 * time.Second)
	d := g.Admit("tcp:9090")
	if d.Admitted {
		t.Fatal("expected cooldown denial")
	}
	if d.Reason != "cooldown" {
		t.Fatalf("expected reason cooldown, got %s", d.Reason)
	}
}

func TestAdmit_AfterCooldown_AdmitsAgain(t *testing.T) {
	now := fixedNow
	g := newWithClock(defaultCfg(), func() time.Time { return now })
	for i := 0; i < 4; i++ {
		g.Admit("tcp:7070")
	}
	// Advance past cooldown AND window.
	now = now.Add(2 * time.Minute)
	d := g.Admit("tcp:7070")
	if !d.Admitted {
		t.Fatalf("expected admit after cooldown expired, got %s", d.Reason)
	}
}

func TestAdmit_IndependentKeys(t *testing.T) {
	g := newWithClock(defaultCfg(), fixedClock)
	for i := 0; i < 3; i++ {
		g.Admit("tcp:1111")
	}
	// tcp:2222 should still be admitted.
	d := g.Admit("tcp:2222")
	if !d.Admitted {
		t.Fatal("expected independent key to be admitted")
	}
}

func TestReset_ClearsState(t *testing.T) {
	g := newWithClock(defaultCfg(), fixedClock)
	for i := 0; i < 4; i++ {
		g.Admit("tcp:5555")
	}
	g.Reset("tcp:5555")
	d := g.Admit("tcp:5555")
	if !d.Admitted {
		t.Fatalf("expected admit after reset, got %s", d.Reason)
	}
}
