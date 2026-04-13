// Package eventsink provides a thread-safe, bounded sink for port events
// that supports fan-out delivery to multiple registered consumers.
package eventsink

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alerter"
)

// ErrSinkFull is returned when the internal buffer is at capacity and a
// non-blocking push is attempted.
var ErrSinkFull = errors.New("eventsink: buffer full")

// Consumer is a function that receives a single alerter.Event for processing.
type Consumer func(alerter.Event)

// Sink fans out incoming events to all registered consumers.
type Sink struct {
	mu        sync.RWMutex
	consumers []Consumer
	ch        chan alerter.Event
	wg        sync.WaitGroup
}

// New creates a Sink with the given channel buffer capacity.
// Panics if capacity is less than 1.
func New(capacity int) *Sink {
	if capacity < 1 {
		panic("eventsink: capacity must be >= 1")
	}
	return &Sink{
		ch: make(chan alerter.Event, capacity),
	}
}

// Register adds a consumer that will receive every event pushed to the sink.
func (s *Sink) Register(c Consumer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consumers = append(s.consumers, c)
}

// Push enqueues an event. Returns ErrSinkFull if the buffer is at capacity.
func (s *Sink) Push(e alerter.Event) error {
	select {
	case s.ch <- e:
		return nil
	default:
		return ErrSinkFull
	}
}

// Run dispatches events from the internal channel to all registered consumers
// until the channel is closed. Call Close to stop.
func (s *Sink) Run() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for e := range s.ch {
			s.mu.RLock()
			cs := make([]Consumer, len(s.consumers))
			copy(cs, s.consumers)
			s.mu.RUnlock()
			for _, c := range cs {
				c(e)
			}
		}
	}()
}

// Close shuts down the sink and waits for the dispatch goroutine to finish.
func (s *Sink) Close() {
	close(s.ch)
	s.wg.Wait()
}
