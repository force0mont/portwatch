package probe

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func dialSuccess(_ string, _ string, _ time.Duration) (net.Conn, error) {
	c1, c2 := net.Pipe()
	_ = c2.Close()
	return c1, nil
}

func dialFailure(_ string, _ string, _ time.Duration) (net.Conn, error) {
	return nil, fmt.Errorf("connection refused")
}

func TestCheck_TCP_Reachable(t *testing.T) {
	p := newWithDialer(time.Second, dialSuccess)
	r := p.Check("tcp", "127.0.0.1", 8080)
	if !r.Reachable {
		t.Fatal("expected reachable=true")
	}
	if r.Port != 8080 || r.Protocol != "tcp" {
		t.Fatalf("unexpected result fields: %+v", r)
	}
}

func TestCheck_TCP_Unreachable(t *testing.T) {
	p := newWithDialer(time.Second, dialFailure)
	r := p.Check("tcp", "127.0.0.1", 9999)
	if r.Reachable {
		t.Fatal("expected reachable=false")
	}
}

func TestCheck_UDP_Reachable(t *testing.T) {
	p := newWithDialer(time.Second, dialSuccess)
	r := p.Check("udp", "0.0.0.0", 53)
	if !r.Reachable {
		t.Fatal("expected reachable=true for udp")
	}
}

func TestCheck_UnknownProtocol_NotReachable(t *testing.T) {
	p := newWithDialer(time.Second, dialSuccess)
	r := p.Check("sctp", "127.0.0.1", 1234)
	if r.Reachable {
		t.Fatal("unknown protocol should return reachable=false")
	}
}

func TestNew_NotNil(t *testing.T) {
	p := New(500 * time.Millisecond)
	if p == nil {
		t.Fatal("expected non-nil Prober")
	}
}
