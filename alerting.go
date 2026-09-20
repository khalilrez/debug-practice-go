// Package alerting provides small helpers for processing incoming alerts:
// deduplicating repeated alerts, computing severity from an error rate,
// retrying flaky operations with backoff, and parsing raw log lines.
package alerting

import (
	"fmt"
	"strings"
	"time"
)

type Alert struct {
	ID      string
	Message string
}

// seenAlertIDs is used by DedupeAlerts.
// NOTE: each independent call to DedupeAlerts should be treated as a fresh
// batch — it should NOT remember IDs seen in a previous, unrelated call.

// DedupeAlerts returns only the alerts whose ID hasn't been seen before
// within this batch.
func DedupeAlerts(alerts []Alert) []Alert {
	var seenAlertIDs = map[string]bool{}
	result := []Alert{}
	for _, a := range alerts {
		if !seenAlertIDs[a.ID] {
			seenAlertIDs[a.ID] = true
			result = append(result, a)
		}
	}
	return result
}

// ComputeSeverity maps an error rate (0.0-1.0) to a severity label.
//
// Thresholds:
//
//	< 0.05          -> "low"
//	0.05 to < 0.15  -> "medium"
//	0.15 to < 0.30  -> "high"
//	>= 0.30         -> "critical"
func ComputeSeverity(errorRate float64) string {
	if errorRate < 0.05 {
		return "low"
	} else if errorRate < 0.15 {
		return "medium"
	} else if errorRate < 0.30 {
		return "high"
	}
	return "critical"
}

// RetryWithBackoff calls fn up to `retries` times, doubling the delay
// between attempts each time, starting at baseDelay. If fn succeeds, it
// returns nil immediately. If all attempts fail, it should return the
// last error encountered.
func RetryWithBackoff(fn func() error, retries int, baseDelay time.Duration) error {
	delay := baseDelay
	var lastErr error
	for attempt := 0; attempt < retries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		time.Sleep(delay)
		delay *= 2
	}
	// bug: swallows the failure instead of returning lastErr
	return lastErr
}

type LogEvent struct {
	Timestamp string
	Level     string
	Service   string
	Message   string
}

// ParseLogLine parses a line of the form:
//
//	"2024-01-15T10:30:00Z ERROR story-view-api database connection timeout"
//
// The message portion may contain multiple words.
func ParseLogLine(line string) (LogEvent, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return LogEvent{}, fmt.Errorf("malformed log line: %q", line)
	}
	return LogEvent{
		Timestamp: parts[0],
		Level:     parts[1],
		Service:   parts[2],
		Message:   strings.Join(parts[3:], " "), // bug: only grabs the first word of the message
	}, nil
}
