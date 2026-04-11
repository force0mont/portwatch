// Package healthcheck provides a simple HTTP health endpoint that exposes
// daemon liveness and a snapshot of current metrics for external monitoring.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Server is a lightweight HTTP server that serves a /healthz endpoint.
type Server struct {
	addr    string
	metrics *metrics.Metrics
	server  *http.Server
	started time.Time
}

// response is the JSON body returned by the health endpoint.
type response struct {
	Status    string          `json:"status"`
	UptimeSec int64           `json:"uptime_sec"`
	Metrics   metrics.Snapshot `json:"metrics"`
}

// New creates a Server that will listen on addr (e.g. ":9090").
func New(addr string, m *metrics.Metrics) *Server {
	s := &Server{
		addr:    addr,
		metrics: m,
		started: time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown() error {
	return s.server.Close()
}

// Addr returns the configured listen address.
func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	snap := s.metrics.Snapshot()
	resp := response{
		Status:    "ok",
		UptimeSec: int64(time.Since(s.started).Seconds()),
		Metrics:   snap,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
	}
}
