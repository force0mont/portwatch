// Package enricher decorates raw scanner.Port values with contextual
// metadata such as reverse-DNS hostnames and IANA service names.
//
// Usage:
//
//	e := enricher.New()
//	entry := e.Enrich(port)
//	fmt.Println(entry.Label) // e.g. "tcp/22 (ssh)"
//
// The enricher is intentionally best-effort: DNS failures and unknown port
// numbers are silently ignored so that the monitoring pipeline is never
// blocked by external resolver unavailability.
package enricher
