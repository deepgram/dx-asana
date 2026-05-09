package cmd

import (
	"fmt"
	"regexp"
	"time"
)

var dueDatePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// validateDueDate ensures a --due flag value is YYYY-MM-DD and parses as a
// real calendar date. Empty input is treated as "not provided" and skipped.
func validateDueDate(s string) error {
	if s == "" {
		return nil
	}
	if !dueDatePattern.MatchString(s) {
		return fmt.Errorf("--due must be YYYY-MM-DD (got %q)", s)
	}
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return fmt.Errorf("--due is not a valid calendar date: %s", s)
	}
	return nil
}
