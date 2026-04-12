package tagger_test

import (
	"testing"

	"github.com/your-org/portwatch/internal/tagger"
)

func TestTag_WellKnownPort(t *testing.T) {
	tg := tagger.New()
	if got := tg.Tag(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestTag_UnknownPort(t *testing.T) {
	tg := tagger.New()
	if got := tg.Tag(9999); got != "unknown" {
		t.Fatalf("expected unknown, got %q", got)
	}
}

func TestTag_Override_TakesPrecedence(t *testing.T) {
	tg := tagger.New()
	tg.Override(80, "my-app")
	if got := tg.Tag(80); got != "my-app" {
		t.Fatalf("expected my-app, got %q", got)
	}
}

func TestTag_Override_CustomPort(t *testing.T) {
	tg := tagger.New()
	tg.Override(9200, "elasticsearch")
	if got := tg.Tag(9200); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestKnown_WellKnownPort_ReturnsTrue(t *testing.T) {
	tg := tagger.New()
	if !tg.Known(443) {
		t.Fatal("expected 443 to be known")
	}
}

func TestKnown_UnknownPort_ReturnsFalse(t *testing.T) {
	tg := tagger.New()
	if tg.Known(9999) {
		t.Fatal("expected 9999 to be unknown")
	}
}

func TestKnown_OverriddenPort_ReturnsTrue(t *testing.T) {
	tg := tagger.New()
	tg.Override(9999, "my-service")
	if !tg.Known(9999) {
		t.Fatal("expected overridden port to be known")
	}
}

func TestTag_MultipleOverrides_Independent(t *testing.T) {
	tg := tagger.New()
	tg.Override(1000, "svc-a")
	tg.Override(2000, "svc-b")

	if got := tg.Tag(1000); got != "svc-a" {
		t.Fatalf("expected svc-a, got %q", got)
	}
	if got := tg.Tag(2000); got != "svc-b" {
		t.Fatalf("expected svc-b, got %q", got)
	}
}
