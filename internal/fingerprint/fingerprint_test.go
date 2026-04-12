package fingerprint

import (
	"net"
	"testing"
)

func makeEntry(proto string, ip string, port uint16) Entry {
	return Entry{
		Proto: proto,
		Addr:  net.ParseIP(ip),
		Port:  port,
	}
}

func TestDerive_SameEntry_SameFingerprint(t *testing.T) {
	d := New()
	e := makeEntry("tcp", "0.0.0.0", 8080)

	fp1 := d.Derive(e)
	fp2 := d.Derive(e)

	if fp1 != fp2 {
		t.Fatalf("expected same fingerprint, got %q vs %q", fp1, fp2)
	}
}

func TestDerive_DifferentPort_DifferentFingerprint(t *testing.T) {
	d := New()
	e1 := makeEntry("tcp", "0.0.0.0", 8080)
	e2 := makeEntry("tcp", "0.0.0.0", 9090)

	if d.Derive(e1) == d.Derive(e2) {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestDerive_DifferentProtocol_DifferentFingerprint(t *testing.T) {
	d := New()
	e1 := makeEntry("tcp", "0.0.0.0", 53)
	e2 := makeEntry("udp", "0.0.0.0", 53)

	if d.Derive(e1) == d.Derive(e2) {
		t.Fatal("expected different fingerprints for different protocols")
	}
}

func TestDerive_CachesResult(t *testing.T) {
	d := New()
	e := makeEntry("tcp", "127.0.0.1", 443)

	d.Derive(e)

	if got := d.Len(); got != 1 {
		t.Fatalf("expected cache len 1, got %d", got)
	}
}

func TestInvalidate_RemovesFromCache(t *testing.T) {
	d := New()
	e := makeEntry("tcp", "127.0.0.1", 22)

	d.Derive(e)
	if d.Len() != 1 {
		t.Fatal("expected entry in cache before invalidation")
	}

	d.Invalidate(e)
	if d.Len() != 0 {
		t.Fatal("expected empty cache after invalidation")
	}
}

func TestDerive_FingerprintLength(t *testing.T) {
	d := New()
	e := makeEntry("udp", "0.0.0.0", 161)
	fp := d.Derive(e)

	// sha256 first 8 bytes → 16 hex chars
	if len(fp) != 16 {
		t.Fatalf("expected fingerprint length 16, got %d", len(fp))
	}
}

func TestDerive_IndependentInstances_SameResult(t *testing.T) {
	e := makeEntry("tcp", "10.0.0.1", 3306)

	fp1 := New().Derive(e)
	fp2 := New().Derive(e)

	if fp1 != fp2 {
		t.Fatalf("fingerprint should be deterministic across instances: %q vs %q", fp1, fp2)
	}
}
