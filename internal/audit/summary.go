package audit

import (
	"fmt"
	"strings"
	"time"
)

// Summary holds aggregated statistics from audit log entries.
type Summary struct {
	TotalRuns    int
	SuccessRuns  int
	FailedRuns   int
	SecretsAdded int
	SecretsUpdated int
	SecretsDeleted int
	FirstRun     *time.Time
	LastRun      *time.Time
}

// Summarize computes a Summary from a slice of audit entries.
func Summarize(entries []Entry) Summary {
	var s Summary
	for _, e := range entries {
		s.TotalRuns++
		if e.Error == "" {
			s.SuccessRuns++
		} else {
			s.FailedRuns++
		}
		s.SecretsAdded += e.Added
		s.SecretsUpdated += e.Updated
		s.SecretsDeleted += e.Deleted
		t := e.Timestamp
		if s.FirstRun == nil || t.Before(*s.FirstRun) {
			s.FirstRun = &t
		}
		if s.LastRun == nil || t.After(*s.LastRun) {
			s.LastRun = &t
		}
	}
	return s
}

// Format returns a human-readable summary string.
func (s Summary) Format() string {
	if s.TotalRuns == 0 {
		return "No audit entries found."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Total runs:     %d\n", s.TotalRuns)
	fmt.Fprintf(&b, "  Successful:   %d\n", s.SuccessRuns)
	fmt.Fprintf(&b, "  Failed:       %d\n", s.FailedRuns)
	fmt.Fprintf(&b, "Secrets added:  %d\n", s.SecretsAdded)
	fmt.Fprintf(&b, "Secrets updated:%d\n", s.SecretsUpdated)
	fmt.Fprintf(&b, "Secrets deleted:%d\n", s.SecretsDeleted)
	if s.FirstRun != nil {
		fmt.Fprintf(&b, "First run:      %s\n", s.FirstRun.Format(time.RFC3339))
	}
	if s.LastRun != nil {
		fmt.Fprintf(&b, "Last run:       %s\n", s.LastRun.Format(time.RFC3339))
	}
	return b.String()
}
