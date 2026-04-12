package pipeline_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

const fakeTCP = `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 0 1 0000000000000000 100 0 0 10 0
`

func TestIntegration_Pipeline_ProcessesFakeScan(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "net"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "net", "tcp"), []byte(fakeTCP), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	s, err := scanner.NewWithProcPath(dir)
	if err != nil {
		t.Fatalf("scanner: %v", err)
	}
	st := state.New()
	r, _ := rules.New(nil)
	f, _ := filter.New(nil)
	d := dedupe.New(200 * time.Millisecond)
	rl := ratelimit.New(100, time.Second)
	a := alerter.NewWithWriter(&buf)

	p := pipeline.New(pipeline.Config{
		Scanner:   s,
		State:     st,
		Rules:     r,
		Filter:    f,
		Dedupe:    d,
		RateLimit: rl,
		Alerter:   a,
		Interval:  15 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	p.Run(ctx) //nolint:errcheck

	if buf.Len() == 0 {
		t.Fatal("expected at least one alert to be written")
	}
}
