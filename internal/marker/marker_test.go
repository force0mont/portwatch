package marker

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestAck_And_IsAcked(t *testing.T) {
	m := New()
	m.Ack(8080, "tcp", "alice", 0, epoch)
	if !m.IsAcked(8080, "tcp", epoch) {
		t.Fatal("expected port to be acked")
	}
}

func TestIsAcked_UnknownPort_ReturnsFalse(t *testing.T) {
	m := New()
	if m.IsAcked(9999, "tcp", epoch) {
		t.Fatal("expected unacked port to return false")
	}
}

func TestIsAcked_ExpiredTTL_ReturnsFalse(t *testing.T) {
	m := New()
	m.Ack(443, "tcp", "bob", time.Minute, epoch)
	later := epoch.Add(2 * time.Minute)
	if m.IsAcked(443, "tcp", later) {
		t.Fatal("expected expired ack to return false")
	}
}

func TestIsAcked_WithinTTL_ReturnsTrue(t *testing.T) {
	m := New()
	m.Ack(443, "tcp", "bob", time.Hour, epoch)
	soon := epoch.Add(30 * time.Minute)
	if !m.IsAcked(443, "tcp", soon) {
		t.Fatal("expected ack within TTL to return true")
	}
}

func TestRevoke_RemovesEntry(t *testing.T) {
	m := New()
	m.Ack(22, "tcp", "carol", 0, epoch)
	m.Revoke(22, "tcp")
	if m.IsAcked(22, "tcp", epoch) {
		t.Fatal("expected revoked port to return false")
	}
}

func TestProtocol_Distinct(t *testing.T) {
	m := New()
	m.Ack(53, "tcp", "dave", 0, epoch)
	if m.IsAcked(53, "udp", epoch) {
		t.Fatal("tcp ack should not satisfy udp query")
	}
}

func TestAll_ReturnsNonExpired(t *testing.T) {
	m := New()
	m.Ack(80, "tcp", "eve", 0, epoch)               // no expiry
	m.Ack(8080, "tcp", "eve", time.Minute, epoch)    // expires after 1m
	m.Ack(9090, "tcp", "eve", time.Second, epoch)    // already expired

	later := epoch.Add(90 * time.Second)
	entries := m.All(later)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	m := New()
	m.Ack(80, "tcp", "frank", 0, epoch)
	a := m.All(epoch)
	a[0].AckedBy = "tampered"
	b := m.All(epoch)
	if b[0].AckedBy == "tampered" {
		t.Fatal("All should return an independent copy")
	}
}
