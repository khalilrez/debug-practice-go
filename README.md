# Debug Practice Exercise (Go)

A small `alerting` package has a handful of bugs. Your job: get `go test -v ./...`
fully green by fixing `alerting.go` only — do not edit `alerting_test.go`.

## Setup

Requires Go installed (1.21+).

```bash
go test -v ./...
```

You should see some tests failing. Read the failure output carefully — it
tells you the expected vs actual value, which is usually enough to point
you at the bug without needing to guess blindly.

## How to approach it (practice this process, not just the fixes)

1. Run `go test -v ./...` first, before touching any code, and read every
   failure message.
2. Pick one failing test at a time. Read the doc comment above the
   relevant function in `alerting.go` — it describes the *intended*
   behavior.
3. Form a hypothesis about the bug before changing code. Say it out loud
   (or write it down): "I think this fails because ___."
4. Make the smallest possible fix, then rerun `go test -v ./...` to
   confirm — don't fix all four at once and hope.
5. Once everything's green, go back and explain out loud *why* each bug
   happened — this is the part interviewers actually care about, not
   just that you got it working.

## What's being tested (don't peek until you're stuck)

There are 4 distinct bugs, each a somewhat classic/realistic category:

- Shared/package-level state leaking across calls that should be
  independent (a very common Go footgun — package vars holding state
  that should really be scoped to a single call or struct)
- An off-by-one / boundary condition
- A function that fails silently instead of surfacing an error
  (returns `nil` when it should return the real error)
- A string-parsing assumption (`strings.Fields`) that breaks once the
  message has more than one word

If you get stuck for more than ~10-15 minutes on any one, that's a fine
moment to ask for a hint rather than grinding — in a real assessment
you'd likely timebox similarly.

## Bonus, if you finish early

The `DedupeAlerts` bug is fixable a few different ways in Go — a local
map instead of a package var is the minimal fix, but think about how
you'd design this if it needed to track real cross-request state safely
(e.g. concurrent goroutines calling it at once). What would break, and
what would you reach for instead of a plain `map[string]bool`?
