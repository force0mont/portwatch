// Package ledger maintains a frequency table of port appearance and
// disappearance events observed by the portwatch scanner.
//
// Each (protocol, port) pair is tracked independently. Callers record
// events via RecordAppeared and RecordDisappeared, then query totals
// via Get or All.
//
// Ledger is safe for concurrent use.
package ledger
