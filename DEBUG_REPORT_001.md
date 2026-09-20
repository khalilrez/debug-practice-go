# Debug Practice: Go

# Bugs and Issues

## `alerting.DedupeAlerts` doesn't work as expected

The `DedupeAlerts` function doesn't work as expected. It should only dedupe
alerts whose IDs haven't been seen before within a batch. However, it's
currently remembering IDs seen in a previous, unrelated batch.
### Solution
Use a fresh map for each batch.

## `alerting.ComputeSeverity` doesn't work as expected

The `ComputeSeverity` function doesn't work as expected. It should map an
error rate (0.0-1.0) to a severity label.

Thresholds:
	< 0.05          -> "low"
	0.05 to < 0.15  -> "medium"
	0.15 to < 0.30  -> "high"
	>= 0.30         -> "critical"
### Solution
Fix the thresholds.

## `alerting.RetryWithBackoff` doesn't work as expected

The `RetryWithBackoff` function doesn't work as expected. It should call `fn`
up to `retries` times, doubling the delay between attempts each time, starting
at `baseDelay`. If `fn` succeeds, it should return nil immediately. If all
attempts fail, it should return the last error encountered.
### Solution
Fix the implementation by returning the last error encountered instead of
nil.

## `alerting.ParseLogLine` doesn't work as expected

The `ParseLogLine` function doesn't work as expected. It should parse a log
line of the form:

	"2024-01-15T10:30:00Z ERROR story-view-api database connection timeout"

The message portion may contain multiple words.
### Solution
Use `strings.Fields` to split the line into its constituent parts then join them with a space. This will correctly handle the case where the message contains multiple words.

