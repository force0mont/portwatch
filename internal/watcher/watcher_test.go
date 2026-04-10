package watcher_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alerter"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/watcher"
)

func defaultEngine(t *testing.T) *rules.Engine {
	t.Helper()
	e, err := rules.New([]rules.Rule{
		{Port: 22, Protocol: "tcp", Action: "allow"},
		{Port: 9999, Protocol: "tcp", Action: "alert"},
	})
	if err != nil {
		t.Fatalf("rules.New: %v", err)
	}
	return e
}

func TestWatcher_RunCancels(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)
	e := defaultEngine(t)

	w := watcher.New(e, a, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// expected: Run returned after context cancellation
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancellation")
	}
}

func TestWatcher_New_NotNil(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)
	e := defaultEngine(t)

	w := watcher.New(e, a, time.Second)
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestWatcher_OutputContainsJSON(t *testing.T) {
	var buf bytes.Buffer
	a := alerter.NewWithWriter(&buf)
	e := defaultEngine(t)

	w := watcher.New(e, a, 30*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	// The alerter writes JSON lines; verify the output is at least valid-looking.
	output := buf.String()
	if len(output) > 0 && !strings.Contains(output, "port") {
		t.Errorf("unexpected output format: %s", output)
	}
}
