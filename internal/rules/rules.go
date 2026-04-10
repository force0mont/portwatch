// Package rules provides configurable allow/deny rules for port monitoring.
package rules

import (
	"fmt"
	"strings"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionAlert Action = "alert"
)

// Rule represents a single port monitoring rule.
type Rule struct {
	Port     uint16 `yaml:"port"`
	Protocol string `yaml:"protocol"` // tcp or udp
	Action   Action `yaml:"action"`
	Comment  string `yaml:"comment,omitempty"`
}

// Engine evaluates open ports against a set of rules.
type Engine struct {
	rules []Rule
}

// New creates a new Engine with the provided rules.
func New(rules []Rule) *Engine {
	return &Engine{rules: rules}
}

// Evaluate checks a port/protocol pair against the rule set.
// It returns the matching Action and whether a rule was matched.
// If no rule matches, ActionAlert is returned as the default.
func (e *Engine) Evaluate(port uint16, protocol string) (Action, bool) {
	proto := strings.ToLower(protocol)
	for _, r := range e.rules {
		if r.Port == port && strings.ToLower(r.Protocol) == proto {
			return r.Action, true
		}
	}
	return ActionAlert, false
}

// Validate checks that all rules have valid fields.
func Validate(rules []Rule) error {
	for i, r := range rules {
		if r.Port == 0 {
			return fmt.Errorf("rule[%d]: port must be non-zero", i)
		}
		proto := strings.ToLower(r.Protocol)
		if proto != "tcp" && proto != "udp" {
			return fmt.Errorf("rule[%d]: protocol must be 'tcp' or 'udp', got %q", i, r.Protocol)
		}
		if r.Action != ActionAllow && r.Action != ActionAlert {
			return fmt.Errorf("rule[%d]: action must be 'allow' or 'alert', got %q", i, r.Action)
		}
	}
	return nil
}
