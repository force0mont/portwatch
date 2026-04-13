package buffer_test

import (
	"sync"
	"testing"

	"github.com/example/portwatch/internal/buffer"
)

func TestConcurrent_Push_NoRace(t *testing.T) {
	b := buffer.New[int](64)
	var wg sync.WaitGroup
	workers := 8
	iterations := 200

	for w := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range iterations {
				b.Push("key", id*iterations+i)
			}
		}(w)
	}
	wg.Wait()

	if b.Len() > 64 {
		t.Fatalf("len %d exceeds capacity 64", b.Len())
	}
}

func TestConcurrent_PushAndAll_NoRace(t *testing.T) {
	b := buffer.New[string](32)
	var wg sync.WaitGroup

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				b.Push("k", "v")
			}
		}()
	}
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				_ = b.All()
			}
		}()
	}
	wg.Wait()
}
