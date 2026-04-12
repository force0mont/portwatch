package preselector_test

import (
	"testing"

	"github.com/user/portwatch/internal/preselector"
	"github.com/user/portwatch/internal/scanner"
)

func ports(pairs ...any) []scanner.Port {
	var out []scanner.Port
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, scanner.Port{
			Port:     pairs[i].(uint16),
			Protocol: pairs[i+1].(string),
		})
	}
	return out
}

func TestFilter_EmptyIgnoreList_PassesAll(t *testing.T) {
	p := preselector.New()
	input := ports(uint16(80), "tcp", uint16(443), "tcp")
	got := p.Filter(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestFilter_IgnoredPort_IsDropped(t *testing.T) {
	p := preselector.New()
	p.Ignore(uint16(22), "tcp")

	input := ports(uint16(22), "tcp", uint16(8080), "tcp")
	got := p.Filter(input)

	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", got[0].Port)
	}
}

func TestFilter_ProtocolMismatch_NotDropped(t *testing.T) {
	p := preselector.New()
	p.Ignore(uint16(53), "tcp")

	input := ports(uint16(53), "udp")
	got := p.Filter(input)

	if len(got) != 1 {
		t.Fatalf("expected udp/53 to pass through, got %d entries", len(got))
	}
}

func TestRemove_UnregistersIgnoredPair(t *testing.T) {
	p := preselector.New()
	p.Ignore(uint16(3306), "tcp")
	p.Remove(uint16(3306), "tcp")

	input := ports(uint16(3306), "tcp")
	got := p.Filter(input)

	if len(got) != 1 {
		t.Fatalf("expected port to pass after Remove, got %d entries", len(got))
	}
}

func TestLen_ReflectsIgnoredCount(t *testing.T) {
	p := preselector.New()
	if p.Len() != 0 {
		t.Fatalf("expected 0, got %d", p.Len())
	}
	p.Ignore(uint16(80), "tcp")
	p.Ignore(uint16(443), "tcp")
	if p.Len() != 2 {
		t.Fatalf("expected 2, got %d", p.Len())
	}
	p.Remove(uint16(80), "tcp")
	if p.Len() != 1 {
		t.Fatalf("expected 1 after Remove, got %d", p.Len())
	}
}

func TestFilter_EmptyInput_ReturnsNonNil(t *testing.T) {
	p := preselector.New()
	got := p.Filter([]scanner.Port{})
	if got == nil {
		t.Fatal("expected non-nil slice for empty input")
	}
}
