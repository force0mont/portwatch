// Package digest provides snapshot fingerprinting for portwatch.
//
// A Digester computes a deterministic SHA-256 hash over a set of open-port
// observations (proto, address, port triples) and tracks whether the hash
// has changed between consecutive calls. This lets the watcher skip
// downstream processing — rule evaluation, alerting, state diffing — when
// nothing about the listening set has changed since the last scan cycle.
//
// Entries are sorted before hashing so that the digest is independent of
// the order in which the kernel returns rows from /proc/net/tcp.
//
// Usage:
//
//	d := digest.New()
//	if d.Changed(entries) {
//		// process new snapshot
//	}
package digest
