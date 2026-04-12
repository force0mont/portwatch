package dispatcher

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

// stubNotifier records every message it receives.
type stubNotifier struct {
	msgs []Message
	err  error
}

func (s *stubNotifier) Send(_ context.Context, msg Message) error {
	s.msgs = append(s.msgs, msg)
	return s.err
}

var fixedMsg = Message{Level: "alert", Title: "unexpected port", Body: "port 4444/tcp appeared"}

func TestDispatch_NoNotifiers_ReturnsError(t *testing.T) {
	d := New()
	if err := d.Dispatch(context.Background(), fixedMsg); err == nil {
		t.Fatal("expected error when no notifiers registered")
	}
}

func TestDispatch_SingleNotifier_ReceivesMessage(t *testing.T) {
	d := New()
	n := &stubNotifier{}
	d.Register(n)

	if err := d.Dispatch(context.Background(), fixedMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(n.msgs) != 1 {
		t.Fatalf("want 1 message, got %d", len(n.msgs))
	}
	if n.msgs[0] != fixedMsg {
		t.Errorf("want %+v, got %+v", fixedMsg, n.msgs[0])
	}
}

func TestDispatch_MultipleNotifiers_AllReceiveMessage(t *testing.T) {
	d := New()
	a, b := &stubNotifier{}, &stubNotifier{}
	d.Register(a)
	d.Register(b)

	_ = d.Dispatch(context.Background(), fixedMsg)

	if len(a.msgs) != 1 || len(b.msgs) != 1 {
		t.Errorf("want each notifier to receive 1 message; got a=%d b=%d", len(a.msgs), len(b.msgs))
	}
}

func TestDispatch_NotifierError_DoesNotAbortOthers(t *testing.T) {
	var buf bytes.Buffer
	d := newWithWriter(&buf)

	bad := &stubNotifier{err: errors.New("send failed")}
	good := &stubNotifier{}
	d.Register(bad)
	d.Register(good)

	if err := d.Dispatch(context.Background(), fixedMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(good.msgs) != 1 {
		t.Error("good notifier should still receive the message")
	}
	if buf.Len() == 0 {
		t.Error("expected error to be logged")
	}
}

func TestLen_ReflectsRegistrations(t *testing.T) {
	d := New()
	if d.Len() != 0 {
		t.Fatalf("want 0, got %d", d.Len())
	}
	d.Register(&stubNotifier{})
	d.Register(&stubNotifier{})
	if d.Len() != 2 {
		t.Fatalf("want 2, got %d", d.Len())
	}
}
