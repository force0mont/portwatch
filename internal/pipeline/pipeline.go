// Package pipeline wires together the scanning, filtering, deduplication,
// rate-limiting and alerting stages into a single reusable processing chain.
package pipeline

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Stage is a function that processes a scanned port entry and returns whether
// processing should continue down the pipeline.
type Stage func(entry scanner.Entry) bool

// Pipeline holds the ordered chain of stages and supporting components.
type Pipeline struct {
	scanner  *scanner.Scanner
	state    *state.State
	rules    *rules.Engine
	filter   *filter.Filter
	dedupe   *dedupe.Dedupe
	ratelimit *ratelimit.Limiter
	alerter  *alerter.Alerter
	interval time.Duration
}

// Config carries the dependencies needed to construct a Pipeline.
type Config struct {
	Scanner   *scanner.Scanner
	State     *state.State
	Rules     *rules.Engine
	Filter    *filter.Filter
	Dedupe    *dedupe.Dedupe
	RateLimit *ratelimit.Limiter
	Alerter   *alerter.Alerter
	Interval  time.Duration
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) *Pipeline {
	return &Pipeline{
		scanner:   cfg.Scanner,
		state:     cfg.State,
		rules:     cfg.Rules,
		filter:    cfg.Filter,
		dedupe:    cfg.Dedupe,
		ratelimit: cfg.RateLimit,
		alerter:   cfg.Alerter,
		interval:  cfg.Interval,
	}
}

// Run starts the pipeline loop, scanning at the configured interval until ctx
// is cancelled.
func (p *Pipeline) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			p.tick()
		}
	}
}

func (p *Pipeline) tick() {
	entries, err := p.scanner.Scan()
	if err != nil {
		return
	}
	events := p.state.Diff(entries)
	for _, ev := range events {
		if p.filter.Suppressed(ev.Entry) {
			continue
		}
		if p.dedupe.IsDuplicate(ev.Entry) {
			continue
		}
		key := ev.Entry.Addr + ev.Entry.Protocol
		if !p.ratelimit.Allow(key) {
			continue
		}
		action := p.rules.Evaluate(ev.Entry)
		if action == rules.ActionAlert {
			p.alerter.EmitAlert(ev)
		} else {
			p.alerter.EmitInfo(ev)
		}
	}
}
