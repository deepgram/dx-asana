package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestEveryCommandHelpRunsWithoutPanic(t *testing.T) {
	commands := []string{
		"asana-cli",
		"asana-cli list",
		"asana-cli view",
		"asana-cli create",
		"asana-cli update",
		"asana-cli complete",
		"asana-cli delete",
		"asana-cli search",
		"asana-cli me",
		"asana-cli my-tasks",
		"asana-cli comment",
		"asana-cli workspaces",
		"asana-cli projects",
		"asana-cli sections",
		"asana-cli subtasks",
		"asana-cli config",
		"asana-cli config get",
		"asana-cli config set",
		"asana-cli config unset",
		"asana-cli config project",
		"asana-cli config project add",
		"asana-cli config project list",
		"asana-cli config project remove",
		"asana-cli config project switch",
		"asana-cli sync",
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
