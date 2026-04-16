package pacemaker

import (
	"testing"
	"time"
)

func BenchmarkBeat_Sequential(b *testing.B) {
	p := New(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Beat()
	}
}

func BenchmarkMissed_Sequential(b *testing.B) {
	p := New(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Missed()
	}
}
