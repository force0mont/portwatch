// Package probe attempts a TCP or UDP dial against a discovered port to
// verify it is genuinely accepting connections, reducing false-positive alerts
// caused by transient /proc entries.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Port     uint16
	Protocol string
	Addr     string
	Reachable bool
}

// Prober dials ports to confirm they are live.
type Prober struct {
	timeout time.Duration
	dial    func(network, addr string, timeout time.Duration) (net.Conn, error)
}

// New returns a Prober with the given dial timeout.
func New(timeout time.Duration) *Prober {
	return &Prober{
		timeout: timeout,
		dial:    net.DialTimeout,
	}
}

// newWithDialer is used in tests to inject a fake dialer.
func newWithDialer(timeout time.Duration, dial func(string, string, time.Duration) (net.Conn, error)) *Prober {
	return &Prober{timeout: timeout, dial: dial}
}

// Check dials addr:port using the given protocol ("tcp" or "udp").
// UDP probes are best-effort: a successful dial only confirms the socket
// exists locally; true reachability cannot be guaranteed without a response.
func (p *Prober) Check(protocol string, addr string, port uint16) Result {
	network := protocol
	if protocol != "tcp" && protocol != "udp" {
		return Result{Port: port, Protocol: protocol, Addr: addr, Reachable: false}
	}
	target := fmt.Sprintf("%s:%d", addr, port)
	conn, err := p.dial(network, target, p.timeout)
	if err != nil {
		return Result{Port: port, Protocol: protocol, Addr: addr, Reachable: false}
	}
	_ = conn.Close()
	return Result{Port: port, Protocol: protocol, Addr: addr, Reachable: true}
}
