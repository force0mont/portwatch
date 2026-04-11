package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
	"github.com/user/portwatch/internal/metrics"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func TestIntegration_HealthServer_LiveRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	m := metrics.New()
	m.IncScans()
	m.IncAlerts()
	m.AddPorts(3)

	addr := freePort(t)
	s := healthcheck.New(addr, m)

	go func() { _ = s.ListenAndServe() }()
	defer s.Shutdown() //nolint:errcheck

	// Wait for the server to be ready.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("status = %v, want ok", body["status"])
	}

	metricsBlock, ok := body["metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("metrics field missing or wrong type")
	}
	if metricsBlock["total_scans"].(float64) < 1 {
		t.Error("expected total_scans >= 1")
	}
}
