package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/jwhittle933/portwatch/internal/alerter"
	"github.com/jwhittle933/portwatch/internal/metrics"
	"github.com/jwhittle933/portwatch/internal/ratelimit"
	"github.com/jwhittle933/portwatch/internal/rules"
	"github.com/jwhittle933/portwatch/internal/scanner"
	"github.com/jwhittle933/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning, rule evaluation, and alerting.
type Watcher struct {
	scanner  *scanner.Scanner
	rules    *rules.Engine
	alerter  *alerter.Alerter
	state    *state.State
	metrics  *metrics.Metrics
	limiter  *ratelimit.Limiter
	interval time.Duration
}

// New creates a Watcher with default rate-limit settings (burst=3, window=60s).
func New(
	sc *scanner.Scanner,
	re *rules.Engine,
	al *alerter.Alerter,
	st *state.State,
	me *metrics.Metrics,
	interval time.Duration,
) *Watcher {
	return &Watcher{
		scanner:  sc,
		rules:    re,
		alerter:  al,
		state:    st,
		metrics:  me,
		limiter:  ratelimit.New(60*time.Second, 3),
		interval: interval,
	}
}

// Run starts the watch loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				return fmt.Errorf("watcher tick: %w", err)
			}
		}
	}
}

func (w *Watcher) tick() error {
	ports, err := w.scanner.Scan()
	if err != nil {
		return err
	}
	w.metrics.IncScans()
	w.metrics.AddPorts(len(ports))
	w.limiter.Purge()

	events := w.state.Diff(ports)
	for _, ev := range events {
		action := w.rules.Evaluate(ev.Port)
		key := fmt.Sprintf("%s:%d", ev.Port.Protocol, ev.Port.Port)
		if action == rules.ActionAlert && w.limiter.Allow(key) {
			w.metrics.IncAlerts()
			w.alerter.EmitAlert(ev)
		} else {
			w.alerter.EmitInfo(ev)
		}
	}
	return nil
}
