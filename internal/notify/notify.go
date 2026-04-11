// Package notify provides pluggable notification backends for portwatch alerts.
// Supported channels: stderr (default), webhook (HTTP POST).
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelAlert Level = "alert"
	LevelInfo  Level = "info"
)

// Message is the payload sent to a notification backend.
type Message struct {
	Level     Level     `json:"level"`
	Port      uint16    `json:"port"`
	Protocol  string    `json:"protocol"`
	Addr      string    `json:"addr"`
	Timestamp time.Time `json:"timestamp"`
	Detail    string    `json:"detail,omitempty"`
}

// Notifier sends a notification message to a backend.
type Notifier interface {
	Send(msg Message) error
}

// StderrNotifier writes formatted alerts to an io.Writer (defaults to os.Stderr).
type StderrNotifier struct {
	w io.Writer
}

// NewStderr returns a StderrNotifier writing to os.Stderr.
func NewStderr() *StderrNotifier { return &StderrNotifier{w: os.Stderr} }

// NewStderrWithWriter returns a StderrNotifier writing to w (for testing).
func NewStderrWithWriter(w io.Writer) *StderrNotifier { return &StderrNotifier{w: w} }

// Send writes a human-readable line to the configured writer.
func (s *StderrNotifier) Send(msg Message) error {
	_, err := fmt.Fprintf(s.w, "[%s] %s %s:%d %s\n",
		msg.Timestamp.Format(time.RFC3339),
		msg.Level,
		msg.Addr, msg.Port,
		msg.Protocol,
	)
	return err
}

// WebhookNotifier POSTs JSON-encoded messages to a remote URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhook returns a WebhookNotifier that posts to url.
func NewWebhook(url string) *WebhookNotifier {
	return &WebhookNotifier{url: url, client: &http.Client{Timeout: 5 * time.Second}}
}

// Send marshals msg to JSON and POSTs it to the webhook URL.
func (w *WebhookNotifier) Send(msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("notify: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
