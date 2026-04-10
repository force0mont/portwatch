// Package watcher ties together the scanner, rules engine, and alerter
// into a polling loop that monitors open ports at a configurable interval.
package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
)

// Watcher polls the system for open ports and evaluates each one against
// the configured rule set, emitting alerts via the Alerter.
type Watcher struct {
	scanner  *scanner.Scanner
	rules    *rules.Engine
	alerter  *alerter.Alerter
	interval time.Duration
}

// New creates a Watcher with the given rules engine, alerter, and poll interval.
func New(r *rules.Engine, a *alerter.Alerter, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  scanner.New(),
		rules:    r,
		alerter:  a,
		interval: interval,
	}
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("portwatch: starting watcher (interval=%s)", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch: watcher stopped")
			return
		case <-ticker.C:
			w.scan()
		}
	}
}

// scan performs a single scan cycle.
func (w *Watcher) scan() {
	ports, err := w.scanner.Scan()
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return
	}

	for _, p := range ports {
		action := w.rules.Evaluate(p)
		switch action {
		case rules.ActionAlert:
			w.alerter.EmitAlert(p)
		case rules.ActionAllow:
			w.alerter.EmitInfo(p)
		}
	}
}
