package alerting

import (
	"errors"
	"testing"
	"time"
)

func TestDedupeRemovesRepeatedIDs(t *testing.T) {
	alerts := []Alert{
		{ID: "a1", Message: "x"},
		{ID: "a1", Message: "x"},
		{ID: "a2", Message: "y"},
	}
	result := DedupeAlerts(alerts)
	if len(result) != 2 {
		t.Errorf("expected 2 deduped alerts, got %d", len(result))
	}
}

func TestDedupeStartsFreshEachCall(t *testing.T) {
	firstBatch := []Alert{{ID: "b1", Message: "first"}}
	secondBatch := []Alert{{ID: "b1", Message: "first"}} // same ID, independent call

	DedupeAlerts(firstBatch)
	result := DedupeAlerts(secondBatch)

	if len(result) != 1 {
		t.Errorf("a fresh call should not be affected by a previous call's state, got %d results", len(result))
	}
}

func TestComputeSeverityThresholds(t *testing.T) {
	cases := []struct {
		rate     float64
		expected string
	}{
		{0.01, "low"},
		{0.049, "low"},
		{0.05, "medium"},
		{0.10, "medium"},
		{0.15, "high"}, // boundary: 0.15 should already be "high", not "medium"
		{0.29, "high"},
		{0.30, "critical"},
		{0.50, "critical"},
	}
	for _, c := range cases {
		got := ComputeSeverity(c.rate)
		if got != c.expected {
			t.Errorf("ComputeSeverity(%.3f) = %q, want %q", c.rate, got, c.expected)
		}
	}
}

func TestRetrySucceedsEventually(t *testing.T) {
	calls := 0
	fn := func() error {
		calls++
		if calls < 2 {
			return errors.New("not yet")
		}
		return nil
	}
	err := RetryWithBackoff(fn, 3, time.Millisecond)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestRetryReturnsErrorAfterExhaustingAttempts(t *testing.T) {
	fn := func() error {
		return errors.New("nope")
	}
	err := RetryWithBackoff(fn, 3, time.Millisecond)
	if err == nil {
		t.Error("expected an error after exhausting retries, got nil")
	}
}

func TestParseLogLineBasic(t *testing.T) {
	line := "2024-01-15T10:30:00Z ERROR story-view-api database connection timeout"
	event, err := ParseLogLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("timestamp = %q", event.Timestamp)
	}
	if event.Level != "ERROR" {
		t.Errorf("level = %q", event.Level)
	}
	if event.Service != "story-view-api" {
		t.Errorf("service = %q", event.Service)
	}
	if event.Message != "database connection timeout" {
		t.Errorf("message = %q, want %q", event.Message, "database connection timeout")
	}
}
