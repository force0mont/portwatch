package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// TestThrottle_MultiplePortsUnderLoad verifies independent counters for
// several ports running concurrently within the same window.
func TestThrottle_MultiplePortsUnderLoad(t *testing.T) {
	th := throttle.New(3, time.Minute)

	ports := []uint16{80, 443, 8080, 9090}
	for _, p := range ports {
		for i := 0; i < 3; i++ {
			if !th.Allow(p, "tcp") {
				t.Errorf("port %d call %d: expected Allow=true", p, i+1)
			}
		}
		// 4th call should be throttled
		if th.Allow(p, "tcp") {
			t.Errorf("port %d: expected Allow=false on 4th call", p)
		}
	}
}

// TestThrottle_ResetRestoresAllPorts verifies Reset clears all keys.
func TestThrottle_ResetRestoresAllPorts(t *testing.T) {
	th := throttle.New(1, time.Minute)

	ports := []uint16{22, 25, 110}
	for _, p := range ports {
		th.Allow(p, "tcp")
		th.Allow(p, "tcp") // throttle each
	}

	th.Reset()

	for _, p := range ports {
		if !th.Allow(p, "tcp") {
			t.Errorf("port %d: expected Allow=true after Reset", p)
		}
	}
}
