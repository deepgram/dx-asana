package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestEveryCommandHelpRunsWithoutPanic(t *testing.T) {
	commands := []string{
		"dx-asana",
		"dx-asana list",
		"dx-asana view",
		"dx-asana create",
		"dx-asana update",
		"dx-asana complete",
		"dx-asana delete",
		"dx-asana search",
		"dx-asana me",
		"dx-asana my-tasks",
		"dx-asana comment",
		"dx-asana workspaces",
		"dx-asana projects",
		"dx-asana sections",
		"dx-asana subtasks",
		"dx-asana config",
		"dx-asana config get",
		"dx-asana config set",
		"dx-asana config unset",
		"dx-asana config project",
		"dx-asana config project add",
		"dx-asana config project list",
		"dx-asana config project remove",
		"dx-asana config project switch",
		"dx-asana sync",
	}

	for _, line := range commands {
		t.Run(line, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			parts := strings.Fields(line)
			args := append(parts[1:], "--help")
			rootCmd.SetOut(&stdout)
			rootCmd.SetErr(&stderr)
			rootCmd.SetArgs(args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("%s --help returned error: %v\nstderr: %s", line, err, stderr.String())
			}
			out := stdout.String()
			if out == "" {
				t.Errorf("%s --help produced no output", line)
			}
			if !strings.Contains(out, "Usage:") && !strings.Contains(out, "Available Commands:") {
				t.Errorf("%s --help missing Usage/Commands section:\n%s", line, out)
			}
		})
	}
}
