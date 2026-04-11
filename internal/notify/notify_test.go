package notify_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rodrwan/portwatch/internal/notify"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixedMsg(level notify.Level) notify.Message {
	return notify.Message{
		Level:     level,
		Port:      8080,
		Protocol:  "tcp",
		Addr:      "0.0.0.0",
		Timestamp: fixedTime,
		Detail:    "unexpected listener",
	}
}

func TestStderrNotifier_Send_Alert(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewStderrWithWriter(&buf)
	if err := n.Send(fixedMsg(notify.LevelAlert)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alert") {
		t.Errorf("expected 'alert' in output, got: %q", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port '8080' in output, got: %q", out)
	}
}

func TestStderrNotifier_Send_Info(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewStderrWithWriter(&buf)
	if err := n.Send(fixedMsg(notify.LevelInfo)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "info") {
		t.Errorf("expected 'info' in output")
	}
}

func TestWebhookNotifier_Send_Success(t *testing.T) {
	var received notify.Message
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewWebhook(ts.URL)
	msg := fixedMsg(notify.LevelAlert)
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Level != notify.LevelAlert {
		t.Errorf("expected level alert, got %s", received.Level)
	}
}

func TestWebhookNotifier_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewWebhook(ts.URL)
	if err := n.Send(fixedMsg(notify.LevelAlert)); err == nil {
		t.Error("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_Send_BadURL(t *testing.T) {
	n := notify.NewWebhook("http://127.0.0.1:0/no-server")
	if err := n.Send(fixedMsg(notify.LevelAlert)); err == nil {
		t.Error("expected connection error, got nil")
	}
}
