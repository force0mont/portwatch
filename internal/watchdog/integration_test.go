package watchdog_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestIntegration_Watchdog_FiresOnRealStall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	wd := watchdog.New(80 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go wd.Run(ctx)
	// Intentionally do not call Beat — expect the watchdog to fire.
	select {
	case <-wd.Stuck():
		// good
	case <-ctx.Done():
		t.Fatal("watchdog never fired during real stall")
	}
}

func TestIntegration_Watchdog_SilentWithRegularBeats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	wd := watchdog.New(150 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Run(ctx)

	end := time.After(400 * time.Millisecond)
	beat := time.NewTicker(40 * time.Millisecond)
	defer beat.Stop()
	for {
		select {
		case <-beat.C:
			wd.Beat()
		case <-wd.Stuck():
			t.Fatal("watchdog fired despite regular beats")
		case <-end:
			return
		}
	}
}
