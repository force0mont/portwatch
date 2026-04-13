package anomaly_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/joeshaw/portwatch/internal/anomaly"
	"github.com/joeshaw/portwatch/internal/scanner"
)

func TestConcurrent_Check_NoRace(t *testing.T) {
	det := anomaly.New(5 * time.Minute)

	ports := make([]scanner.Port, 20)
	for i := range ports {
		ports[i] = scanner.Port{
			Port:     uint16(8000 + i),
			Protocol: "tcp",
			Address:  "0.0.0.0",
		}
	}

	// Seed the detector with a known baseline.
	det.Learn(ports[:10])

	var wg sync.WaitGroup
	const goroutines = 16

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				_ = det.Check(ports)
			}
		}(g)
	}

	wg.Wait()
}

func TestConcurrent_LearnAndCheck_NoRace(t *testing.T) {
	det := anomaly.New(5 * time.Minute)

	make := func(base int) []scanner.Port {
		out := make([]scanner.Port, 5)
		for i := range out {
			out[i] = scanner.Port{
				Port:     uint16(base + i),
				Protocol: "tcp",
				Address:  fmt.Sprintf("127.0.0.%d", i+1),
			}
		}
		return out
	}

	var wg sync.WaitGroup

	for g := 0; g < 8; g++ {
		wg.Add(1)
		base := g * 100
		go func(b int) {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				ports := make(b + i)
				det.Learn(ports)
				_ = det.Check(ports)
			}
		}(base)
	}

	wg.Wait()
}

func TestAnomaly_LearnThenCheck_NoFalsePositives(t *testing.T) {
	det := anomaly.New(10 * time.Minute)

	known := []scanner.Port{
		{Port: 22, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 80, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 443, Protocol: "tcp", Address: "0.0.0.0"},
	}

	det.Learn(known)

	anomalies := det.Check(known)
	if len(anomalies) != 0 {
		t.Fatalf("expected no anomalies for learned ports, got %d", len(anomalies))
	}
}

func TestAnomaly_UnlearnedPorts_AreAnomalies(t *testing.T) {
	det := anomaly.New(10 * time.Minute)

	known := []scanner.Port{
		{Port: 22, Protocol: "tcp", Address: "0.0.0.0"},
	}
	det.Learn(known)

	unexpected := []scanner.Port{
		{Port: 22, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 9999, Protocol: "tcp", Address: "0.0.0.0"},
	}

	anomalies := det.Check(unexpected)
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Port != 9999 {
		t.Errorf("expected anomaly on port 9999, got %d", anomalies[0].Port)
	}
}
