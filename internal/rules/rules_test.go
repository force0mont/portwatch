package rules_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/rules"
)

func defaultRules() []rules.Rule {
	return []rules.Rule{
		{Port: 22, Protocol: "tcp", Action: rules.ActionAllow, Comment: "SSH"},
		{Port: 80, Protocol: "tcp", Action: rules.ActionAllow, Comment: "HTTP"},
		{Port: 9999, Protocol: "tcp", Action: rules.ActionAlert, Comment: "suspicious"},
	}
}

func TestEvaluate_AllowedPort(t *testing.T) {
	e := rules.New(defaultRules())
	action, matched := e.Evaluate(22, "tcp")
	if !matched {
		t.Fatal("expected rule to match port 22/tcp")
	}
	if action != rules.ActionAllow {
		t.Errorf("expected allow, got %q", action)
	}
}

func TestEvaluate_AlertPort(t *testing.T) {
	e := rules.New(defaultRules())
	action, matched := e.Evaluate(9999, "tcp")
	if !matched {
		t.Fatal("expected rule to match port 9999/tcp")
	}
	if action != rules.ActionAlert {
		t.Errorf("expected alert, got %q", action)
	}
}

func TestEvaluate_NoMatch_DefaultsToAlert(t *testing.T) {
	e := rules.New(defaultRules())
	action, matched := e.Evaluate(12345, "tcp")
	if matched {
		t.Fatal("expected no rule to match port 12345/tcp")
	}
	if action != rules.ActionAlert {
		t.Errorf("expected default alert, got %q", action)
	}
}

func TestEvaluate_ProtocolMismatch(t *testing.T) {
	e := rules.New(defaultRules())
	// port 22 is defined for tcp only
	_, matched := e.Evaluate(22, "udp")
	if matched {
		t.Fatal("expected no match for port 22/udp")
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := rules.Validate(defaultRules()); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_ZeroPort(t *testing.T) {
	bad := []rules.Rule{{Port: 0, Protocol: "tcp", Action: rules.ActionAllow}}
	if err := rules.Validate(bad); err == nil {
		t.Error("expected error for zero port")
	}
}

func TestValidate_BadProtocol(t *testing.T) {
	bad := []rules.Rule{{Port: 80, Protocol: "icmp", Action: rules.ActionAllow}}
	if err := rules.Validate(bad); err == nil {
		t.Error("expected error for invalid protocol")
	}
}

func TestValidate_BadAction(t *testing.T) {
	bad := []rules.Rule{{Port: 80, Protocol: "tcp", Action: "block"}}
	if err := rules.Validate(bad); err == nil {
		t.Error("expected error for invalid action")
	}
}
