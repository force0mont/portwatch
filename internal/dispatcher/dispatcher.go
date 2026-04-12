// Package dispatcher routes alert events to one or more registered
// notifiers, applying per-notifier error handling without blocking
// the caller.
package dispatcher

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Message is the payload sent to each notifier.
type Message struct {
	Level   string // "alert" | "info"
	Title   string
	Body    string
}

// Notifier is anything that can receive a Message.
type Notifier interface {
	Send(ctx context.Context, msg Message) error
}

// Dispatcher fans a Message out to all registered Notifiers.
type Dispatcher struct {
	mu        sync.RWMutex
	notifiers []Notifier
	logger    *log.Logger
}

// New returns a Dispatcher that writes error logs to stderr.
func New() *Dispatcher {
	return newWithWriter(os.Stderr)
}

func newWithWriter(w io.Writer) *Dispatcher {
	return &Dispatcher{
		logger: log.New(w, "dispatcher: ", 0),
	}
}

// Register adds n to the set of notifiers. Duplicate registration is
// allowed; the notifier will simply be called multiple times.
func (d *Dispatcher) Register(n Notifier) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.notifiers = append(d.notifiers, n)
}

// Dispatch sends msg to every registered notifier. Each notifier is
// called sequentially; errors are logged but do not abort delivery to
// subsequent notifiers. Dispatch returns an error only when no
// notifiers are registered.
func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) error {
	d.mu.RLock()
	ns := make([]Notifier, len(d.notifiers))
	copy(ns, d.notifiers)
	d.mu.RUnlock()

	if len(ns) == 0 {
		return fmt.Errorf("dispatcher: no notifiers registered")
	}

	for _, n := range ns {
		if err := n.Send(ctx, msg); err != nil {
			d.logger.Printf("notifier error: %v", err)
		}
	}
	return nil
}

// Len returns the number of registered notifiers.
func (d *Dispatcher) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.notifiers)
}
