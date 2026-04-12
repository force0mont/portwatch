package enricher_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/scanner"
)

func TestNew_EnrichesRealPort(t *testing.T) {
	e := enricher.New()
	p := scanner.Port{
		Protocol: "tcp",
		IP:       "127.0.0.1",
		Port:     22,
	}
	ent := e.Enrich(p)

	// Service name must be populated from the built-in table.
	if ent.ServiceName != "ssh" {
		t.Fatalf("expected ssh, got %q", ent.ServiceName)
	}
	// Label must contain the protocol and port.
	if !strings.Contains(ent.Label, "tcp/22") {
		t.Fatalf("label missing tcp/22: %q", ent.Label)
	}
}

func TestNew_Concurrent_NoRace(t *testing.T) {
	e := enricher.New()
	ports := []scanner.Port{
		{Protocol: "tcp", IP: "0.0.0.0", Port: 80},
		{Protocol: "tcp", IP: "0.0.0.0", Port: 443},
		{Protocol: "udp", IP: "0.0.0.0", Port: 53},
		{Protocol: "tcp", IP: "0.0.0.0", Port: 9999},
	}

	done := make(chan struct{}, len(ports)*4)
	for i := 0; i < 4; i++ {
		for _, p := range ports {
			go func(port scanner.Port) {
				e.Enrich(port)
				done <- struct{}{}
			}(p)
		}
	}
	for i := 0; i < len(ports)*4; i++ {
		<-done
	}
}
