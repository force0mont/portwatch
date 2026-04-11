// Package watcher ties together the scanner, rules engine, state tracker,
// alerter, filter, throttle, and notifier into a single polling loop.
package watcher

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/throttle"
)

// Watcher polls open ports and emits alerts according to configured rules.
type Watcher struct {
	scanner  *scanner.Scanner
	rules    *rules.Engine
	state    *state.State
	alerter  *alerter.Alerter
	filter   *filter.Filter
	throttle *throttle.Throttle
	interval time.Duration
}

// New constructs a Watcher with all dependencies.
func New(
	s *scanner.Scanner,
	e *rules.Engine,
	st *state.State,
	a *alerter.Alerter,
	f *filter.Filter,
	th *throttle.Throttle,
	interval time.Duration,
) *Watcher {
	return &Watcher{
		scanner:  s,
		rules:    e,
		state:    st,
		alerter:  a,
		filter:   f,
		throttle: th,
		interval: interval,
	}
}

// Run starts the polling loop; it returns when ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.tick()
		}
	}
}

func (w *Watcher) tick() {
	ports, err := w.scanner.Scan()
	if err != nil {
		w.alerter.EmitInfo("scan error: " + err.Error())
		return
	}

	events := w.state.Diff(ports)
	for _, ev := range events {
		if w.filter.Suppressed(ev.Port, ev.Proto, "") {
			continue
		}
		action := w.rules.Evaluate(ev.Port, ev.Proto)
		if action == rules.ActionAlert {
			if w.throttle.Allow(ev.Port, ev.Proto) {
				w.alerter.EmitAlert(ev)
			}
		} else {
			w.alerter.EmitInfo(ev)
		}
	}
}
