// Package main is the entry point for the portwatch CLI daemon.
//
// Usage:
//
//	portwatch [flags]
//
// Flags:
//
//	-config string
//		Path to a JSON configuration file. When omitted, built-in
//		defaults are used (see internal/config for details).
//
//	-version
//		Print the build version and exit.
//
// Signals:
//
//	SIGINT / SIGTERM  — graceful shutdown.
//
// portwatch reads /proc/net/tcp (and /proc/net/udp) on each scan interval,
// compares the open-port snapshot against the previous one, evaluates each
// port against the configured allow/alert rules, and emits structured JSON
// events to stdout.
package main
