// Package reporter provides periodic summary reporting for portwatch.
//
// A Reporter consumes a metrics.Snapshot and writes a human-readable or
// machine-readable (JSON) summary to any io.Writer, making it easy to
// redirect output to stdout, a log file, or a network sink.
//
// Supported formats:
//
//	 reporter.FormatText  – single-line key=value output
//	 reporter.FormatJSON  – newline-delimited JSON object
//
// Example:
//
//	r := reporter.New(reporter.FormatJSON)
//	r.Report(metricsInstance.Snapshot())
package reporter
