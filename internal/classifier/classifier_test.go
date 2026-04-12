package classifier_test

import (
	"testing"

	"github.com/jwhittle933/portwatch/internal/classifier"
	"github.com/jwhittle933/portwatch/internal/scanner"
)

func makePort(port uint16, addr string) scanner.Port {
	return scanner.Port{Port: port, Address: addr, Protocol: "tcp"}
}

func TestClassify_EphemeralPort_IsHighRisk(t *testing.T) {
	c := classifier.New()
	tier := c.Classify(makePort(45000, "0.0.0.0"))
	if tier != classifier.TierHigh {
		t.Fatalf("expected TierHigh, got %s", tier)
	}
}

func TestClassify_WellKnownPort_Loopback_IsLowRisk(t *testing.T) {
	c := classifier.New()
	tier := c.Classify(makePort(80, "127.0.0.1"))
	if tier != classifier.TierLow {
		t.Fatalf("expected TierLow, got %s", tier)
	}
}

func TestClassify_WellKnownPort_PublicAddr_IsMediumRisk(t *testing.T) {
	c := classifier.New()
	tier := c.Classify(makePort(443, "0.0.0.0"))
	if tier != classifier.TierMedium {
		t.Fatalf("expected TierMedium, got %s", tier)
	}
}

func TestClassify_RegisteredPort_IsMediumRisk(t *testing.T) {
	c := classifier.New()
	tier := c.Classify(makePort(8080, "0.0.0.0"))
	if tier != classifier.TierMedium {
		t.Fatalf("expected TierMedium, got %s", tier)
	}
}

func TestClassify_CustomEphemeralRange(t *testing.T) {
	c := classifier.NewWithEphemeralRange(1024, 2048)
	tier := c.Classify(makePort(1500, "0.0.0.0"))
	if tier != classifier.TierHigh {
		t.Fatalf("expected TierHigh for custom ephemeral range, got %s", tier)
	}
}

func TestTier_String(t *testing.T) {
	cases := []struct {
		tier classifier.Tier
		want string
	}{
		{classifier.TierLow, "low"},
		{classifier.TierMedium, "medium"},
		{classifier.TierHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.tier.String(); got != tc.want {
			t.Errorf("Tier(%d).String() = %q, want %q", tc.tier, got, tc.want)
		}
	}
}
