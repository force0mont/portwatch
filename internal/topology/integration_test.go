package topology_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/topology"
)

func TestConcurrent_BuildAndGroups_NoRace(t *testing.T) {
	topo := topology.New()
	ports := makePorts()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			topo.Build(ports)
		}()
		go func() {
			defer wg.Done()
			_ = topo.Groups()
		}()
	}
	wg.Wait()
}

func TestTopology_EmptyBuild_ZeroLen(t *testing.T) {
	topo := topology.New()
	topo.Build(makePorts())
	topo.Build([]scanner.Port{})
	if topo.Len() != 0 {
		t.Fatalf("expected 0 after empty build, got %d", topo.Len())
	}
	if len(topo.Groups()) != 0 {
		t.Fatal("expected no groups after empty build")
	}
}
