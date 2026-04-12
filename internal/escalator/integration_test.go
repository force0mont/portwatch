package escalator_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stevezaluk/portwatch/internal/escalator"
)

func TestConcurrent_Record_NoRace(t *testing.T) {
	e := escalator.New(time.Minute, 3, 6)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.Record("tcp:8080")
		}()
	}
	wg.Wait()
}

func TestEscalation_ProgressesThroughLevels(t *testing.T) {
	e := escalator.New(time.Minute, 2, 4)

	levels := make([]escalator.Level, 5)
	for i := range levels {
		levels[i] = e.Record("tcp:22")
	}

	if levels[0] != escalator.LevelNone {
		t.Errorf("hit 1: expected LevelNone, got %d", levels[0])
	}
	if levels[1] != escalator.LevelWarning {
		t.Errorf("hit 2: expected LevelWarning, got %d", levels[1])
	}
	if levels[3] != escalator.LevelCritical {
		t.Errorf("hit 4: expected LevelCritical, got %d", levels[3])
	}
	if levels[4] != escalator.LevelCritical {
		t.Errorf("hit 5: expected LevelCritical (sticky), got %d", levels[4])
	}
}
