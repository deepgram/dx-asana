package cmd

import (
	"testing"
)

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd is nil")
	}

	if rootCmd.Use != "dx-asana" {
		t.Errorf("unexpected command name: %s", rootCmd.Use)
	}
}

func TestPersistentFlags(t *testing.T) {
	flags := rootCmd.PersistentFlags()

	if flags.Lookup("json") == nil {
		t.Error("json flag not found")
	}

	if flags.Lookup("token") == nil {
		t.Error("token flag not found")
	}

	if flags.Lookup("workspace") == nil {
		t.Error("workspace flag not found")
	}

	if flags.Lookup("project") == nil {
		t.Error("project flag not found")
	}
}