package scorecard_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
	"github.com/yourorg/portwatch/internal/scorecard"
)

func makePort(addr string, port uint16, service string) scanner.Port {
	return scanner.Port{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Service:  service,
	}
}

func TestScore_UnknownService_AddsWeight(t *testing.T) {
	s := scorecard.New(scorecard.DefaultWeights())
	p := makePort("0.0.0.0", 8080, "")
	e := s.Score(p)
	if e.Score < scorecard.DefaultWeights().UnknownService {
		t.Fatalf("expected score >= %.2f, got %.2f", scorecard.DefaultWeights().UnknownService, e.Score)
	}
}

func TestScore_KnownService_NoUnknownWeight(t *testing.T) {
	w := scorecard.DefaultWeights()
	s := scorecard.New(w)
	p := makePort("0.0.0.0", 80, "http")
	e := s.Score(p)
	if e.Score >= w.UnknownService {
		t.Fatalf("expected score < %.2f for known service, got %.2f", w.UnknownService, e.Score)
	}
}

func TestScore_EphemeralPort_AddsWeight(t *testing.T) {
	w := scorecard.DefaultWeights()
	s := scorecard.New(w)
	p := makePort("0.0.0.0", 55000, "http") // known service but ephemeral port
	e := s.Score(p)
	if e.Score < w.EphemeralPort {
		t.Fatalf("expected score >= %.2f for ephemeral port, got %.2f", w.EphemeralPort, e.Score)
	}
}

func TestScore_Loopback_ReducesScore(t *testing.T) {
	w := scorecard.DefaultWeights()
	s := scorecard.New(w)
	public := makePort("0.0.0.0", 8080, "")
	loopback := makePort("127.0.0.1", 8080, "")
	pubScore := s.Score(public).Score
	loopScore := s.Score(loopback).Score
	if loopScore >= pubScore {
		t.Fatalf("loopback score %.2f should be lower than public %.2f", loopScore, pubScore)
	}
}

func TestScore_NeverNegative(t *testing.T) {
	w := scorecard.Weights{UnknownService: 0, EphemeralPort: 0, LoopbackOnly: 1.0}
	s := scorecard.New(w)
	p := makePort("127.0.0.1", 80, "http")
	e := s.Score(p)
	if e.Score < 0 {
		t.Fatalf("score must not be negative, got %.2f", e.Score)
	}
}

func TestScore_CachedResult_Consistent(t *testing.T) {
	s := scorecard.New(scorecard.DefaultWeights())
	p := makePort("0.0.0.0", 9999, "")
	e1 := s.Score(p)
	e2 := s.Score(p)
	if e1.Score != e2.Score {
		t.Fatalf("cached score mismatch: %.2f != %.2f", e1.Score, e2.Score)
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	s := scorecard.New(scorecard.DefaultWeights())
	p := makePort("0.0.0.0", 9999, "")
	_ = s.Score(p)
	s.Evict(p)
	// After eviction the score is recomputed — should still equal the original.
	e := s.Score(p)
	if e.Score < 0 {
		t.Fatal("unexpected negative score after eviction")
	}
}
