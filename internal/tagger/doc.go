// Package tagger maps port numbers to human-readable service labels.
//
// It ships with a built-in table of well-known ports (SSH, HTTP, HTTPS,
// MySQL, PostgreSQL, Redis, MongoDB, …) and supports caller-supplied
// overrides so that site-specific services can be labelled clearly in
// alert output.
//
// Usage:
//
//	tg := tagger.New()
//	tg.Override(9200, "elasticsearch")
//
//	label := tg.Tag(443)   // "https"
//	label  = tg.Tag(9200)  // "elasticsearch"
//	label  = tg.Tag(9999)  // "unknown"
//
// Tagger is safe for concurrent use.
package tagger
