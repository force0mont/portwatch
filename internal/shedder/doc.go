// Package shedder provides a load-shedding gate for port scan events.
//
// Under normal conditions all events pass through. When the rate of
// incoming events exceeds the configured per-window maximums, lower-
// priority events are dropped so that high-priority alerts are never
// lost and downstream consumers are not overwhelmed.
//
// Three priority tiers are defined:
//
//	PriorityHigh   – never shed; reserved for critical alerts.
//	PriorityNormal – shed once the normal quota is exhausted.
//	PriorityLow    – shed once the low quota is exhausted.
//
// Each Shedder is safe for concurrent use.
package shedder
