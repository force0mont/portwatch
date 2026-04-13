package anomaly

import (
	"testing"
	"time"

	"github.com/jwhittle933/portwatch/internal/scanner"
)

var fixedNow = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port}
}

func TestCheck_NoKnown_AllAreAnomalies(t *testing.T) {
	d := newWithClock(nil, func() time.Time { return fixedNow })
	ports := []scanner.Port{makePort("tcp", 8080), makePort("tcp", 443)}

	anoms := d.Check(ports)

	if len(anoms) != 2 {
		t.Fatalf("expected 2 anomalies, got %d", len(anoms))
	}
}

func TestCheck_KnownPort_NotFlagged(t *testing.T) {
	known := []scanner.Port{makePort("tcp", 443)}
	d := newWithClock(known, func() time.Time { return fixedNow })

	ports := []scanner.Port{makePort("tcp", 443)}
	anoms := d.Check(ports)

	if len(anoms) != 0 {
		t.Fatalf("expected no anomalies, got %d", len(anoms))
	}
}

func TestCheck_MixedPorts_OnlyUnknownFlagged(t *testing.T) {
	known := []scanner.Port{makePort("tcp", 80)}
	d := newWithClock(known, func() time.Time { return fixedNow })

	ports := []scanner.Port{makePort("tcp", 80), makePort("tcp", 9999)}
	anoms := d.Check(ports)

	if len(anoms) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anoms))
	}
	if anoms[0].Port.Port != 9999 {
		t.Errorf("expected port 9999, got %d", anoms[0].Port.Port)
	}
}

func TestCheck_ProtocolDistinct_SamePorts(t *testing.T) {
	known := []scanner.Port{makePort("tcp", 53)}
	d := newWithClock(known, func() time.Time { return fixedNow })

	// UDP/53 is NOT in the known set even though TCP/53 is.
	ports := []scanner.Port{makePort("udp", 53)}
	anoms := d.Check(ports)

	if len(anoms) != 1 {
		t.Fatalf("expected 1 anomaly for udp:53, got %d", len(anoms))
	}
}

func TestAdd_RegistersPort_NoLongerFlagged(t *testing.T) {
	d := newWithClock(nil, func() time.Time { return fixedNow })
	p := makePort("tcp", 2222)

	d.Add(p)
	anoms := d.Check([]scanner.Port{p})

	if len(anoms) != 0 {
		t.Fatalf("expected no anomalies after Add, got %d", len(anoms))
	}
}

func TestAnomaly_String_ContainsPortAndReason(t *testing.T) {
	a := Anomaly{
		Port:       makePort("tcp", 8080),
		Reason:     "port not in known-good baseline",
		DetectedAt: fixedNow,
	}
	s := a.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
