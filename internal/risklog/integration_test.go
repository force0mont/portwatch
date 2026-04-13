package risklog

import (
	"sync"
	"testing"

	"github.com/iamcathal/portwatch/internal/scanner"
)

func TestConcurrent_Record_NoRace(t *testing.T) {
	l := New()
	ports := []scanner.Port{
		{Protocol: "tcp", Address: "0.0.0.0:80"},
		{Protocol: "tcp", Address: "0.0.0.0:443"},
		{Protocol: "udp", Address: "0.0.0.0:53"},
		{Protocol: "tcp", Address: "0.0.0.0:22"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p := ports[idx%len(ports)]
			l.Record(p, float64(idx%10)/10.0)
		}(i)
	}
	wg.Wait()

	if l.Len() == 0 {
		t.Error("expected at least one entry after concurrent records")
	}
}

func TestConcurrent_TopNAndRemove_NoRace(t *testing.T) {
	l := New()
	for i := 0; i < 20; i++ {
		l.Record(scanner.Port{
			Protocol: "tcp",
			Address:  "0.0.0.0:" + itoa(i+1000),
		}, float64(i)/20.0)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = l.TopN(5)
		}()
		go func(idx int) {
			defer wg.Done()
			l.Remove(scanner.Port{
				Protocol: "tcp",
				Address:  "0.0.0.0:" + itoa(idx+1000),
			})
		}(i)
	}
	wg.Wait()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
