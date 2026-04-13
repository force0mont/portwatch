// Package buffer provides a fixed-capacity ring buffer that retains the
// most recent N port events, evicting the oldest entry when full.
package buffer

import "sync"

// Entry holds a single buffered item with an associated key.
type Entry[T any] struct {
	Key   string
	Value T
}

// Buffer is a thread-safe, fixed-capacity ring buffer.
type Buffer[T any] struct {
	mu       sync.Mutex
	slots    []Entry[T]
	head     int
	count    int
	capacity int
}

// New returns a Buffer with the given capacity. Panics if capacity < 1.
func New[T any](capacity int) *Buffer[T] {
	if capacity < 1 {
		panic("buffer: capacity must be >= 1")
	}
	return &Buffer[T]{
		slots:    make([]Entry[T], capacity),
		capacity: capacity,
	}
}

// Push adds a new entry, evicting the oldest if the buffer is full.
func (b *Buffer[T]) Push(key string, value T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := (b.head + b.count) % b.capacity
	b.slots[idx] = Entry[T]{Key: key, Value: value}

	if b.count < b.capacity {
		b.count++
	} else {
		// overwrite oldest: advance head
		b.head = (b.head + 1) % b.capacity
	}
}

// All returns a copy of all buffered entries in insertion order (oldest first).
func (b *Buffer[T]) All() []Entry[T] {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]Entry[T], b.count)
	for i := range b.count {
		out[i] = b.slots[(b.head+i)%b.capacity]
	}
	return out
}

// Len returns the current number of entries in the buffer.
func (b *Buffer[T]) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset clears all entries without releasing the underlying memory.
func (b *Buffer[T]) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.count = 0
}
