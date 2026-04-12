package labeler_test

import (
	"sync"
	"testing"

	"github.com/patrickward/portwatch/internal/labeler"
	"github.com/patrickward/portwatch/internal/scanner"
)

func TestConcurrent_LabelAndAdd_NoRace(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Port: 22, Protocol: "tcp", Label: "ssh"},
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			port := uint16(1000 + i)
			_ = l.Add(labeler.Rule{Port: port, Protocol: "tcp", Label: "custom"})
			_ = l.Label(scanner.PortEntry{Protocol: "tcp", Port: port})
		}(i)
	}
	wg.Wait()
}

func TestConcurrent_RemoveAndLabel_NoRace(t *testing.T) {
	rules := make([]labeler.Rule, 20)
	for i := range rules {
		rules[i] = labeler.Rule{Port: uint16(2000 + i), Protocol: "udp", Label: "svc"}
	}
	l, _ := labeler.New(rules)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			l.Remove("udp", uint16(2000+i))
		}(i)
		go func(i int) {
			defer wg.Done()
			_ = l.Label(scanner.PortEntry{Protocol: "udp", Port: uint16(2000 + i)})
		}(i)
	}
	wg.Wait()
}
