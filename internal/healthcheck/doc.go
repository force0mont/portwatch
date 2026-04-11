// Package healthcheck provides a minimal HTTP server that exposes a /healthz
// endpoint for liveness probing and metric inspection.
//
// Usage:
//
//	s := healthcheck.New(":9090", metricsInstance)
//	go s.ListenAndServe()
//	// later…
//	s.Shutdown()
//
// The /healthz endpoint returns HTTP 200 with a JSON body:
//
//	{
//	  "status": "ok",
//	  "uptime_sec": 42,
//	  "metrics": { … }
//	}
package healthcheck
