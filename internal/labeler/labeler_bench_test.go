package labeler_test

import (
	"testing"

	"github.com/patrickward/portwatch/internal/labeler"
	"github.com/patrickward/portwatch/internal/scanner"
)

func BenchmarkLabel_Hit(b *testing.B) {
	l, _ := labeler.New([]labeler.Rule{
		{Port: 80, Protocol: "tcp", Label: "http"},
	})
	p := scanner.PortEntry{Protocol: "tcp", Port: 80}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Label(p)
	}
}

func BenchmarkLabel_Miss(b *testing.B) {
	l, _ := labeler.New(nil)
	p := scanner.PortEntry{Protocol: "tcp", Port: 9999}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Label(p)
	}
}
