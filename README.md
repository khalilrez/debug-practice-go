# Alerting Utilities in Go

A small Go package for practicing debugging with realistic alert-processing
helpers. It covers four common tasks: deduplicating alerts, assigning severity
from an error rate, retrying transient failures, and parsing structured log
lines.

The repository contains the corrected implementation and a test suite that
documents its expected behavior. A summary of the original bugs and their
fixes is available in [`DEBUG_REPORT_001.md`](DEBUG_REPORT_001.md).
Go back in commit history to see the original bugs and their fixes.

## Requirements

- Go 1.22.2 or later

No third-party dependencies are required.

## Run the tests

From the repository root:

```bash
go test -v ./...
```

## Package overview

### Deduplicate alerts

`DedupeAlerts` removes repeated alert IDs within a single batch while
preserving the first occurrence and its input order. Each call starts with a
fresh deduplication scope.

```go
alerts := []alerting.Alert{
	{ID: "db-1", Message: "database unavailable"},
	{ID: "db-1", Message: "database unavailable"},
	{ID: "api-1", Message: "high latency"},
}

unique := alerting.DedupeAlerts(alerts) // two alerts
```

### Compute severity

`ComputeSeverity` maps an error rate to a label using these boundaries:

| Error rate | Severity |
| --- | --- |
| `< 0.05` | `low` |
| `0.05` to `< 0.15` | `medium` |
| `0.15` to `< 0.30` | `high` |
| `>= 0.30` | `critical` |

### Retry an operation

`RetryWithBackoff` calls a function up to the requested number of attempts.
The delay starts at `baseDelay` and doubles after each failed attempt. It
returns immediately on success or returns the last error when all attempts
fail.

```go
err := alerting.RetryWithBackoff(func() error {
	return sendAlert()
}, 3, 100*time.Millisecond)
```

### Parse a log line

`ParseLogLine` converts a whitespace-delimited line into a `LogEvent`:

```text
2024-01-15T10:30:00Z ERROR story-view-api database connection timeout
```

The first three fields become the timestamp, level, and service. Everything
after them becomes the message, so multi-word messages are preserved. Lines
with fewer than four fields return an error.

## Project layout

```text
.
├── alerting.go          # Package implementation
├── alerting_test.go     # Behavioral tests
├── DEBUG_REPORT_001.md  # Original bug report and fixes
└── go.mod
```
