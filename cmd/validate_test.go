package cmd

import "testing"

func TestValidateDueDate(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"", false},
		{"2026-12-31", false},
		{"2026-01-01", false},
		{"2026-02-29", true},  // not a leap year
		{"2024-02-29", false}, // is a leap year
		{"2026-13-01", true},
		{"2026-12-32", true},
		{"tomorrow", true},
		{"2026/12/31", true},
		{"26-12-31", true},
		{"2026-12-31T00:00:00Z", true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			err := validateDueDate(tc.in)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateDueDate(%q) err=%v, wantErr=%v", tc.in, err, tc.wantErr)
			}
		})
	}
}
