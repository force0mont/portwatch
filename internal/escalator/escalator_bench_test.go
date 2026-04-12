package escalator

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkRecord_SingleKey(b *testing.B) {
	e := New(time.Minute, 5, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Record("tcp:8080")
	}
}

func BenchmarkRecord_ManyKeys(b *testing.B) {
	e := New(time.Minute, 5, 10)
	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("tcp:%d", 1024+i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Record(keys[i%len(keys)])
	}
}
