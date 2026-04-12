package enricher

import (
	"fmt"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto, ip string, port uint16) scanner.Port {
	return scanner.Port{
		Protocol: scanner.Protocol(proto),
		IP:       ip,
		Port:     port,
	}
}

func resolverOK(addr string) ([]string, error) {
	return []string{"host.example.com."}, nil
}

func resolverFail(addr string) ([]string, error) {
	return nil, fmt.Errorf("no PTR")
}

func portLookupStub(network, service string) (int, error) {
	return 0, fmt.Errorf("unused")
}

func TestEnrich_WellKnownPort_HasServiceName(t *testing.T) {
	e := newWithResolvers(resolverFail, portLookupStub)
	ent := e.Enrich(makePort("tcp", "127.0.0.1", 22))

	if ent.ServiceName != "ssh" {
		t.Fatalf("expected ssh, got %q", ent.ServiceName)
	}
	if ent.Label != "tcp/22 (ssh)" {
		t.Fatalf("unexpected label %q", ent.Label)
	}
}

func TestEnrich_UnknownPort_LabelHasNoParens(t *testing.T) {
	e := newWithResolvers(resolverFail, portLookupStub)
	ent := e.Enrich(makePort("tcp", "0.0.0.0", 9999))

	if ent.ServiceName != "" {
		t.Fatalf("expected empty service name, got %q", ent.ServiceName)
	}
	if ent.Label != "tcp/9999" {
		t.Fatalf("unexpected label %q", ent.Label)
	}
}

func TestEnrich_ResolvedHostname(t *testing.T) {
	e := newWithResolvers(resolverOK, portLookupStub)
	ent := e.Enrich(makePort("tcp", "93.184.216.34", 80))

	if ent.Hostname != "host.example.com." {
		t.Fatalf("expected hostname, got %q", ent.Hostname)
	}
}

func TestEnrich_DNSFailure_HostnameEmpty(t *testing.T) {
	e := newWithResolvers(resolverFail, portLookupStub)
	ent := e.Enrich(makePort("tcp", "10.0.0.1", 443))

	if ent.Hostname != "" {
		t.Fatalf("expected empty hostname, got %q", ent.Hostname)
	}
}

func TestEnrich_PortFieldPreserved(t *testing.T) {
	e := newWithResolvers(resolverFail, portLookupStub)
	p := makePort("udp", "0.0.0.0", 53)
	ent := e.Enrich(p)

	if ent.Port.Port != 53 || ent.Port.Protocol != "udp" {
		t.Fatalf("original port not preserved: %+v", ent.Port)
	}
}

func TestEnrich_UDP_WellKnown(t *testing.T) {
	e := newWithResolvers(resolverFail, portLookupStub)
	ent := e.Enrich(makePort("udp", "0.0.0.0", 53))

	if ent.ServiceName != "domain" {
		t.Fatalf("expected domain, got %q", ent.ServiceName)
	}
}
